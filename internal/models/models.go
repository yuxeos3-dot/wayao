package models

import "time"

type Domain struct {
	ID              int64     `json:"id"`
	Domain          string    `json:"domain"`
	TemplateID      int64     `json:"template_id"`
	Market          string    `json:"market"`
	Language        string    `json:"language"`
	KeywordType     string    `json:"keyword_type"`
	PrimaryKeyword  string    `json:"primary_keyword"`
	RedirectURL     string    `json:"redirect_url"`
	ServerIP        string    `json:"server_ip"`
	ServerUser      string    `json:"server_user"`
	ServerPath      string    `json:"server_path"`
	Status          string    `json:"status"`
	CloudflareZone  string    `json:"cloudflare_zone"`
	ContentSourceID *int64    `json:"content_source_id"`
	ClusterID       *int64    `json:"cluster_id"`
	HasContent      bool      `json:"has_content"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Template struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CSSPrefix   string    `json:"css_prefix"`
	Path        string    `json:"path"`
	Thumbnail   string    `json:"thumbnail"`
	SupportedKW string    `json:"supported_kw_types"` // comma-separated
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type Content struct {
	ID               int64     `json:"id"`
	DomainID         int64     `json:"domain_id"`
	SiteName         string    `json:"site_name"`
	SiteSlogan       string    `json:"site_slogan"`
	BrandColor       string    `json:"brand_color"`
	LogoText         string    `json:"logo_text"`
	HeroTitle        string    `json:"hero_title"`
	HeroDescription  string    `json:"hero_description"`
	MainContent      string    `json:"main_content"`
	SidebarWidget    string    `json:"sidebar_widget"`
	CTAButtonText    string    `json:"cta_button_text"`
	CTAButtonURL     string    `json:"cta_button_url"`
	MetaTitle        string    `json:"meta_title"`
	MetaDescription  string    `json:"meta_description"`
	FAQJson          string    `json:"faq_json"`
	ProsConsJson     string    `json:"pros_cons_json"`
	ExtraFieldsJson  string    `json:"extra_fields_json"`
	ContentQuality   int       `json:"content_quality"`
	ContentHash      string    `json:"content_hash"`
	AIGenerated      bool      `json:"ai_generated"`
	ReviewedBy       string    `json:"reviewed_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Keyword struct {
	ID          int64   `json:"id"`
	Keyword     string  `json:"keyword"`
	Category    string  `json:"category"`
	Market      string  `json:"market"`
	Volume      int     `json:"volume"`
	KD          int     `json:"kd"`
	CPC         float64 `json:"cpc"`
	Top1DR      int     `json:"top1_dr"`
	AssignedTo  *int64  `json:"assigned_to"`
	IsAssigned  bool    `json:"is_assigned"`
}

type Click struct {
	ID        int64     `json:"id"`
	SiteID    string    `json:"site_id"`
	Domain    string    `json:"domain"`
	Action    string    `json:"action"`
	Label     string    `json:"label"`
	IP        string    `json:"ip"`
	UA        string    `json:"ua"`
	Referrer  string    `json:"referrer"`
	Country   string    `json:"country"`
	IsFraud   bool      `json:"is_fraud"`
	CreatedAt time.Time `json:"created_at"`
}

type PageView struct {
	ID        int64     `json:"id"`
	SiteID    string    `json:"site_id"`
	Domain    string    `json:"domain"`
	Path      string    `json:"path"`
	IP        string    `json:"ip"`
	UA        string    `json:"ua"`
	Referrer  string    `json:"referrer"`
	CreatedAt time.Time `json:"created_at"`
}

type BuildLog struct {
	ID        int64     `json:"id"`
	DomainID  int64     `json:"domain_id"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Duration  int       `json:"duration_ms"`
	CreatedAt time.Time `json:"created_at"`
}

type Setting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RankingHistory struct {
	ID        int64     `json:"id"`
	DomainID  int64     `json:"domain_id"`
	Keyword   string    `json:"keyword"`
	Position  int       `json:"position"`
	URL       string    `json:"url"`
	CheckedAt time.Time `json:"checked_at"`
}

type IndexNowKey struct {
	ID        int64     `json:"id"`
	DomainID  int64     `json:"domain_id"`
	APIKey    string    `json:"api_key"`
	CreatedAt time.Time `json:"created_at"`
}

type IndexSubmission struct {
	ID        int64     `json:"id"`
	DomainID  int64     `json:"domain_id"`
	URL       string    `json:"url"`
	Engine    string    `json:"engine"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CityMatrix struct {
	ID          int64     `json:"id"`
	DomainID    int64     `json:"domain_id"`
	CityName    string    `json:"city_name"`
	CitySlug    string    `json:"city_slug"`
	ExtraTitle  string    `json:"extra_title"`
	ExtraDesc   string    `json:"extra_desc"`
	IsBuilt     bool      `json:"is_built"`
	CreatedAt   time.Time `json:"created_at"`
}

type TitleVariant struct {
	ID          int64  `json:"id"`
	KeywordType string `json:"keyword_type"`
	Pattern     string `json:"pattern"`
	IsActive    bool   `json:"is_active"`
}

type SiteCluster struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Strategy  string    `json:"strategy"`
	CreatedAt time.Time `json:"created_at"`
}

type SiteClusterMember struct {
	ID        int64  `json:"id"`
	ClusterID int64  `json:"cluster_id"`
	DomainID  int64  `json:"domain_id"`
	Role      string `json:"role"`
}

type ContentRefreshSchedule struct {
	ID         int64     `json:"id"`
	DomainID   int64     `json:"domain_id"`
	CronExpr   string    `json:"cron_expr"`
	LastRun    time.Time `json:"last_run"`
	NextRun    time.Time `json:"next_run"`
	IsActive   bool      `json:"is_active"`
}

type IPRule struct {
	ID        int64     `json:"id"`
	IP        string    `json:"ip"`
	Action    string    `json:"action"` // block, allow
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type UaRule struct {
	ID        int64     `json:"id"`
	Pattern   string    `json:"pattern"`
	Action    string    `json:"action"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type StatsOverview struct {
	TotalDomains   int            `json:"total_domains"`
	ActiveDomains  int            `json:"active_domains"`
	TotalKeywords  int            `json:"total_keywords"`
	TotalClicks    int            `json:"total_clicks"`
	TotalPV        int            `json:"total_pv"`
	TodayClicks    int            `json:"today_clicks"`
	TodayPV        int            `json:"today_pv"`
	AvgPosition    float64        `json:"avg_position"`
	TopDomains     []DomainStats  `json:"top_domains"`
	ClicksByAction map[string]int `json:"clicks_by_action"`
}

type DomainStats struct {
	Domain string `json:"domain"`
	Clicks int    `json:"clicks"`
	PV     int    `json:"pv"`
}
