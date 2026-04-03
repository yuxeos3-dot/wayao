package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataDir string) (*sql.DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	dbPath := filepath.Join(dataDir, "brandsite.db")
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	// seed title_pool from data/titles/*.json
	seedTitlePool(db, dataDir)
	log.Printf("[DB] initialized at %s", dbPath)
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}

// seedTitlePool loads title patterns from JSON files into title_pool table
func seedTitlePool(db *sql.DB, dataDir string) {
	titlesDir := filepath.Join(dataDir, "titles")
	files, err := os.ReadDir(titlesDir)
	if err != nil {
		return
	}
	inserted := 0
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), "_titles.json") {
			continue
		}
		kwType := strings.TrimSuffix(f.Name(), "_titles.json")
		data, err := os.ReadFile(filepath.Join(titlesDir, f.Name()))
		if err != nil {
			continue
		}
		var patterns []string
		if json.Unmarshal(data, &patterns) != nil {
			continue
		}
		for _, p := range patterns {
			res, _ := db.Exec("INSERT OR IGNORE INTO title_pool(keyword_type, slot, template, market) VALUES(?, 'title', ?, 'zh-TW')", kwType, p)
			if res != nil {
				if n, _ := res.RowsAffected(); n > 0 {
					inserted++
				}
			}
		}
	}
	if inserted > 0 {
		log.Printf("[DB] seeded %d title patterns into title_pool", inserted)
	}
}

