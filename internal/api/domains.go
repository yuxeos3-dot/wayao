package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func (app *App) ListDomains(w http.ResponseWriter, r *http.Request) {
	market := getQuery(r, "market", "")
	status := getQuery(r, "status", "")
	kwType := getQuery(r, "keyword_type", "")
	page := getQueryInt(r, "page", 1)
	size := cappedSize(r, 50, 500)
	offset := (page - 1) * size

	where := "1=1"
	args := []interface{}{}
	if market != "" {
		where += " AND d.market=?"
		args = append(args, market)
	}
	if status != "" {
		where += " AND d.status=?"
		args = append(args, status)
	}
	if kwType != "" {
		where += " AND d.keyword_type=?"
		args = append(args, kwType)
	}

	var total int
	app.DB.QueryRow("SELECT COUNT(*) FROM domains d WHERE "+where, args...).Scan(&total)

	query := fmt.Sprintf(`SELECT d.id, d.domain, d.template_id, d.market, d.language, d.keyword_type,
		d.primary_keyword, d.redirect_url, d.status, d.cloudflare_zone, d.created_at, d.updated_at,
		COALESCE(t.name,'') as template_name,
		CASE WHEN c.id IS NOT NULL THEN 1 ELSE 0 END as has_content
		FROM domains d
		LEFT JOIN templates t ON d.template_id=t.id
		LEFT JOIN contents c ON c.domain_id=d.id
		WHERE %s ORDER BY d.id DESC LIMIT ? OFFSET ?`, where)
	args = append(args, size, offset)

	rows, err := app.DB.Query(query, args...)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id, templateID int64
		var domain, market, lang, kwType, pk, redirect, status, cfZone, createdAt, updatedAt, tplName string
		var hasContent int
		rows.Scan(&id, &domain, &templateID, &market, &lang, &kwType, &pk, &redirect, &status, &cfZone, &createdAt, &updatedAt, &tplName, &hasContent)
		list = append(list, map[string]interface{}{
			"id": id, "domain": domain, "template_id": templateID, "market": market,
			"language": lang, "keyword_type": kwType, "primary_keyword": pk,
			"redirect_url": redirect, "status": status, "cloudflare_zone": cfZone,
			"created_at": createdAt, "updated_at": updatedAt,
			"template_name": tplName, "has_content": hasContent == 1,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonList(w, list, total)
}

func (app *App) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain         string `json:"domain"`
		TemplateID     int64  `json:"template_id"`
		Market         string `json:"market"`
		Language       string `json:"language"`
		KeywordType    string `json:"keyword_type"`
		PrimaryKeyword string `json:"primary_keyword"`
		RedirectURL    string `json:"redirect_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Domain == "" {
		jsonError(w, 400, "domain required")
		return
	}
	// validate domain format to prevent path traversal
	domainRe := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`)
	if !domainRe.MatchString(req.Domain) || strings.Contains(req.Domain, "..") {
		jsonError(w, 400, "invalid domain format")
		return
	}
	if req.Market == "" {
		req.Market = "zh-TW"
	}
	if req.Language == "" {
		req.Language = req.Market
	}
	if req.KeywordType == "" {
		req.KeywordType = "brand"
	}
	// auto-generate site_id from domain if not provided
	siteID := strings.ReplaceAll(strings.ReplaceAll(req.Domain, ".", "_"), "-", "_")
	res, err := app.DB.Exec(`INSERT INTO domains(domain,template_id,market,language,keyword_type,primary_keyword,redirect_url,site_id) VALUES(?,?,?,?,?,?,?,?)`,
		req.Domain, req.TemplateID, req.Market, req.Language, req.KeywordType, req.PrimaryKeyword, req.RedirectURL, siteID)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	// auto-create empty content row
	app.DB.Exec("INSERT OR IGNORE INTO contents(domain_id) VALUES(?)", id)
	jsonOK(w, map[string]interface{}{"id": id})
}

