package api

import (
	"fmt"
	"net/http"
)

func (app *App) StatsOverview(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{}

	var totalDomains, activeDomains, totalKW int
	app.DB.QueryRow("SELECT COUNT(*) FROM domains").Scan(&totalDomains)
	app.DB.QueryRow("SELECT COUNT(*) FROM domains WHERE status='active'").Scan(&activeDomains)
	app.DB.QueryRow("SELECT COUNT(*) FROM keywords").Scan(&totalKW)

	var totalClicks, totalPV, todayClicks, todayPV int
	app.DB.QueryRow("SELECT COUNT(*) FROM clicks").Scan(&totalClicks)
	app.DB.QueryRow("SELECT COUNT(*) FROM pageviews").Scan(&totalPV)
	app.DB.QueryRow("SELECT COUNT(*) FROM clicks WHERE date(created_at)=date('now')").Scan(&todayClicks)
	app.DB.QueryRow("SELECT COUNT(*) FROM pageviews WHERE date(created_at)=date('now')").Scan(&todayPV)

	result["total_domains"] = totalDomains
	result["active_domains"] = activeDomains
	result["total_keywords"] = totalKW
	result["total_clicks"] = totalClicks
	result["total_pv"] = totalPV
	result["today_clicks"] = todayClicks
	result["today_pv"] = todayPV

	// top domains by clicks (last 7 days)
	rows, _ := app.DB.Query(`SELECT site_id, COUNT(*) as cnt FROM clicks
		WHERE created_at >= datetime('now','-7 days') GROUP BY site_id ORDER BY cnt DESC LIMIT 10`)
	var topDomains []map[string]interface{}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var siteID string
			var cnt int
			rows.Scan(&siteID, &cnt)
			topDomains = append(topDomains, map[string]interface{}{"domain": siteID, "clicks": cnt})
		}
	}
	if topDomains == nil {
		topDomains = []map[string]interface{}{}
	}
	result["top_domains"] = topDomains

	// clicks by action
	rows2, _ := app.DB.Query("SELECT action, COUNT(*) FROM clicks WHERE date(created_at)=date('now') GROUP BY action")
	byAction := map[string]int{}
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var action string
			var cnt int
			rows2.Scan(&action, &cnt)
			byAction[action] = cnt
		}
	}
	result["clicks_by_action"] = byAction

	// hourly PV (last 24 hours) for Dashboard chart
	rows3, _ := app.DB.Query(`SELECT strftime('%H', created_at) as hour, COUNT(*) as cnt
		FROM pageviews WHERE created_at >= datetime('now', '-24 hours')
		GROUP BY hour ORDER BY hour`)
	hourPV := map[string]int{}
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var h string
			var c int
			rows3.Scan(&h, &c)
			hourPV[h] = c
		}
	}
	result["hour_pv"] = hourPV

	jsonOK(w, result)
}

func (app *App) GetClicks(w http.ResponseWriter, r *http.Request) {
	siteID := getQuery(r, "site_id", "")
	days := getQueryInt(r, "days", 7)
	if days < 1 || days > 365 {
		days = 7
	}
	page := getQueryInt(r, "page", 1)
	size := cappedSize(r, 50, 500)
	offset := (page - 1) * size

	where := fmt.Sprintf("created_at >= datetime('now','-%d days')", days)
	args := []interface{}{}
	if siteID != "" {
		where += " AND site_id=?"
		args = append(args, siteID)
	}

	var total int
	app.DB.QueryRow("SELECT COUNT(*) FROM clicks WHERE "+where, args...).Scan(&total)

	query := fmt.Sprintf("SELECT id, site_id, action, ip, ua, referer, page_url, is_fraud, created_at FROM clicks WHERE %s ORDER BY id DESC LIMIT ? OFFSET ?", where)
	args = append(args, size, offset)
	rows, err := app.DB.Query(query, args...)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id int64
		var siteID, action, ip, ua, referer, pageURL, createdAt string
		var isFraud int
		rows.Scan(&id, &siteID, &action, &ip, &ua, &referer, &pageURL, &isFraud, &createdAt)
		list = append(list, map[string]interface{}{
			"id": id, "site_id": siteID, "action": action,
			"ip": ip, "ua": ua, "referer": referer, "page_url": pageURL,
			"is_fraud": isFraud == 1, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonList(w, list, total)
}

func (app *App) GetBuildLog(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	rows, err := app.DB.Query("SELECT id, action, status, log_output, duration_ms, created_at FROM build_logs WHERE domain_id=? ORDER BY id DESC LIMIT 20", id)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var lid int64
		var action, status, logOutput, createdAt string
		var dur int
		rows.Scan(&lid, &action, &status, &logOutput, &dur, &createdAt)
		list = append(list, map[string]interface{}{
			"id": lid, "action": action, "status": status, "log_output": logOutput,
			"duration_ms": dur, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}