const schema = `
CREATE TABLE IF NOT EXISTS settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO settings(key, value) VALUES
    ('site_title', 'BrandSite Pro'),
    ('api_token', ''),
    ('tracker_url', ''),
    ('hugo_path', '/usr/local/bin/hugo'),
    ('sites_root', '/var/www/sites'),
    ('default_server_ip', 'localhost'),
    ('default_server_user', 'www-data'),
    ('claude_api_key', ''),
    ('ahrefs_api_key', ''),
    ('serpapi_key', ''),
    ('bing_wmt_key', ''),
    ('cc_limit_per_min', '120'),
    ('indexnow_auto', '1'),
    ('content_refresh_on', '1'),
    ('city_matrix_on', '0'),
    ('og_image_enabled', '1'),
    ('support_pages_auto', '1'),
    ('css_noise_enabled', '1'),
    ('schema_shuffle', '1');

CREATE TABLE IF NOT EXISTS site_clusters (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS templates (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    category    TEXT NOT NULL DEFAULT 'brand',
    description TEXT DEFAULT '',
    css_prefix  TEXT DEFAULT '',
    path        TEXT NOT NULL,
    thumbnail   TEXT DEFAULT '',
    supported_kw_types TEXT DEFAULT 'brand,game,sports,generic,promo,payment,affiliate,strategy,app,register,region,credit,live,community,terms',
    is_active   INTEGER DEFAULT 1,
    created_by  TEXT DEFAULT 'admin',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tpl_category ON templates(category);
CREATE INDEX IF NOT EXISTS idx_tpl_active ON templates(is_active);

CREATE TABLE IF NOT EXISTS domains (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    domain           TEXT NOT NULL UNIQUE,
    template_id      INTEGER REFERENCES templates(id),
    market           TEXT NOT NULL DEFAULT 'zh-TW',
    language         TEXT NOT NULL DEFAULT 'zh-TW',
    keyword_type     TEXT NOT NULL DEFAULT 'brand',
    primary_keyword  TEXT DEFAULT '',
    redirect_url     TEXT DEFAULT '',
    redirect_mode    TEXT DEFAULT 'tracker',
    site_id          TEXT NOT NULL DEFAULT '',
    server_ip        TEXT DEFAULT '',
    server_user      TEXT DEFAULT 'root',
    server_path      TEXT DEFAULT '',
    status           TEXT DEFAULT 'draft',
    cloudflare_zone  TEXT DEFAULT '',
    content_source_id INTEGER REFERENCES domains(id) ON DELETE SET NULL,
    cluster_id       INTEGER REFERENCES site_clusters(id) ON DELETE SET NULL,
    total_pv         INTEGER DEFAULT 0,
    total_click      INTEGER DEFAULT 0,
    last_built_at    DATETIME,
    last_deployed_at DATETIME,
    build_error      TEXT DEFAULT '',
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_domain_status ON domains(status);
CREATE INDEX IF NOT EXISTS idx_domain_market ON domains(market);

CREATE TABLE IF NOT EXISTS contents (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id        INTEGER NOT NULL UNIQUE REFERENCES domains(id) ON DELETE CASCADE,
    keyword_type     TEXT DEFAULT '',
    target_keyword   TEXT DEFAULT '',
    page_title       TEXT DEFAULT '',
    meta_desc        TEXT DEFAULT '',
    h1               TEXT DEFAULT '',
    brand_name       TEXT DEFAULT '',
    brand_color      TEXT DEFAULT '#1976D2',
    cta_text         TEXT DEFAULT '',
    cta_sub          TEXT DEFAULT '',
    hero_title       TEXT DEFAULT '',
    hero_subtitle    TEXT DEFAULT '',
    feature_1_icon   TEXT DEFAULT '',
    feature_1_title  TEXT DEFAULT '',
    feature_1_desc   TEXT DEFAULT '',
    feature_2_icon   TEXT DEFAULT '',
    feature_2_title  TEXT DEFAULT '',
    feature_2_desc   TEXT DEFAULT '',
    feature_3_icon   TEXT DEFAULT '',
    feature_3_title  TEXT DEFAULT '',
    feature_3_desc   TEXT DEFAULT '',
    intro_text       TEXT DEFAULT '',
    body_content     TEXT DEFAULT '',
    conclusion       TEXT DEFAULT '',
    faq_title        TEXT DEFAULT '常見問題',
    faq_items        TEXT DEFAULT '[]',
    extra_data       TEXT DEFAULT '{}',
    author_name      TEXT DEFAULT '',
    author_title     TEXT DEFAULT '',
    author_bio       TEXT DEFAULT '',
    author_avatar    TEXT DEFAULT '',
    last_updated     TEXT DEFAULT '',
    last_updated_iso TEXT DEFAULT '',
    review_count     TEXT DEFAULT '',
    trust_badges     TEXT DEFAULT '[]',
    disclosure       TEXT DEFAULT '',
    disclaimer       TEXT DEFAULT '',
    related_pages    TEXT DEFAULT '[]',
    content_angle    TEXT DEFAULT '',
    status           TEXT DEFAULT 'draft',
    content_hash     TEXT DEFAULT '',
    ai_generated     INTEGER DEFAULT 0,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at       DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_content_domain ON contents(domain_id);
CREATE INDEX IF NOT EXISTS idx_content_type ON contents(keyword_type);

CREATE TABLE IF NOT EXISTS keywords (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    keyword      TEXT NOT NULL,
    category     TEXT NOT NULL,
    market       TEXT NOT NULL DEFAULT 'zh-TW',
    monthly_vol  INTEGER DEFAULT 0,
    difficulty   INTEGER DEFAULT 0,
    cpc          REAL DEFAULT 0,
    top1_dr      INTEGER DEFAULT 0,
    domain_id    INTEGER REFERENCES domains(id) ON DELETE SET NULL,
    status       TEXT DEFAULT 'unassigned',
    current_rank INTEGER DEFAULT 0,
    best_rank    INTEGER DEFAULT 0,
    last_checked DATETIME,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(keyword, market)
);
CREATE INDEX IF NOT EXISTS idx_kw_cat ON keywords(category);
CREATE INDEX IF NOT EXISTS idx_kw_market ON keywords(market);
CREATE INDEX IF NOT EXISTS idx_kw_domain ON keywords(domain_id);

CREATE TABLE IF NOT EXISTS clicks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id    TEXT NOT NULL,
    domain_id  INTEGER REFERENCES domains(id),
    action     TEXT DEFAULT 'click',
    label      TEXT DEFAULT '',
    ip         TEXT DEFAULT '',
    ua         TEXT DEFAULT '',
    referer    TEXT DEFAULT '',
    page_url   TEXT DEFAULT '',
    is_fraud   INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_click_site ON clicks(site_id);
CREATE INDEX IF NOT EXISTS idx_click_date ON clicks(created_at);
CREATE INDEX IF NOT EXISTS idx_click_domain ON clicks(domain_id);

CREATE TABLE IF NOT EXISTS pageviews (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    site_id    TEXT NOT NULL,
    domain_id  INTEGER REFERENCES domains(id),
    ip         TEXT DEFAULT '',
    ua         TEXT DEFAULT '',
    referer    TEXT DEFAULT '',
    page_url   TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_pv_site ON pageviews(site_id);
CREATE INDEX IF NOT EXISTS idx_pv_date ON pageviews(created_at);

CREATE TABLE IF NOT EXISTS build_logs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id   INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    action      TEXT NOT NULL DEFAULT 'build',
    status      TEXT DEFAULT 'pending',
    log_output  TEXT DEFAULT '',
    duration_ms INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_blog_domain ON build_logs(domain_id);
CREATE INDEX IF NOT EXISTS idx_blog_date ON build_logs(created_at);

CREATE TABLE IF NOT EXISTS ranking_history (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id  INTEGER REFERENCES domains(id) ON DELETE CASCADE,
    keyword_id INTEGER REFERENCES keywords(id),
    keyword    TEXT NOT NULL,
    rank       INTEGER DEFAULT 0,
    market     TEXT DEFAULT 'zh-TW',
    engine     TEXT DEFAULT 'google',
    checked_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_rank_domain ON ranking_history(domain_id);
CREATE INDEX IF NOT EXISTS idx_rank_keyword ON ranking_history(keyword_id);
CREATE INDEX IF NOT EXISTS idx_rank_date ON ranking_history(checked_at);

CREATE TABLE IF NOT EXISTS indexnow_keys (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id  INTEGER NOT NULL UNIQUE REFERENCES domains(id) ON DELETE CASCADE,
    api_key    TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS index_submissions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id  INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    engine     TEXT DEFAULT 'indexnow',
    status     TEXT DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_idxsub_domain ON index_submissions(domain_id);

CREATE TABLE IF NOT EXISTS city_matrix (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id   INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    city_name   TEXT NOT NULL,
    city_slug   TEXT NOT NULL,
    extra_title TEXT DEFAULT '',
    extra_desc  TEXT DEFAULT '',
    is_built    INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(domain_id, city_slug)
);
CREATE INDEX IF NOT EXISTS idx_city_domain ON city_matrix(domain_id);

CREATE TABLE IF NOT EXISTS title_pool (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    keyword_type TEXT NOT NULL,
    slot         TEXT NOT NULL DEFAULT 'title',
    template     TEXT NOT NULL,
    market       TEXT DEFAULT 'zh-TW',
    weight       INTEGER DEFAULT 1,
    is_active    INTEGER DEFAULT 1,
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tpool_type ON title_pool(keyword_type, slot, market);

CREATE TABLE IF NOT EXISTS site_cluster_members (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id INTEGER NOT NULL REFERENCES site_clusters(id) ON DELETE CASCADE,
    domain_id  INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    role       TEXT DEFAULT 'member',
    UNIQUE(cluster_id, domain_id)
);

CREATE TABLE IF NOT EXISTS content_refresh_schedule (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    domain_id       INTEGER NOT NULL UNIQUE REFERENCES domains(id) ON DELETE CASCADE,
    refresh_type    TEXT DEFAULT 'timestamp',
    frequency_days  INTEGER DEFAULT 14,
    last_refreshed  DATETIME,
    next_refresh    DATETIME,
    is_active       INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS ip_rules (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_cidr    TEXT NOT NULL,
    rule_type  TEXT DEFAULT 'block',
    reason     TEXT DEFAULT '',
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ua_rules (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern    TEXT NOT NULL,
    rule_type  TEXT DEFAULT 'block',
    is_active  INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 预置UA黑名单
INSERT OR IGNORE INTO ua_rules (pattern, rule_type, is_active) VALUES
    ('ahrefsbot', 'block', 1),
    ('semrushbot', 'block', 1),
    ('mj12bot', 'block', 1),
    ('dotbot', 'block', 1),
    ('blexbot', 'block', 1),
    ('siteexplorer', 'block', 1);
`
