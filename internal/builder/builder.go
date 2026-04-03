package builder

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuxeos3-dot/wayao/internal/config"
	"github.com/yuxeos3-dot/wayao/internal/schema"
	"github.com/yuxeos3-dot/wayao/internal/variation"
)

type Builder struct {
	DB  *sql.DB
	Cfg *config.Config
}

func New(db *sql.DB, cfg *config.Config) *Builder {
	return &Builder{DB: db, Cfg: cfg}
}

type siteInfo struct {
	ID             int64
	Domain         string
	TemplateID     int64
	TemplateSlug   string
	TemplatePath   string
	Market         string
	Language        string
	KeywordType    string
	PrimaryKeyword string
	RedirectURL    string
	BrandName       string
	BrandColor      string
	CTAText         string
	CTASub          string
	HeroTitle       string
	HeroSubtitle    string
	IntroText       string
	BodyContent     string
	Conclusion      string
	PageTitle       string
	MetaDesc        string
	H1              string
	FAQTitle        string
	FAQItems        string
	ExtraData       string
	AuthorName      string
	AuthorTitle     string
	AuthorBio       string
	TrustBadges     string
	Disclosure      string
	Disclaimer      string
	Feature1Icon    string
	Feature1Title   string
	Feature1Desc    string
	Feature2Icon    string
	Feature2Title   string
	Feature2Desc    string
	Feature3Icon    string
	Feature3Title   string
	Feature3Desc    string
}

func (b *Builder) BuildSite(domainID int64) error {
	start := time.Now()
	logID := b.createBuildLog(domainID, "building")

	site, err := b.loadSiteInfo(domainID)
	if err != nil {
		b.finishBuildLog(logID, "error", err.Error(), time.Since(start))
		return fmt.Errorf("load site: %w", err)
	}

	// 1. Create temp build directory
	buildDir := filepath.Join(b.Cfg.DataDir, "builds", site.Domain)
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		b.finishBuildLog(logID, "error", err.Error(), time.Since(start))
		return err
	}

	// 2. Copy template to build dir
	templateSrc := filepath.Join(b.Cfg.TemplDir, site.TemplateSlug)
	if err := copyDir(templateSrc, buildDir); err != nil {
		b.finishBuildLog(logID, "error", fmt.Sprintf("copy template: %v", err), time.Since(start))
		return err
	}

	// 3. Generate Hugo config with all content injected
	if err := b.generateConfig(buildDir, site); err != nil {
		b.finishBuildLog(logID, "error", fmt.Sprintf("gen config: %v", err), time.Since(start))
		return err
	}

	// 4. Generate support pages
	if err := b.generateSupportPages(buildDir, site); err != nil {
		log.Printf("[BUILD] warn: support pages: %v", err)
	}

	// 5. Run Hugo
	outputDir := filepath.Join(b.Cfg.BuildOut, site.Domain)
	os.MkdirAll(outputDir, 0755)

	hugoPath := b.Cfg.HugoPath
	if hugoPath == "" {
		hugoPath = "hugo"
	}

	cmd := exec.Command(hugoPath, "--minify", "-d", outputDir)
	cmd.Dir = buildDir
	cmd.Env = append(os.Environ(), "HUGO_ENVIRONMENT=production")

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.finishBuildLog(logID, "error", fmt.Sprintf("hugo: %s\n%v", string(output), err), time.Since(start))
		return fmt.Errorf("hugo build: %w\n%s", err, string(output))
	}

	// 6. Generate robots.txt and sitemap
	robotsTxt := schema.GenerateRobotsTxt(site.Domain)
	os.WriteFile(filepath.Join(outputDir, "robots.txt"), []byte(robotsTxt), 0644)

	// Generate sitemap.xml
	sitemapPages := []string{"/", "/about/", "/methodology/", "/privacy/", "/disclaimer/", "/responsible-gambling/", "/contact/"}
	sitemapXML := schema.GenerateSitemapIndex(site.Domain, sitemapPages)
	os.WriteFile(filepath.Join(outputDir, "sitemap.xml"), []byte(sitemapXML), 0644)

	// 6b. Write IndexNow key file (for Bing verification)
	var indexNowKey string
	b.DB.QueryRow("SELECT api_key FROM indexnow_keys WHERE domain_id=?", domainID).Scan(&indexNowKey)
	if indexNowKey != "" {
		os.WriteFile(filepath.Join(outputDir, indexNowKey+".txt"), []byte(indexNowKey), 0644)
	}

	// 7. Update domain status
	b.DB.Exec("UPDATE domains SET status='built', updated_at=CURRENT_TIMESTAMP WHERE id=?", domainID)

	duration := time.Since(start)
	b.finishBuildLog(logID, "success", string(output), duration)
	log.Printf("[BUILD] %s completed in %v", site.Domain, duration)
	return nil
}

