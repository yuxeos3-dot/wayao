package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yuxeos3-dot/wayao/internal/config"
)

type App struct {
	DB         *sql.DB
	Cfg        *config.Config
	BuildFunc  func(int64) error
	DeployFunc func(int64) error
}

func NewApp(db *sql.DB, cfg *config.Config) *App {
	return &App{DB: db, Cfg: cfg}
}

func jsonOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": data})
}

func jsonList(w http.ResponseWriter, list interface{}, total int) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"code": 0, "data": list, "total": total})
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"code": -1, "error": msg})
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
}

func getQuery(r *http.Request, key, fallback string) string {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	return v
}

func getQueryInt(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

// maxPageSize caps the page size to prevent DoS via huge queries
func cappedSize(r *http.Request, fallback, max int) int {
	n := getQueryInt(r, "size", fallback)
	if n < 1 {
		return fallback
	}
	if n > max {
		return max
	}
	return n
}

func getSetting(db *sql.DB, key string) string {
	var val string
	db.QueryRow("SELECT value FROM settings WHERE key=?", key).Scan(&val)
	return val
}
