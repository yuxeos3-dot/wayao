package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// ========== IndexNow ==========

func (app *App) SubmitIndexNow(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		URLs []string `json:"urls"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	// ensure domain has an IndexNow key
	var apiKey string
	err = app.DB.QueryRow("SELECT api_key FROM indexnow_keys WHERE domain_id=?", id).Scan(&apiKey)
	if err != nil {
		b := make([]byte, 16)
		rand.Read(b)
		apiKey = hex.EncodeToString(b)
		app.DB.Exec("INSERT INTO indexnow_keys(domain_id, api_key) VALUES(?,?)", id, apiKey)
	}
	// get domain name
	var domain string
	app.DB.QueryRow("SELECT domain FROM domains WHERE id=?", id).Scan(&domain)

	for _, u := range req.URLs {
		app.DB.Exec("INSERT INTO index_submissions(domain_id, url, engine, status) VALUES(?,?,'indexnow','pending')", id, u)
	}

	// Actually send to IndexNow API
	if domain != "" && len(req.URLs) > 0 {
		go func() {
			payload := map[string]interface{}{
				"host":        domain,
				"key":         apiKey,
				"keyLocation": fmt.Sprintf("https://%s/%s.txt", domain, apiKey),
				"urlList":     req.URLs,
			}
			body, _ := json.Marshal(payload)
			resp, err := http.Post("https://api.indexnow.org/IndexNow", "application/json; charset=utf-8", bytes.NewReader(body))
			if err != nil {
				log.Printf("[IndexNow] send failed for %s: %v", domain, err)
				return
			}
			defer resp.Body.Close()
			status := "submitted"
			if resp.StatusCode >= 400 {
				status = fmt.Sprintf("error_%d", resp.StatusCode)
			}
			log.Printf("[IndexNow] %s: %d URLs sent, status=%s", domain, len(req.URLs), status)
			for _, u := range req.URLs {
				app.DB.Exec("UPDATE index_submissions SET status=? WHERE domain_id=? AND url=?", status, id, u)
			}
		}()
	}

	jsonOK(w, map[string]interface{}{"api_key": apiKey, "submitted": len(req.URLs)})
}

func (app *App) GetIndexNowRecords(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	rows, err := app.DB.Query("SELECT id, url, engine, status, created_at FROM index_submissions WHERE domain_id=? ORDER BY id DESC LIMIT 100", id)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var sid int64
		var url, engine, status, createdAt string
		rows.Scan(&sid, &url, &engine, &status, &createdAt)
		list = append(list, map[string]interface{}{
			"id": sid, "url": url, "engine": engine, "status": status, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

// ========== City Matrix ==========

func (app *App) GetCityMatrix(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	rows, err := app.DB.Query("SELECT id, city_name, city_slug, extra_title, extra_desc, is_built, created_at FROM city_matrix WHERE domain_id=? ORDER BY city_name", id)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var cid int64
		var name, slug, title, desc, createdAt string
		var isBuilt int
		rows.Scan(&cid, &name, &slug, &title, &desc, &isBuilt, &createdAt)
		list = append(list, map[string]interface{}{
			"id": cid, "city_name": name, "city_slug": slug,
			"extra_title": title, "extra_desc": desc, "is_built": isBuilt == 1, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) SaveCityMatrix(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		Cities []struct {
			CityName   string `json:"city_name"`
			CitySlug   string `json:"city_slug"`
			ExtraTitle string `json:"extra_title"`
			ExtraDesc  string `json:"extra_desc"`
		} `json:"cities"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	app.DB.Exec("DELETE FROM city_matrix WHERE domain_id=?", id)
	for _, c := range req.Cities {
		app.DB.Exec("INSERT INTO city_matrix(domain_id, city_name, city_slug, extra_title, extra_desc) VALUES(?,?,?,?,?)",
			id, c.CityName, c.CitySlug, c.ExtraTitle, c.ExtraDesc)
	}
	jsonOK(w, map[string]interface{}{"count": len(req.Cities)})
}

// ========== Title Pool ==========

