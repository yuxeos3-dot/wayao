package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port       string
	DataDir    string
	HugoPath   string
	BuildOut   string
	TemplDir   string
	AdminDir   string
	TrackerURL string
}

func Load() *Config {
	c := &Config{
		Port:     getEnv("PORT", "8080"),
		DataDir:  getEnv("DATA_DIR", "./data"),
		HugoPath: getEnv("HUGO_PATH", "/usr/local/bin/hugo"),
		BuildOut: getEnv("BUILD_OUTPUT", "/var/www/sites"),
		TemplDir: getEnv("TEMPLATE_DIR", "./templates"),
		AdminDir: getEnv("ADMIN_DIR", "./frontend/dist"),
	}
	c.DataDir, _ = filepath.Abs(c.DataDir)
	c.TemplDir, _ = filepath.Abs(c.TemplDir)
	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