func (app *App) GetDomain(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var domain, market, lang, kwType, pk, redirect, sip, suser, spath, status, cfZone, createdAt, updatedAt string
	var templateID int64
	err = app.DB.QueryRow(`SELECT domain, template_id, market, language, keyword_type, primary_keyword,
		redirect_url, server_ip, server_user, server_path, status, cloudflare_zone, created_at, updated_at
		FROM domains WHERE id=?`, id).Scan(&domain, &templateID, &market, &lang, &kwType, &pk,
		&redirect, &sip, &suser, &spath, &status, &cfZone, &createdAt, &updatedAt)
	if err != nil {
		jsonError(w, 404, "domain not found")
		return
	}
	jsonOK(w, map[string]interface{}{
		"id": id, "domain": domain, "template_id": templateID, "market": market,
		"language": lang, "keyword_type": kwType, "primary_keyword": pk,
		"redirect_url": redirect, "server_ip": sip, "server_user": suser,
		"server_path": spath, "status": status, "cloudflare_zone": cfZone,
		"created_at": createdAt, "updated_at": updatedAt,
	})
}

func (app *App) UpdateDomain(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	allowed := map[string]bool{
		"domain": true, "template_id": true, "market": true, "language": true,
		"keyword_type": true, "primary_keyword": true, "redirect_url": true,
		"server_ip": true, "server_user": true, "server_path": true,
		"status": true, "cloudflare_zone": true, "content_source_id": true, "cluster_id": true,
	}
	for k, v := range req {
		if !allowed[k] {
			continue
		}
		app.DB.Exec(fmt.Sprintf("UPDATE domains SET %s=?, updated_at=CURRENT_TIMESTAMP WHERE id=?", k), v, id)
	}
	jsonOK(w, "ok")
}

func (app *App) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	// clean up all related data (contents handled by CASCADE)
	app.DB.Exec("DELETE FROM clicks WHERE site_id=(SELECT domain FROM domains WHERE id=?)", id)
	app.DB.Exec("DELETE FROM pageviews WHERE site_id=(SELECT domain FROM domains WHERE id=?)", id)
	app.DB.Exec("DELETE FROM ranking_history WHERE domain_id=?", id)
	app.DB.Exec("DELETE FROM site_cluster_members WHERE domain_id=?", id)
	app.DB.Exec("DELETE FROM domains WHERE id=?", id)
	jsonOK(w, "deleted")
}

func (app *App) BindTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		TemplateID int64 `json:"template_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	app.DB.Exec("UPDATE domains SET template_id=?, updated_at=CURRENT_TIMESTAMP WHERE id=?", req.TemplateID, id)
	jsonOK(w, "ok")
}

func (app *App) BatchDomainOp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs    []int64 `json:"ids"`
		Action string  `json:"action"` // build, deploy, build_deploy, delete
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	results := make([]map[string]interface{}, 0, len(req.IDs))
	for _, id := range req.IDs {
		res := map[string]interface{}{"id": id}
		switch req.Action {
		case "build":
			err := app.BuildSiteByID(id)
			res["status"] = "ok"
			if err != nil {
				res["status"] = "error"
				res["error"] = err.Error()
			}
		case "deploy":
			err := app.DeploySiteByID(id)
			res["status"] = "ok"
			if err != nil {
				res["status"] = "error"
				res["error"] = err.Error()
			}
		case "build_deploy":
			if err := app.BuildSiteByID(id); err != nil {
				res["status"] = "error"
				res["error"] = err.Error()
			} else if err := app.DeploySiteByID(id); err != nil {
				res["status"] = "error"
				res["error"] = err.Error()
			} else {
				res["status"] = "ok"
			}
		case "delete":
			app.DB.Exec("DELETE FROM clicks WHERE site_id=(SELECT domain FROM domains WHERE id=?)", id)
			app.DB.Exec("DELETE FROM pageviews WHERE site_id=(SELECT domain FROM domains WHERE id=?)", id)
			app.DB.Exec("DELETE FROM ranking_history WHERE domain_id=?", id)
			app.DB.Exec("DELETE FROM site_cluster_members WHERE domain_id=?", id)
			app.DB.Exec("DELETE FROM domains WHERE id=?", id)
			res["status"] = "deleted"
		default:
			res["status"] = "unknown action"
		}
		results = append(results, res)
	}
	jsonOK(w, results)
}
