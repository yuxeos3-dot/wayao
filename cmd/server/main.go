package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuxeos3-dot/wayao/internal/api"
	"github.com/yuxeos3-dot/wayao/internal/builder"
	"github.com/yuxeos3-dot/wayao/internal/config"
	"github.com/yuxeos3-dot/wayao/internal/db"
	"github.com/yuxeos3-dot/wayao/internal/middleware"
	"github.com/yuxeos3-dot/wayao/internal/tracker"
)

func main() {
	cfg := config.Load()
<<<<<<< HEAD
	log.Printf("[MAIN] BrandSite Pro starting on :%s", cfg.Port)
=======
	log.Printf("[MAIN] Wayao CMS starting on :%s", cfg.Port)
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
	log.Printf("[MAIN] data=%s templates=%s hugo=%s", cfg.DataDir, cfg.TemplDir, cfg.HugoPath)

	database, err := db.InitDB(cfg.DataDir)
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	app := api.NewApp(database, cfg)
	bld := builder.New(database, cfg)
	trk := tracker.New(database)

	// wire build/deploy methods into app
	app.BuildFunc = bld.BuildSite
	app.DeployFunc = bld.DeploySite

	// Start background schedulers
	go app.StartContentRefreshScheduler()
	go app.StartHealthCheckScheduler()

	r := mux.NewRouter()

	// Rate limiter for tracker endpoints
	trackLimiter := middleware.NewRateLimiter(120, time.Minute)

	// === Public endpoints (no auth) ===
	loginLimiter := middleware.NewRateLimiter(10, time.Minute)
	r.Handle("/api/auth/login", loginLimiter.Middleware(http.HandlerFunc(app.HandleLogin))).Methods("POST")
	r.Handle("/t", trackLimiter.Middleware(http.HandlerFunc(trk.HandleTrack))).Methods("POST", "OPTIONS")
	r.Handle("/pv", trackLimiter.Middleware(http.HandlerFunc(trk.HandlePageView))).Methods("POST", "OPTIONS")
	r.HandleFunc("/health", trk.HandleHealth).Methods("GET")
	r.HandleFunc("/og/{site_id}.svg", app.HandleOGImage).Methods("GET")

	// === Authenticated API ===
	apiR := r.PathPrefix("/api/v1").Subrouter()
	apiR.Use(middleware.Auth(database))

	// Auth
	apiR.HandleFunc("/auth/logout", app.HandleLogout).Methods("POST")
	apiR.HandleFunc("/auth/change-password", app.HandleChangePassword).Methods("POST")

	// Stats
	apiR.HandleFunc("/stats/overview", app.StatsOverview).Methods("GET")
	apiR.HandleFunc("/stats/clicks", app.GetClicks).Methods("GET")

	// Templates
	apiR.HandleFunc("/templates", app.ListTemplates).Methods("GET")
	apiR.HandleFunc("/templates", app.CreateTemplate).Methods("POST")
	apiR.HandleFunc("/templates/{id}", app.UpdateTemplate).Methods("PUT")
	apiR.HandleFunc("/templates/{id}", app.DeleteTemplate).Methods("DELETE")

	// Domains
	apiR.HandleFunc("/domains", app.ListDomains).Methods("GET")
	apiR.HandleFunc("/domains", app.CreateDomain).Methods("POST")
	apiR.HandleFunc("/domains/{id}", app.GetDomain).Methods("GET")
	apiR.HandleFunc("/domains/{id}", app.UpdateDomain).Methods("PUT")
	apiR.HandleFunc("/domains/{id}", app.DeleteDomain).Methods("DELETE")
	apiR.HandleFunc("/domains/{id}/bind-template", app.BindTemplate).Methods("POST")
	apiR.HandleFunc("/domains/batch", app.BatchDomainOp).Methods("POST")

	// Content
	apiR.HandleFunc("/content/{id}", app.GetContent).Methods("GET")
	apiR.HandleFunc("/content/{id}", app.SaveContent).Methods("PUT", "POST")
	apiR.HandleFunc("/domains/{id}/ai-content", app.GenerateAIContent).Methods("POST")

	// Keywords
	apiR.HandleFunc("/keywords", app.ListKeywords).Methods("GET")
	apiR.HandleFunc("/keywords/categories", app.KeywordCategories).Methods("GET")
	apiR.HandleFunc("/keywords/import", app.ImportKeywords).Methods("POST")
	apiR.HandleFunc("/keywords/{id}/assign", app.AssignKeyword).Methods("POST")
	apiR.HandleFunc("/keywords/{id}", app.DeleteKeyword).Methods("DELETE")

	// Build / Deploy
	apiR.HandleFunc("/build/{id}", app.HandleBuild).Methods("POST")
	apiR.HandleFunc("/build/{id}/deploy", app.HandleDeploy).Methods("POST")
	apiR.HandleFunc("/build/{id}/full", app.HandleBuildAndDeploy).Methods("POST")
	apiR.HandleFunc("/build/{id}/status", app.GetBuildStatus).Methods("GET")
	apiR.HandleFunc("/build/{id}/log", app.GetBuildLog).Methods("GET")
	apiR.HandleFunc("/build/batch", app.HandleBatchBuild).Methods("POST")

	// Stats (extended)
	apiR.HandleFunc("/stats/domain/{id}", app.GetDomainStats).Methods("GET")
	apiR.HandleFunc("/stats/summary", app.GetDailySummary).Methods("GET")

	// Rankings
	apiR.HandleFunc("/rankings", app.ListRankings).Methods("GET")
	apiR.HandleFunc("/rankings/check/{id}", app.CheckRanking).Methods("POST")

	// Settings
	apiR.HandleFunc("/settings", app.GetSettings).Methods("GET")
	apiR.HandleFunc("/settings", app.SaveSettings).Methods("PUT", "POST")

	// V4: IndexNow
	apiR.HandleFunc("/indexnow/{id}/submit", app.SubmitIndexNow).Methods("POST")
	apiR.HandleFunc("/indexnow/{id}/records", app.GetIndexNowRecords).Methods("GET")

	// V4: City Matrix
	apiR.HandleFunc("/city-matrix/{id}", app.GetCityMatrix).Methods("GET")
	apiR.HandleFunc("/city-matrix/{id}", app.SaveCityMatrix).Methods("PUT")

	// V4: Title Pool
	apiR.HandleFunc("/title-pool", app.ListTitlePool).Methods("GET")
	apiR.HandleFunc("/title-pool", app.CreateTitleVariant).Methods("POST")
	apiR.HandleFunc("/title-pool/{id}", app.DeleteTitleVariant).Methods("DELETE")

	// V4: Site Clusters
	apiR.HandleFunc("/clusters", app.ListClusters).Methods("GET")
	apiR.HandleFunc("/clusters", app.CreateCluster).Methods("POST")
	apiR.HandleFunc("/clusters/{id}", app.DeleteCluster).Methods("DELETE")
	apiR.HandleFunc("/clusters/{id}/members", app.AddClusterMember).Methods("POST")
	apiR.HandleFunc("/clusters/{id}/members", app.RemoveClusterMember).Methods("DELETE")

	// V4: Content Refresh
	apiR.HandleFunc("/refresh-schedule", app.ListRefreshSchedule).Methods("GET")
	apiR.HandleFunc("/refresh-schedule", app.SaveRefreshSchedule).Methods("POST")
	apiR.HandleFunc("/refresh/{id}/run", app.RunRefreshNow).Methods("POST")

	// V4: Health Check
	apiR.HandleFunc("/health-check/{id}", app.CheckSiteHealth).Methods("GET")
	apiR.HandleFunc("/health-check/batch", app.BatchHealthCheck).Methods("POST")
	apiR.HandleFunc("/health-check/alerts", app.GetHealthAlerts).Methods("GET")

	// V4: Index Status
	apiR.HandleFunc("/index-status/{id}", app.CheckIndexStatus).Methods("GET")
	apiR.HandleFunc("/index-status/batch", app.BatchCheckIndex).Methods("POST")

	// V4: CTR Score
	apiR.HandleFunc("/ctr-score", app.GetCTRScore).Methods("POST")

	// V4: Batch Export
	apiR.HandleFunc("/export/batch", app.ExportBatch).Methods("POST")
	apiR.HandleFunc("/import", app.ImportDomain).Methods("POST")

	// V4: IP/UA Rules
	apiR.HandleFunc("/ip-rules", app.ListIPRules).Methods("GET")
	apiR.HandleFunc("/ip-rules", app.AddIPRule).Methods("POST")
	apiR.HandleFunc("/ip-rules/{id}", app.DeleteIPRule).Methods("DELETE")
	apiR.HandleFunc("/ua-rules", app.ListUARules).Methods("GET")
	apiR.HandleFunc("/ua-rules", app.AddUARule).Methods("POST")
	apiR.HandleFunc("/ua-rules/{id}", app.DeleteUARule).Methods("DELETE")

	// V4: Export
	apiR.HandleFunc("/export/{id}", app.ExportDomain).Methods("GET")

	// === Admin SPA ===
	adminDir := cfg.AdminDir
	if _, err := os.Stat(adminDir); err == nil {
		spa := spaHandler{staticPath: adminDir, indexPath: "index.html"}
		r.PathPrefix("/admin").Handler(http.StripPrefix("/admin", spa))
	}
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin/", http.StatusFound)
	})

	// wrap with CORS and body size limit (10MB max)
	handler := middleware.CORS(http.MaxBytesHandler(r, 10<<20))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("[MAIN] listening on :%s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// spaHandler serves the Vue SPA
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	fi, err := os.Stat(path)
	if err != nil || fi.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
