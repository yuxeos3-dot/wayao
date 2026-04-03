package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (app *App) ListKeywords(w http.ResponseWriter, r *http.Request) {
	cat := getQuery(r, "category", "")
	market := getQuery(r, "market", "")
	assigned := getQuery(r, "assigned", "")
	search := getQuery(r, "search", "")
	page := getQueryInt(r, "page", 1)
	size := cappedSize(r, 50, 500)
	offset := (page - 1) * size

	where := "1=1"
	args := []interface{}{}
	if cat != "" {
		where += " AND category=?"
		args = append(args, cat)
	}
	if market != "" {
		where += " AND market=?"
		args = append(args, market)
	}
	if assigned == "true" {
		where += " AND status='assigned'"
	} else if assigned == "false" {
		where += " AND status='unassigned'"
	}
	if search != "" {
		where += " AND keyword LIKE ?"
		args = append(args, "%"+search+"%")
	}

	var total int
	app.DB.QueryRow("SELECT COUNT(*) FROM keywords WHERE "+where, args...).Scan(&total)

	query := fmt.Sprintf("SELECT id, keyword, category, market, monthly_vol, difficulty, cpc, top1_dr, domain_id, status FROM keywords WHERE %s ORDER BY monthly_vol DESC LIMIT ? OFFSET ?", where)
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
		var keyword, category, mkt string
		var vol, difficulty, dr int
		var cpc float64
		var domainID *int64
		var status string
		rows.Scan(&id, &keyword, &category, &mkt, &vol, &difficulty, &cpc, &dr, &domainID, &status)
		list = append(list, map[string]interface{}{
			"id": id, "keyword": keyword, "category": category, "market": mkt,
			"monthly_vol": vol, "difficulty": difficulty, "cpc": cpc, "top1_dr": dr,
			"domain_id": domainID, "status": status,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonList(w, list, total)
}

func (app *App) KeywordCategories(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT category, COUNT(*) as cnt FROM keywords GROUP BY category ORDER BY cnt DESC")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	var list []map[string]interface{}
	for rows.Next() {
		var cat string
		var cnt int
		rows.Scan(&cat, &cnt)
		list = append(list, map[string]interface{}{"category": cat, "count": cnt})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonOK(w, list)
}

func (app *App) ImportKeywords(w http.ResponseWriter, r *http.Request) {
	// Accept CSV upload: keyword,category,market,monthly_vol,difficulty,cpc
	file, _, err := r.FormFile("file")
	if err != nil {
		jsonError(w, 400, "file required")
		return
	}
	defer file.Close()

	market := r.FormValue("market")
	if market == "" {
		market = "zh-TW"
	}

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// skip header
	reader.Read()

	tx, err := app.DB.Begin()
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	stmt, err := tx.Prepare("INSERT OR IGNORE INTO keywords(keyword,category,market,monthly_vol,difficulty,cpc) VALUES(?,?,?,?,?,?)")
	if err != nil {
		tx.Rollback()
		jsonError(w, 500, err.Error())
		return
	}
	defer stmt.Close()

	inserted := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if len(record) < 2 {
			continue
		}
		keyword := strings.TrimSpace(record[0])
		category := strings.TrimSpace(record[1])
		if keyword == "" || category == "" {
			continue
		}
		mkt := market
		if len(record) > 2 && record[2] != "" {
			mkt = record[2]
		}
		vol := 0
		difficulty := 0
		cpc := 0.0
		if len(record) > 3 {
			vol, _ = strconv.Atoi(record[3])
		}
		if len(record) > 4 {
			difficulty, _ = strconv.Atoi(record[4])
		}
		if len(record) > 5 {
			cpc, _ = strconv.ParseFloat(record[5], 64)
		}
		if _, err := stmt.Exec(keyword, category, mkt, vol, difficulty, cpc); err == nil {
			inserted++
		}
	}
	tx.Commit()
	jsonOK(w, map[string]interface{}{"inserted": inserted})
}

func (app *App) AssignKeyword(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	var req struct {
		DomainID int64 `json:"domain_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	app.DB.Exec("UPDATE keywords SET domain_id=?, status='assigned' WHERE id=?", req.DomainID, id)
	jsonOK(w, "assigned")
}

func (app *App) DeleteKeyword(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	app.DB.Exec("DELETE FROM keywords WHERE id=?", id)
	jsonOK(w, "deleted")
}
