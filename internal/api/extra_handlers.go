package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// HandleLogout POST /api/v1/auth/logout
func (app *App) HandleLogout(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"status": "ok"})
}

// HandleOGImage GET /og/{site_id}.svg
func (app *App) HandleOGImage(w http.ResponseWriter, r *http.Request) {
	siteID := r.URL.Path[len("/og/"):]
	siteID = strings.TrimSuffix(siteID, ".svg")

	var domain, brandName string
	app.DB.QueryRow("SELECT d.domain, COALESCE(c.brand_name, d.domain) FROM domains d LEFT JOIN contents c ON c.domain_id=d.id WHERE d.site_id=? OR d.domain=?", siteID, siteID).Scan(&domain, &brandName)
	if domain == "" {
		domain = siteID
		brandName = siteID
	}

	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="630" viewBox="0 0 1200 630">
  <rect width="1200" height="630" fill="#1a1a2e"/>
  <text x="600" y="280" text-anchor="middle" fill="#fff" font-size="48" font-family="Arial, sans-serif" font-weight="bold">%s</text>
  <text x="600" y="350" text-anchor="middle" fill="#aaa" font-size="24" font-family="Arial, sans-serif">%s</text>
</svg>`, brandName, domain)

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write([]byte(svg))
}

// HandleBuildAndDeploy POST /api/v1/build/{id}/full
func (app *App) HandleBuildAndDeploy(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	if err := app.BuildSiteByID(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	if err := app.DeploySiteByID(id); err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	jsonOK(w, "build and deploy complete")
}

// GetBuildStatus GET /api/v1/build/{id}/status
func (app *App) GetBuildStatus(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var status, buildError string
	var lastBuilt, lastDeployed *string
	app.DB.QueryRow("SELECT status, build_error, last_built_at, last_deployed_at FROM domains WHERE id=?", id).
		Scan(&status, &buildError, &lastBuilt, &lastDeployed)
	jsonOK(w, map[string]interface{}{
		"status": status, "build_error": buildError,
		"last_built_at": lastBuilt, "last_deployed_at": lastDeployed,
	})
}

// HandleBatchBuild POST /api/v1/build/batch
func (app *App) HandleBatchBuild(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs    []int64 `json:"ids"`
		Action string  `json:"action"` // build, deploy, full
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Action == "" {
		req.Action = "full"
	}
	go func() {
		sem := make(chan struct{}, 3) // max 3 concurrent builds
		var wg sync.WaitGroup
		for _, id := range req.IDs {
			id := id
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				switch req.Action {
				case "build":
					app.BuildSiteByID(id)
				case "deploy":
					app.DeploySiteByID(id)
				case "full":
					if app.BuildSiteByID(id) == nil {
						app.DeploySiteByID(id)
					}
				}
			}()
		}
		wg.Wait()
	}()
	jsonOK(w, map[string]interface{}{"status": "started", "count": len(req.IDs)})
}

// GetDomainStats GET /api/v1/stats/domain/{id}
func (app *App) GetDomainStats(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var siteID string
	app.DB.QueryRow("SELECT site_id FROM domains WHERE id=?", id).Scan(&siteID)

	var totalPV, totalClick, uniqueIP int
	app.DB.QueryRow("SELECT COUNT(*) FROM pageviews WHERE site_id=?", siteID).Scan(&totalPV)
	app.DB.QueryRow("SELECT COUNT(*) FROM clicks WHERE site_id=?", siteID).Scan(&totalClick)
	app.DB.QueryRow("SELECT COUNT(DISTINCT ip) FROM clicks WHERE site_id=?", siteID).Scan(&uniqueIP)

	// by action
	rows, _ := app.DB.Query("SELECT action, COUNT(*) FROM clicks WHERE site_id=? GROUP BY action", siteID)
	actions := map[string]int{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var a string
			var c int
			rows.Scan(&a, &c)
			actions[a] = c
		}
	}

	jsonOK(w, map[string]interface{}{
		"total_pv": totalPV, "total_click": totalClick, "unique_ip": uniqueIP,
		"actions": actions,
	})
}

// GetDailySummary GET /api/v1/stats/summary
func (app *App) GetDailySummary(w http.ResponseWriter, r *http.Request) {
	siteID := getQuery(r, "site_id", "")
	from := getQuery(r, "from", "")
	to := getQuery(r, "to", "")

	where := "1=1"
	args := []interface{}{}
	if siteID != "" {
		where += " AND site_id=?"
		args = append(args, siteID)
	}
	if from != "" {
		where += " AND created_at >= ?"
		args = append(args, from)
	}
	if to != "" {
		where += " AND created_at <= ?"
		args = append(args, to+" 23:59:59")
	}

	rows, _ := app.DB.Query(
		"SELECT date(created_at), site_id, COUNT(*) FROM clicks WHERE "+where+" GROUP BY date(created_at), site_id ORDER BY 1 DESC", args...)
	var summary []map[string]interface{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var date, sid string
			var cnt int
			rows.Scan(&date, &sid, &cnt)
			summary = append(summary, map[string]interface{}{"date": date, "site_id": sid, "clicks": cnt})
		}
	}
	if summary == nil {
		summary = []map[string]interface{}{}
	}
	jsonOK(w, map[string]interface{}{"summary": summary})
}

// ListRankings GET /api/v1/rankings
func (app *App) ListRankings(w http.ResponseWriter, r *http.Request) {
	domainID := getQuery(r, "domain_id", "")
	limit := getQueryInt(r, "limit", 100)

	where := "1=1"
	args := []interface{}{}
	if domainID != "" {
		where += " AND r.domain_id=?"
		args = append(args, domainID)
	}
	args = append(args, limit)

	rows, err := app.DB.Query(fmt.Sprintf(
		"SELECT r.id, r.domain_id, r.keyword, r.rank, r.market, r.engine, r.checked_at FROM ranking_history r WHERE %s ORDER BY r.checked_at DESC LIMIT ?", where), args...)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, domainID int64
		var keyword, market, engine, checkedAt string
		var rank int
		rows.Scan(&id, &domainID, &keyword, &rank, &market, &engine, &checkedAt)
		list = append(list, map[string]interface{}{
			"id": id, "domain_id": domainID, "keyword": keyword,
			"rank": rank, "market": market, "engine": engine, "checked_at": checkedAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

// CheckRanking POST /api/v1/rankings/check/{id}
func (app *App) CheckRanking(w http.ResponseWriter, r *http.Request) {
	id, _ := parseID(r)
	jsonOK(w, map[string]interface{}{
		"status":  "use scripts/check_rankings.py for batch ranking checks",
		"domain_id": id,
	})
}

// BatchHealthCheck POST /api/v1/health-check/batch
func (app *App) BatchHealthCheck(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT id, domain, status FROM domains WHERE status IN ('active','built') ORDER BY id")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	type domainInfo struct {
		ID     int64
		Domain string
		Status string
	}
	var domains []domainInfo
	for rows.Next() {
		var d domainInfo
		rows.Scan(&d.ID, &d.Domain, &d.Status)
		domains = append(domains, d)
	}

	// real HTTP check with concurrency limit
	client := &http.Client{Timeout: 10 * time.Second}
	sem := make(chan struct{}, 5)
	var mu sync.Mutex
	var results []map[string]interface{}

	var wg sync.WaitGroup
	for _, d := range domains {
		d := d
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			score := 100
			var issues []string
			ttfb := 0

			start := time.Now()
			resp, httpErr := client.Get("https://" + d.Domain + "/")
			ttfb = int(time.Since(start).Milliseconds())

			if httpErr != nil {
				score -= 40
				issues = append(issues, "HTTP error: "+httpErr.Error())
			} else {
				resp.Body.Close()
				if resp.StatusCode != 200 {
					score -= 20
					issues = append(issues, fmt.Sprintf("HTTP %d", resp.StatusCode))
				}
				if ttfb > 3000 {
					score -= 10
					issues = append(issues, fmt.Sprintf("Slow: %dms", ttfb))
				}
			}
			if score < 0 {
				score = 0
			}

			mu.Lock()
			results = append(results, map[string]interface{}{
				"id": d.ID, "domain": d.Domain, "status": d.Status,
				"score": score, "ttfb_ms": ttfb, "issues": issues,
			})
			mu.Unlock()
		}()
	}
	wg.Wait()

	if results == nil {
		results = []map[string]interface{}{}
	}
	jsonOK(w, results)
}

// CheckIndexStatus GET /api/v1/index-status/{id}
func (app *App) CheckIndexStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := parseID(r)
	var domain string
	app.DB.QueryRow("SELECT domain FROM domains WHERE id=?", id).Scan(&domain)
	jsonOK(w, map[string]interface{}{
		"domain": domain,
		"status": "check via Google Search Console or SerpAPI",
		"tip":    "use scripts/check_rankings.py --domain " + domain,
	})
}

// BatchCheckIndex POST /api/v1/index-status/batch
func (app *App) BatchCheckIndex(w http.ResponseWriter, r *http.Request) {
	jsonOK(w, map[string]string{"status": "use scripts/check_rankings.py for batch checks"})
}

// GetCTRScore POST /api/v1/ctr-score
func (app *App) GetCTRScore(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title   string `json:"title"`
		Desc    string `json:"desc"`
		Keyword string `json:"keyword"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	score := 50
	// title length check
	titleLen := len([]rune(req.Title))
	if titleLen >= 30 && titleLen <= 60 {
		score += 10
	}
	// contains keyword
	if strings.Contains(req.Title, req.Keyword) {
		score += 15
	}
	// desc length
	descLen := len([]rune(req.Desc))
	if descLen >= 100 && descLen <= 160 {
		score += 10
	}
	// contains numbers/symbols for CTR
	for _, c := range req.Title {
		if c >= '0' && c <= '9' || c == '✓' || c == '★' || c == '|' {
			score += 5
			break
		}
	}
	if score > 100 {
		score = 100
	}

	jsonOK(w, map[string]interface{}{
		"score": score,
		"tips":  []string{"標題包含關鍵詞", "描述控制在120-160字", "使用數字和符號提高CTR"},
	})
}