func (b *Builder) generateConfig(buildDir string, site *siteInfo) error {
	cssPrefix := variation.GetCSSPrefix(site.Domain)
	trackFn := variation.GetTrackFnName(site.Domain)
	cssNoise := variation.GenerateCSSNoise(site.Domain)
	brandColorDark := variation.DarkenColor(site.BrandColor)
	publishDate := variation.RandomPublishDate(site.Domain)
	buildID := fmt.Sprintf("%d", time.Now().Unix())

	trackerURL := ""
	b.DB.QueryRow("SELECT value FROM settings WHERE key='tracker_url'").Scan(&trackerURL)

	// load title patterns for keyword type
	titlePatterns := b.loadTitlePatterns(site.KeywordType)
	pageTitle := site.PageTitle
	if pageTitle == "" && len(titlePatterns) > 0 {
		pageTitle = variation.PickTitle(titlePatterns, site.Domain, site.PrimaryKeyword)
	}
	if pageTitle == "" {
		pageTitle = site.PrimaryKeyword
	}

	// use content author if set, else pick from pool
	authorName := site.AuthorName
	authorBio := site.AuthorBio
	authorTitle := site.AuthorTitle
	if authorName == "" {
		authorsFile := filepath.Join(b.Cfg.DataDir, "authors", "authors.json")
		if data, err := os.ReadFile(authorsFile); err == nil {
			authors := variation.LoadAuthors(data)
			author := variation.PickAuthor(authors, site.Domain)
			authorName = author.Name
			authorBio = author.Bio
		}
	}
	if authorName == "" {
		authorName = "Editor"
		authorBio = "Content Editor"
		authorTitle = "Content Editor"
	}

	// generate schema JSON
	var extraFields map[string]interface{}
	json.Unmarshal([]byte(site.ExtraData), &extraFields)

	schemaJSON := schema.GenerateSchema(schema.SiteData{
		Domain:         site.Domain,
		SiteName:       site.BrandName,
		MetaTitle:      pageTitle,
		MetaDesc:       site.MetaDesc,
		MainContent:    site.BodyContent,
		BrandColor:     site.BrandColor,
		RedirectURL:    site.RedirectURL,
		PrimaryKeyword: site.PrimaryKeyword,
		KeywordType:    site.KeywordType,
		FAQJson:        site.FAQItems,
		ExtraFields:    extraFields,
	})

	// parse extra fields for template params
	extraParams := ""
	if extraFields != nil {
		for _, k := range variation.SortedKeys(extraFields) {
			v := extraFields[k]
			vJSON, _ := json.Marshal(v)
			extraParams += fmt.Sprintf("  %s: %s\n", k, string(vJSON))
		}
	}

	config := fmt.Sprintf(`baseURL: "https://%s/"
languageCode: "%s"
title: "%s"
theme: ""

params:
  brand_name: "%s"
  brand_color: "%s"
  brand_color_dark: "%s"
  page_title: "%s"
  meta_desc: "%s"
  h1: "%s"
  hero_title: "%s"
  hero_subtitle: "%s"
  cta_text: "%s"
  cta_sub: "%s"
  intro_text: "%s"
  body_content: |
    %s
  conclusion: "%s"
  target_keyword: "%s"
  primary_keyword: "%s"
  keyword_type: "%s"
  redirect_url: "%s"
  market: "%s"
  css_prefix: "%s"
  track_fn_name: "%s"
  css_noise: "%s"
  tracker_url: "%s"
  site_id: "%s"
  build_id: "%s"
  publish_date: "%s"
  site_year: "%d"
  author_name: "%s"
  author_title: "%s"
  author_bio: "%s"
  faq_title: "%s"
  trust_badges: '%s'
  disclosure: "%s"
  disclaimer: "%s"
  feature_1_icon: "%s"
  feature_1_title: "%s"
  feature_1_desc: "%s"
  feature_2_icon: "%s"
  feature_2_title: "%s"
  feature_2_desc: "%s"
  feature_3_icon: "%s"
  feature_3_title: "%s"
  feature_3_desc: "%s"
  schema_json: |
    %s
  faq_items: '%s'
  extra_data: '%s'
%s
`,
		site.Domain, site.Language, escYAML(site.BrandName),
		escYAML(site.BrandName), site.BrandColor, brandColorDark,
		escYAML(pageTitle), escYAML(site.MetaDesc), escYAML(site.H1),
		escYAML(site.HeroTitle), escYAML(site.HeroSubtitle),
		escYAML(site.CTAText), escYAML(site.CTASub),
		escYAML(site.IntroText), indentYAML(site.BodyContent, "    "), escYAML(site.Conclusion),
		escYAML(site.PrimaryKeyword), escYAML(site.PrimaryKeyword),
		site.KeywordType, escYAML(site.RedirectURL),
		site.Market, cssPrefix, trackFn, escYAML(cssNoise),
		escYAML(trackerURL), escYAML(site.Domain),
		buildID, publishDate, time.Now().Year(),
		escYAML(authorName), escYAML(authorTitle), escYAML(authorBio),
		escYAML(site.FAQTitle),
		escYAML(site.TrustBadges), escYAML(site.Disclosure), escYAML(site.Disclaimer),
		escYAML(site.Feature1Icon), escYAML(site.Feature1Title), escYAML(site.Feature1Desc),
		escYAML(site.Feature2Icon), escYAML(site.Feature2Title), escYAML(site.Feature2Desc),
		escYAML(site.Feature3Icon), escYAML(site.Feature3Title), escYAML(site.Feature3Desc),
		strings.ReplaceAll(schemaJSON, "\n", "\n    "),
		escYAML(site.FAQItems), escYAML(site.ExtraData),
		extraParams,
	)

	return os.WriteFile(filepath.Join(buildDir, "hugo.yaml"), []byte(config), 0644)
}

