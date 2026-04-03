package tracker

import (
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Tracker struct {
	DB       *sql.DB
	ipRules  map[string]string // ip_cidr -> rule_type
	uaRules  []uaRule
	mu       sync.RWMutex
	reloadAt time.Time
}

type uaRule struct {
	Pattern  string
	RuleType string
}

func New(db *sql.DB) *Tracker {
	t := &Tracker{DB: db, ipRules: make(map[string]string)}
	t.loadRules()
	return t
}

func (t *Tracker) loadRules() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.ipRules = make(map[string]string)
	rows, err := t.DB.Query("SELECT ip_cidr, rule_type FROM ip_rules WHERE expires_at IS NULL OR expires_at > datetime('now')")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cidr, ruleType string
			rows.Scan(&cidr, &ruleType)
			t.ipRules[cidr] = ruleType
		}
	}

	t.uaRules = nil
	rows2, err := t.DB.Query("SELECT pattern, rule_type FROM ua_rules WHERE is_active=1")
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var pattern, ruleType string
			rows2.Scan(&pattern, &ruleType)
			t.uaRules = append(t.uaRules, uaRule{Pattern: pattern, RuleType: ruleType})
		}
	}
	t.reloadAt = time.Now()
}

func (t *Tracker) isBlocked(ip, ua string) bool {
	t.mu.RLock()
	needReload := time.Since(t.reloadAt) > 5*time.Minute
	t.mu.RUnlock()
	if needReload {
		t.loadRules()
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	// check IP rules
	for cidr, ruleType := range t.ipRules {
		if ruleType == "block" {
			if ip == cidr || strings.HasPrefix(ip, strings.TrimSuffix(cidr, ".*")) {
				return true
			}
		}
	}
	// check UA rules
	uaLower := strings.ToLower(ua)
	for _, rule := range t.uaRules {
		if rule.RuleType == "block" && strings.Contains(uaLower, strings.ToLower(rule.Pattern)) {
			return true
		}
	}
	return false
}

func (t *Tracker) HandleTrack(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		SiteID  string `json:"s"`
		Action  string `json:"a"`
		Referer string `json:"r"`
		PageURL string `json:"p"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	if req.SiteID == "" {
		http.Error(w, "missing site_id", 400)
		return
	}
	if req.Action == "" {
		req.Action = "click"
	}

	ip := realIP(r)
	ua := r.UserAgent()

	if t.isBlocked(ip, ua) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	isFraud := 0
	if isKnownBot(ua) {
		isFraud = 1
	}

	// columns match db.go clicks table: site_id, action, ip, ua, referer, page_url, is_fraud
	t.DB.Exec(`INSERT INTO clicks(site_id, action, ip, ua, referer, page_url, is_fraud) VALUES(?,?,?,?,?,?,?)`,
		req.SiteID, req.Action, ip, ua, req.Referer, req.PageURL, isFraud)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":1}`))
}

func (t *Tracker) HandlePageView(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var req struct {
		SiteID  string `json:"s"`
		PageURL string `json:"p"`
		Referer string `json:"r"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", 400)
		return
	}
	if req.SiteID == "" {
		http.Error(w, "missing site_id", 400)
		return
	}
	if req.PageURL == "" {
		req.PageURL = "/"
	}

	ip := realIP(r)
	ua := r.UserAgent()

	if t.isBlocked(ip, ua) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// columns match db.go pageviews table: site_id, ip, ua, referer, page_url
	t.DB.Exec(`INSERT INTO pageviews(site_id, ip, ua, referer, page_url) VALUES(?,?,?,?,?)`,
		req.SiteID, ip, ua, req.Referer, req.PageURL)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":1}`))
}

func (t *Tracker) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); net.ParseIP(ip) != nil {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func isKnownBot(ua string) bool {
	ua = strings.ToLower(ua)
	bots := []string{"bot", "spider", "crawl", "slurp", "semrush", "ahrefs", "mj12bot", "dotbot", "petalbot"}
	for _, b := range bots {
		if strings.Contains(ua, b) {
			return true
		}
	}
	return false
}