// ExportBatch POST /api/v1/export/batch
func (app *App) ExportBatch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	var exports []map[string]interface{}
	for _, id := range req.IDs {
		var domain, market, kwType string
		app.DB.QueryRow("SELECT domain, market, keyword_type FROM domains WHERE id=?", id).Scan(&domain, &market, &kwType)
		if domain != "" {
			exports = append(exports, map[string]interface{}{"id": id, "domain": domain, "market": market, "keyword_type": kwType})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=batch-export.json")
	json.NewEncoder(w).Encode(exports)
}

// ImportDomain POST /api/v1/import
func (app *App) ImportDomain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain      string                 `json:"domain"`
		Market      string                 `json:"market"`
		KeywordType string                 `json:"keyword_type"`
		RedirectURL string                 `json:"redirect_url"`
		SiteID      string                 `json:"site_id"`
		Content     map[string]interface{} `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Domain == "" {
		jsonError(w, 400, "domain required")
		return
	}
	if req.Market == "" {
		req.Market = "zh-TW"
	}
	if req.SiteID == "" {
		req.SiteID = strings.ReplaceAll(req.Domain, ".", "_")
	}

	res, err := app.DB.Exec("INSERT INTO domains(domain, market, keyword_type, redirect_url, site_id) VALUES(?,?,?,?,?)",
		req.Domain, req.Market, req.KeywordType, req.RedirectURL, req.SiteID)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()

	if req.Content != nil {
		// whitelist must match db.go contents table columns
		allowed := map[string]bool{
			"keyword_type": true, "target_keyword": true, "page_title": true, "meta_desc": true,
			"h1": true, "brand_name": true, "brand_color": true, "cta_text": true, "cta_sub": true,
			"hero_title": true, "hero_subtitle": true, "intro_text": true, "body_content": true,
			"conclusion": true, "faq_title": true, "faq_items": true, "extra_data": true,
			"author_name": true, "author_title": true, "author_bio": true,
			"trust_badges": true, "disclosure": true, "disclaimer": true,
			"feature_1_icon": true, "feature_1_title": true, "feature_1_desc": true,
			"feature_2_icon": true, "feature_2_title": true, "feature_2_desc": true,
			"feature_3_icon": true, "feature_3_title": true, "feature_3_desc": true,
		}
		app.DB.Exec("INSERT OR IGNORE INTO contents(domain_id) VALUES(?)", id)
		for k, v := range req.Content {
			if !allowed[k] {
				continue
			}
			app.DB.Exec(fmt.Sprintf("UPDATE contents SET %s=? WHERE domain_id=?", k), v, id)
		}
	}

	jsonOK(w, map[string]interface{}{"id": id})
}

// RunRefreshNow POST /api/v1/refresh/{id}/run
func (app *App) RunRefreshNow(w http.ResponseWriter, r *http.Request) {
	id, _ := parseID(r)
	// update last_updated timestamp in contents
	app.DB.Exec("UPDATE contents SET last_updated=?, last_updated_iso=?, updated_at=CURRENT_TIMESTAMP WHERE domain_id=?",
		time.Now().Format("2006年1月"), time.Now().Format("2006-01-02"), id)
	app.DB.Exec("UPDATE content_refresh_schedule SET last_refreshed=CURRENT_TIMESTAMP WHERE domain_id=?", id)
	jsonOK(w, "refreshed")
}

// StartContentRefreshScheduler runs in background
func (app *App) StartContentRefreshScheduler() {
	for {
		time.Sleep(1 * time.Hour)
		enabled := getSetting(app.DB, "content_refresh_on")
		if enabled != "1" {
			continue
		}
		rows, err := app.DB.Query("SELECT domain_id FROM content_refresh_schedule WHERE is_active=1 AND next_refresh <= CURRENT_TIMESTAMP")
		if err != nil {
			continue
		}
		for rows.Next() {
			var domainID int64
			rows.Scan(&domainID)
			// update timestamps
			app.DB.Exec("UPDATE contents SET last_updated=?, last_updated_iso=?, updated_at=CURRENT_TIMESTAMP WHERE domain_id=?",
				time.Now().Format("2006年1月"), time.Now().Format("2006-01-02"), domainID)
			app.DB.Exec("UPDATE content_refresh_schedule SET last_refreshed=CURRENT_TIMESTAMP, next_refresh=datetime('now', '+' || frequency_days || ' days') WHERE domain_id=?", domainID)
			// trigger rebuild so HTML reflects new timestamps
			if app.BuildFunc != nil {
				app.BuildFunc(domainID)
			}
		}
		rows.Close()
	}
}

// StartHealthCheckScheduler runs in background
func (app *App) StartHealthCheckScheduler() {
	client := &http.Client{Timeout: 10 * time.Second}
	for {
		time.Sleep(6 * time.Hour)
		rows, err := app.DB.Query("SELECT id, domain FROM domains WHERE status='active'")
		if err != nil {
			continue
		}
		for rows.Next() {
			var id int64
			var domain string
			rows.Scan(&id, &domain)
			// real HTTP check
			resp, err := client.Get("https://" + domain + "/")
			if err != nil {
				app.DB.Exec("UPDATE domains SET status='error', build_error=? WHERE id=?", err.Error(), id)
			} else {
				resp.Body.Close()
				if resp.StatusCode >= 400 {
					app.DB.Exec("UPDATE domains SET status='error', build_error=? WHERE id=?",
						fmt.Sprintf("HTTP %d", resp.StatusCode), id)
				}
			}
			time.Sleep(2 * time.Second) // rate limit checks
		}
		rows.Close()
	}
}