func (app *App) ListTitlePool(w http.ResponseWriter, r *http.Request) {
	kwType := getQuery(r, "keyword_type", "")
	where := "1=1"
	args := []interface{}{}
	if kwType != "" {
		where += " AND keyword_type=?"
		args = append(args, kwType)
	}
	rows, err := app.DB.Query(fmt.Sprintf("SELECT id, keyword_type, slot, template, is_active FROM title_pool WHERE %s ORDER BY keyword_type, id", where), args...)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var tid int64
		var kt, slot, tmpl string
		var isActive int
		rows.Scan(&tid, &kt, &slot, &tmpl, &isActive)
		list = append(list, map[string]interface{}{
			"id": tid, "keyword_type": kt, "slot": slot, "template": tmpl, "is_active": isActive == 1,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) CreateTitleVariant(w http.ResponseWriter, r *http.Request) {
	var req struct {
		KeywordType string `json:"keyword_type"`
		Slot        string `json:"slot"`
		Template    string `json:"template"`
		Market      string `json:"market"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.KeywordType == "" || req.Template == "" {
		jsonError(w, 400, "keyword_type and template required")
		return
	}
	if req.Slot == "" {
		req.Slot = "title"
	}
	if req.Market == "" {
		req.Market = "zh-TW"
	}
	res, err := app.DB.Exec("INSERT INTO title_pool(keyword_type, slot, template, market) VALUES(?,?,?,?)", req.KeywordType, req.Slot, req.Template, req.Market)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	jsonOK(w, map[string]interface{}{"id": id})
}

func (app *App) DeleteTitleVariant(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	app.DB.Exec("DELETE FROM title_pool WHERE id=?", id)
	jsonOK(w, "deleted")
}

// ========== Site Clusters ==========

func (app *App) ListClusters(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(`SELECT c.id, c.name, c.slug, c.description, c.created_at, COUNT(m.id) as member_count
		FROM site_clusters c LEFT JOIN site_cluster_members m ON c.id=m.cluster_id
		GROUP BY c.id ORDER BY c.id DESC`)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var cid int64
		var name, slug, desc, createdAt string
		var cnt int
		rows.Scan(&cid, &name, &slug, &desc, &createdAt, &cnt)
		list = append(list, map[string]interface{}{
			"id": cid, "name": name, "slug": slug, "description": desc,
			"created_at": createdAt, "member_count": cnt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) CreateCluster(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Name == "" {
		jsonError(w, 400, "name required")
		return
	}
	if req.Slug == "" {
		req.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}
	res, err := app.DB.Exec("INSERT INTO site_clusters(name, slug, description) VALUES(?,?,?)", req.Name, req.Slug, req.Description)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	jsonOK(w, map[string]interface{}{"id": id})
}

func (app *App) DeleteCluster(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	app.DB.Exec("DELETE FROM site_clusters WHERE id=?", id)
	jsonOK(w, "deleted")
}

func (app *App) AddClusterMember(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		DomainID int64  `json:"domain_id"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Role == "" {
		req.Role = "member"
	}
	app.DB.Exec("INSERT OR IGNORE INTO site_cluster_members(cluster_id, domain_id, role) VALUES(?,?,?)", id, req.DomainID, req.Role)
	jsonOK(w, "added")
}

func (app *App) RemoveClusterMember(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		DomainID int64 `json:"domain_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	app.DB.Exec("DELETE FROM site_cluster_members WHERE cluster_id=? AND domain_id=?", id, req.DomainID)
	jsonOK(w, "removed")
}

// ========== Content Refresh ==========

func (app *App) ListRefreshSchedule(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(`SELECT s.id, s.domain_id, d.domain, s.refresh_type, s.frequency_days, s.last_refreshed, s.next_refresh, s.is_active
		FROM content_refresh_schedule s JOIN domains d ON s.domain_id=d.id ORDER BY s.id`)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var sid, did int64
		var domain, refreshType string
		var freqDays, isActive int
		var lastRefreshed, nextRefresh *string
		rows.Scan(&sid, &did, &domain, &refreshType, &freqDays, &lastRefreshed, &nextRefresh, &isActive)
		list = append(list, map[string]interface{}{
			"id": sid, "domain_id": did, "domain": domain, "refresh_type": refreshType,
			"frequency_days": freqDays, "last_refreshed": lastRefreshed,
			"next_refresh": nextRefresh, "is_active": isActive == 1,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) SaveRefreshSchedule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DomainID      int64  `json:"domain_id"`
		RefreshType   string `json:"refresh_type"`
		FrequencyDays int    `json:"frequency_days"`
		IsActive      bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.RefreshType == "" {
		req.RefreshType = "timestamp"
	}
	if req.FrequencyDays < 1 {
		req.FrequencyDays = 14
	}
	active := 0
	if req.IsActive {
		active = 1
	}
	app.DB.Exec(`INSERT INTO content_refresh_schedule(domain_id, refresh_type, frequency_days, is_active, next_refresh)
		VALUES(?,?,?,?, datetime('now', '+' || ? || ' days'))
		ON CONFLICT(domain_id) DO UPDATE SET refresh_type=excluded.refresh_type, frequency_days=excluded.frequency_days, is_active=excluded.is_active`,
		req.DomainID, req.RefreshType, req.FrequencyDays, active, req.FrequencyDays)
	jsonOK(w, "saved")
}

// ========== Health Check ==========

func (app *App) CheckSiteHealth(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var domain, status string
	err = app.DB.QueryRow("SELECT domain, status FROM domains WHERE id=?", id).Scan(&domain, &status)
	if err != nil {
		jsonError(w, 404, "domain not found")
		return
	}

	// real HTTP health check
	checks := map[string]interface{}{}
	score := 100
	var issues []string

	// HTTP check
	client := &http.Client{Timeout: 10 * time.Second}
	start := time.Now()
	resp, httpErr := client.Get("https://" + domain + "/")
	ttfb := int(time.Since(start).Milliseconds())

	if httpErr != nil {
		checks["http"] = "error"
		checks["http_error"] = httpErr.Error()
		score -= 40
		issues = append(issues, "HTTP unreachable: "+httpErr.Error())
	} else {
		resp.Body.Close()
		checks["http"] = resp.StatusCode
		checks["ttfb_ms"] = ttfb
		if resp.StatusCode != 200 {
			score -= 20
			issues = append(issues, fmt.Sprintf("HTTP %d (expected 200)", resp.StatusCode))
		}
		if ttfb > 3000 {
			score -= 10
			issues = append(issues, fmt.Sprintf("Slow TTFB: %dms", ttfb))
		}
		// SSL check (if HTTPS succeeded, SSL is OK)
		if resp.TLS != nil {
			checks["ssl"] = "valid"
		} else {
			checks["ssl"] = "missing"
			score -= 15
			issues = append(issues, "No SSL certificate")
		}
	}

	if score < 0 {
		score = 0
	}

	jsonOK(w, map[string]interface{}{
		"domain": domain, "status": status, "score": score,
		"checks": checks, "issues": issues, "ttfb_ms": ttfb,
	})
}

func (app *App) GetHealthAlerts(w http.ResponseWriter, r *http.Request) {
	// return domains with issues
	rows, err := app.DB.Query("SELECT id, domain, status FROM domains WHERE status IN ('error','down') ORDER BY updated_at DESC")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id int64
		var domain, status string
		rows.Scan(&id, &domain, &status)
		list = append(list, map[string]interface{}{"id": id, "domain": domain, "status": status})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

// ========== IP / UA Rules ==========

func (app *App) ListIPRules(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT id, ip_cidr, rule_type, reason, created_at FROM ip_rules ORDER BY id DESC")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id int64
		var ip, action, reason, createdAt string
		rows.Scan(&id, &ip, &action, &reason, &createdAt)
		list = append(list, map[string]interface{}{
			"id": id, "ip": ip, "action": action, "reason": reason, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) AddIPRule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IP     string `json:"ip"`
		Action string `json:"action"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Action == "" {
		req.Action = "block"
	}
	app.DB.Exec("INSERT INTO ip_rules(ip_cidr, rule_type, reason) VALUES(?,?,?)", req.IP, req.Action, req.Reason)
	jsonOK(w, "added")
}

func (app *App) DeleteIPRule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	app.DB.Exec("DELETE FROM ip_rules WHERE id=?", id)
	jsonOK(w, "deleted")
}

func (app *App) ListUARules(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT id, pattern, rule_type, is_active, created_at FROM ua_rules ORDER BY id DESC")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var id int64
		var pattern, ruleType, createdAt string
		var isActive int
		rows.Scan(&id, &pattern, &ruleType, &isActive, &createdAt)
		list = append(list, map[string]interface{}{
			"id": id, "pattern": pattern, "rule_type": ruleType, "is_active": isActive == 1, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) DeleteUARule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	app.DB.Exec("DELETE FROM ua_rules WHERE id=?", id)
	jsonOK(w, "deleted")
}

func (app *App) AddUARule(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pattern string `json:"pattern"`
		Action  string `json:"action"`
		Reason  string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Action == "" {
		req.Action = "block"
	}
	app.DB.Exec("INSERT INTO ua_rules(pattern, rule_type) VALUES(?,?)", req.Pattern, req.Action)
	jsonOK(w, "added")
}

// ========== Export / Import ==========

func (app *App) ExportDomain(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	// export domain + content as JSON
	var domain, market, lang, kwType, pk, redirect, status string
	var templateID int64
	err = app.DB.QueryRow("SELECT domain, template_id, market, language, keyword_type, primary_keyword, redirect_url, status FROM domains WHERE id=?", id).
		Scan(&domain, &templateID, &market, &lang, &kwType, &pk, &redirect, &status)
	if err != nil {
		jsonError(w, 404, "domain not found")
		return
	}

	var content map[string]interface{}
	var brandName, brandColor, heroTitle, heroSubtitle, bodyContent, pageTitle, metaDesc, faqItems, extraData string
	err = app.DB.QueryRow("SELECT brand_name, brand_color, hero_title, hero_subtitle, body_content, page_title, meta_desc, faq_items, extra_data FROM contents WHERE domain_id=?", id).
		Scan(&brandName, &brandColor, &heroTitle, &heroSubtitle, &bodyContent, &pageTitle, &metaDesc, &faqItems, &extraData)
	if err == nil {
		content = map[string]interface{}{
			"brand_name": brandName, "brand_color": brandColor,
			"hero_title": heroTitle, "hero_subtitle": heroSubtitle, "body_content": bodyContent,
			"page_title": pageTitle, "meta_desc": metaDesc,
			"faq_items": faqItems, "extra_data": extraData,
		}
	}

	export := map[string]interface{}{
		"domain": domain, "template_id": templateID, "market": market,
		"language": lang, "keyword_type": kwType, "primary_keyword": pk,
		"redirect_url": redirect, "status": status, "content": content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", domain))
	json.NewEncoder(w).Encode(export)
}
