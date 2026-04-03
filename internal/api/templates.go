package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *App) ListTemplates(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query(`SELECT id, name, slug, description, css_prefix, path, thumbnail, supported_kw_types, is_active, created_at FROM templates ORDER BY id DESC`)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id int64
		var name, slug, desc, cssPrefix, path, thumb, kwTypes string
		var isActive int
		var createdAt string
		rows.Scan(&id, &name, &slug, &desc, &cssPrefix, &path, &thumb, &kwTypes, &isActive, &createdAt)
		list = append(list, map[string]interface{}{
			"id": id, "name": name, "slug": slug, "description": desc,
			"css_prefix": cssPrefix, "path": path, "thumbnail": thumb,
			"supported_kw_types": kwTypes, "is_active": isActive == 1, "created_at": createdAt,
		})
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	jsonList(w, list, len(list))
}

func (app *App) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		CSSPrefix   string `json:"css_prefix"`
		Path        string `json:"path"`
		Thumbnail   string `json:"thumbnail"`
		SupportedKW string `json:"supported_kw_types"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	if req.Name == "" || req.Slug == "" {
		jsonError(w, 400, "name and slug required")
		return
	}
	if req.SupportedKW == "" {
		req.SupportedKW = "brand,game,sports,generic,promo,payment,affiliate,strategy,app,register,region,credit,live,community,terms"
	}
	res, err := app.DB.Exec(`INSERT INTO templates(name,slug,description,css_prefix,path,thumbnail,supported_kw_types) VALUES(?,?,?,?,?,?,?)`,
		req.Name, req.Slug, req.Description, req.CSSPrefix, req.Path, req.Thumbnail, req.SupportedKW)
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	id, _ := res.LastInsertId()
	jsonOK(w, map[string]interface{}{"id": id})
}

func (app *App) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
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
	allowed := map[string]bool{"name": true, "slug": true, "description": true, "css_prefix": true, "path": true, "thumbnail": true, "supported_kw_types": true, "is_active": true}
	for k, v := range req {
		if !allowed[k] {
			continue
		}
		app.DB.Exec(fmt.Sprintf("UPDATE templates SET %s=? WHERE id=?", k), v, id)
	}
	jsonOK(w, "ok")
}

func (app *App) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, 400, "invalid id")
		return
	}
	// check if any domain uses this template
	var count int
	app.DB.QueryRow("SELECT COUNT(*) FROM domains WHERE template_id=?", id).Scan(&count)
	if count > 0 {
		jsonError(w, 400, fmt.Sprintf("template in use by %d domains", count))
		return
	}
	app.DB.Exec("DELETE FROM templates WHERE id=?", id)
	jsonOK(w, "deleted")
}