func (b *Builder) generateSupportPages(buildDir string, site *siteInfo) error {
	pages := schema.GenerateSupportPages(site.BrandName, site.Domain)
	contentDir := filepath.Join(buildDir, "content")
	os.MkdirAll(contentDir, 0755)

	for name, content := range pages {
		path := filepath.Join(contentDir, name+".md")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	// main index page
	indexContent := fmt.Sprintf(`---
title: "%s"
description: "%s"
---
`, escYAML(site.PageTitle), escYAML(site.MetaDesc))
	return os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(indexContent), 0644)
}

func (b *Builder) loadSiteInfo(domainID int64) (*siteInfo, error) {
	s := &siteInfo{ID: domainID}
	err := b.DB.QueryRow(`SELECT d.domain, d.template_id, d.market, d.language, d.keyword_type, d.primary_keyword, d.redirect_url,
		COALESCE(t.slug,'default') as tpl_slug, COALESCE(t.path,'') as tpl_path
		FROM domains d LEFT JOIN templates t ON d.template_id=t.id WHERE d.id=?`, domainID).
		Scan(&s.Domain, &s.TemplateID, &s.Market, &s.Language, &s.KeywordType, &s.PrimaryKeyword, &s.RedirectURL,
			&s.TemplateSlug, &s.TemplatePath)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	// load content - columns match db.go contents table
	b.DB.QueryRow(`SELECT
		COALESCE(brand_name,''), COALESCE(brand_color,'#1976D2'),
		COALESCE(cta_text,''), COALESCE(cta_sub,''),
		COALESCE(hero_title,''), COALESCE(hero_subtitle,''),
		COALESCE(intro_text,''), COALESCE(body_content,''), COALESCE(conclusion,''),
		COALESCE(page_title,''), COALESCE(meta_desc,''), COALESCE(h1,''),
		COALESCE(faq_title,'常見問題'), COALESCE(faq_items,'[]'), COALESCE(extra_data,'{}'),
		COALESCE(author_name,''), COALESCE(author_title,''), COALESCE(author_bio,''),
		COALESCE(trust_badges,'[]'), COALESCE(disclosure,''), COALESCE(disclaimer,''),
		COALESCE(feature_1_icon,''), COALESCE(feature_1_title,''), COALESCE(feature_1_desc,''),
		COALESCE(feature_2_icon,''), COALESCE(feature_2_title,''), COALESCE(feature_2_desc,''),
		COALESCE(feature_3_icon,''), COALESCE(feature_3_title,''), COALESCE(feature_3_desc,'')
		FROM contents WHERE domain_id=?`, domainID).
		Scan(&s.BrandName, &s.BrandColor,
			&s.CTAText, &s.CTASub,
			&s.HeroTitle, &s.HeroSubtitle,
			&s.IntroText, &s.BodyContent, &s.Conclusion,
			&s.PageTitle, &s.MetaDesc, &s.H1,
			&s.FAQTitle, &s.FAQItems, &s.ExtraData,
			&s.AuthorName, &s.AuthorTitle, &s.AuthorBio,
			&s.TrustBadges, &s.Disclosure, &s.Disclaimer,
			&s.Feature1Icon, &s.Feature1Title, &s.Feature1Desc,
			&s.Feature2Icon, &s.Feature2Title, &s.Feature2Desc,
			&s.Feature3Icon, &s.Feature3Title, &s.Feature3Desc)

	if s.BrandName == "" {
		s.BrandName = s.Domain
	}
	return s, nil
}

