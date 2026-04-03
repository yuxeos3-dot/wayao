package variation

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// GetCSSPrefix generates a unique 4-char CSS class prefix from domain
func GetCSSPrefix(domain string) string {
	h := md5.Sum([]byte(domain))
	return fmt.Sprintf("%x", h[:2])
}

// GetTrackFnName generates a unique tracker function name from domain
func GetTrackFnName(domain string) string {
	h := fnv.New32a()
	h.Write([]byte(domain + "fn"))
	return fmt.Sprintf("_t%x", h.Sum32()&0xFFFF)
}

// DarkenColor takes a hex color and returns a darker variant
func DarkenColor(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "#0D47A1"
	}
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	r = int(float64(r) * 0.7)
	g = int(float64(g) * 0.7)
	b = int(float64(b) * 0.7)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// ShuffleJSONFields actually randomizes JSON key order using domain as seed
func ShuffleJSONFields(jsonStr string, domain string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		return jsonStr
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// deterministic shuffle based on domain
	h := fnv.New32a()
	h.Write([]byte(domain))
	rng := rand.New(rand.NewSource(int64(h.Sum32())))
	rng.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	// build ordered JSON manually
	var sb strings.Builder
	sb.WriteString("{")
	for i, k := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		kJSON, _ := json.Marshal(k)
		vJSON, _ := json.Marshal(m[k])
		sb.Write(kJSON)
		sb.WriteString(":")
		sb.Write(vJSON)
	}
	sb.WriteString("}")
	return sb.String()
}

// GenerateCSSNoise injects subtle CSS variations
func GenerateCSSNoise(domain string) string {
	h := fnv.New32a()
	h.Write([]byte(domain))
	seed := int64(h.Sum32())
	rng := rand.New(rand.NewSource(seed))

	fontSize := 15.5 + rng.Float64()*0.3 // 15.5-15.8px
	letterSpacing := rng.Float64() * 0.02 // 0-0.02em
	lineHeight := 1.55 + rng.Float64()*0.1 // 1.55-1.65

	tag := domain
	if len(tag) > 4 {
		tag = tag[:4]
	}
	return fmt.Sprintf(`/* v:%s */
html { font-size: %.1fpx; letter-spacing: %.4fem; line-height: %.2f; }`,
		tag, fontSize, letterSpacing, lineHeight)
}

// PickTitle selects a title pattern from the pool using domain as seed
func PickTitle(patterns []string, domain, keyword string) string {
	if len(patterns) == 0 {
		return keyword
	}
	h := fnv.New32a()
	h.Write([]byte(domain + keyword))
	idx := int(h.Sum32()) % len(patterns)
	title := patterns[idx]
	title = strings.ReplaceAll(title, "{keyword}", keyword)
	title = strings.ReplaceAll(title, "{year}", fmt.Sprintf("%d", time.Now().Year()))
	return title
}

// RandomPublishDate generates a past date that looks natural
func RandomPublishDate(domain string) string {
	h := fnv.New32a()
	h.Write([]byte(domain + "pub"))
	seed := int64(h.Sum32())
	rng := rand.New(rand.NewSource(seed))
	daysAgo := 30 + rng.Intn(180) // 1-7 months ago
	t := time.Now().AddDate(0, 0, -daysAgo)
	return t.Format("2006-01-02")
}

// PickAuthor selects from author pool
func PickAuthor(authors []Author, domain string) Author {
	if len(authors) == 0 {
		return Author{Name: "Editor", Bio: "Content Editor"}
	}
	h := fnv.New32a()
	h.Write([]byte(domain + "author"))
	idx := int(h.Sum32()) % len(authors)
	return authors[idx]
}

type Author struct {
	Name   string `json:"name"`
	Bio    string `json:"bio"`
	Avatar string `json:"avatar"`
}

// LoadAuthors loads authors from a JSON file content
func LoadAuthors(data []byte) []Author {
	var authors []Author
	json.Unmarshal(data, &authors)
	return authors
}

// GetFontStack generates a unique font-family stack per domain (anti-fingerprint)
func GetFontStack(domain string) string {
	stacks := []string{
		`-apple-system,BlinkMacSystemFont,"Noto Sans TC",Roboto,sans-serif`,
		`"Noto Sans TC",-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif`,
		`-apple-system,"Noto Sans TC",BlinkMacSystemFont,Roboto,"Helvetica Neue",sans-serif`,
		`BlinkMacSystemFont,-apple-system,"Noto Sans TC","Segoe UI",Roboto,sans-serif`,
		`"Segoe UI",-apple-system,BlinkMacSystemFont,"Noto Sans TC",Roboto,sans-serif`,
		`Roboto,"Noto Sans TC",-apple-system,BlinkMacSystemFont,sans-serif`,
	}
	h := fnv.New32a()
	h.Write([]byte(domain + "font"))
	return stacks[int(h.Sum32())%len(stacks)]
}

// GenerateFaviconSVG creates a unique favicon per domain from hash
func GenerateFaviconSVG(domain, brandColor string) string {
	h := fnv.New32a()
	h.Write([]byte(domain + "fav"))
	seed := h.Sum32()

	if brandColor == "" {
		brandColor = "#1976D2"
	}
	// deterministic shape: circle, square, or rounded-rect with first letter
	letter := strings.ToUpper(domain[:1])
	shapes := []string{
		fmt.Sprintf(`<circle cx="16" cy="16" r="14" fill="%s"/>`, brandColor),
		fmt.Sprintf(`<rect x="2" y="2" width="28" height="28" rx="6" fill="%s"/>`, brandColor),
		fmt.Sprintf(`<rect x="2" y="2" width="28" height="28" rx="14" fill="%s"/>`, brandColor),
	}
	shape := shapes[seed%3]

	svg := fmt.Sprintf(`%%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%%3E%s%%3Ctext x='16' y='22' text-anchor='middle' fill='white' font-size='18' font-weight='bold' font-family='sans-serif'%%3E%s%%3C/text%%3E%%3C/svg%%3E`,
		strings.ReplaceAll(strings.ReplaceAll(shape, `"`, `'`), `<`, `%3C`),
		letter)
	// URL-encode the remaining < and >
	svg = strings.ReplaceAll(svg, ">", "%3E")
	return svg
}

// GenerateHreflangTags generates hreflang link tags for multi-market domains
func GenerateHreflangTags(domain, lang string, siblings map[string]string) string {
	if len(siblings) == 0 {
		return ""
	}
	var sb strings.Builder
	// self-referencing hreflang
	sb.WriteString(fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="https://%s/">`, lang, domain))
	sb.WriteString("\n")
	for hrefLang, hrefDomain := range siblings {
		if hrefDomain != domain {
			sb.WriteString(fmt.Sprintf(`<link rel="alternate" hreflang="%s" href="https://%s/">`, hrefLang, hrefDomain))
			sb.WriteString("\n")
		}
	}
	sb.WriteString(`<link rel="alternate" hreflang="x-default" href="https://` + domain + `/">`)
	return sb.String()
}

// SortedKeys returns map keys in sorted order (for deterministic output)
func SortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
