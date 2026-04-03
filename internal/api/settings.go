package api

import (
	"encoding/json"
	"net/http"
)

func (app *App) GetSettings(w http.ResponseWriter, r *http.Request) {
	rows, err := app.DB.Query("SELECT key, value FROM settings")
	if err != nil {
		jsonError(w, 500, err.Error())
		return
	}
	defer rows.Close()
	result := make(map[string]string)
	for rows.Next() {
		var k, v string
		rows.Scan(&k, &v)
		result[k] = v
	}
	// strip sensitive fields from response
	delete(result, "api_token")
	delete(result, "password_hash")
	jsonOK(w, result)
}

func (app *App) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, 400, "invalid json")
		return
	}
	// do not allow overwriting sensitive fields via general settings
	sensitiveKeys := map[string]bool{"api_token": true, "password_hash": true}
	for k, v := range req {
		if sensitiveKeys[k] {
			continue
		}
		app.DB.Exec("INSERT OR REPLACE INTO settings(key,value) VALUES(?,?)", k, v)
	}
	jsonOK(w, "saved")
}