func (b *Builder) loadTitlePatterns(kwType string) []string {
	rows, err := b.DB.Query("SELECT pattern FROM title_variants WHERE keyword_type=? AND is_active=1", kwType)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var patterns []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		patterns = append(patterns, p)
	}
	return patterns
}

func (b *Builder) createBuildLog(domainID int64, status string) int64 {
	res, err := b.DB.Exec("INSERT INTO build_logs(domain_id, status) VALUES(?,?)", domainID, status)
	if err != nil {
		log.Printf("[BUILD] failed to create build log: %v", err)
		return 0
	}
	id, _ := res.LastInsertId()
	return id
}

func (b *Builder) finishBuildLog(logID int64, status, output string, dur time.Duration) {
	b.DB.Exec("UPDATE build_logs SET status=?, log_output=?, duration_ms=? WHERE id=?",
		status, output, dur.Milliseconds(), logID)
}

// indentYAML prepares content for YAML block scalar (|) by indenting each line
func indentYAML(s, prefix string) string {
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if i > 0 {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func escYAML(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func copyDir(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		// if template doesn't exist, create minimal structure
		os.MkdirAll(filepath.Join(dst, "layouts", "_default"), 0755)
		os.MkdirAll(filepath.Join(dst, "layouts", "partials"), 0755)
		os.MkdirAll(filepath.Join(dst, "static", "css"), 0755)
		os.MkdirAll(filepath.Join(dst, "content"), 0755)
		return nil
	}
	cmd := exec.Command("cp", "-r", src+"/.", dst)
	return cmd.Run()
}

// DeploySite copies built files to the target server
func (b *Builder) DeploySite(domainID int64) error {
	var domain, serverIP, serverUser, serverPath string
	err := b.DB.QueryRow("SELECT domain, server_ip, server_user, server_path FROM domains WHERE id=?", domainID).
		Scan(&domain, &serverIP, &serverUser, &serverPath)
	if err != nil {
		return fmt.Errorf("domain not found: %w", err)
	}

	outputDir := filepath.Join(b.Cfg.BuildOut, domain)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("site not built yet")
	}

	if serverIP == "" {
		// local deployment - already in place
		b.DB.Exec("UPDATE domains SET status='active', updated_at=CURRENT_TIMESTAMP WHERE id=?", domainID)
		return nil
	}

	if serverPath == "" {
		serverPath = "/var/www/sites/" + domain
	}
	if serverUser == "" {
		serverUser = "root"
	}

	// rsync to remote server
	target := fmt.Sprintf("%s@%s:%s/", serverUser, serverIP, serverPath)
	cmd := exec.Command("rsync", "-avz", "--delete", outputDir+"/", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("rsync: %s\n%w", string(output), err)
	}

	b.DB.Exec("UPDATE domains SET status='active', updated_at=CURRENT_TIMESTAMP WHERE id=?", domainID)
	log.Printf("[DEPLOY] %s -> %s", domain, target)
	return nil
}
