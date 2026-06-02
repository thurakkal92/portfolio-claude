package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// LoadDotenv loads .env from the current working directory if present.
// Silent if the file does not exist (production sets real env vars).
func LoadDotenv() {
	_ = godotenv.Load()
}

type Config struct {
	Port            string
	DatabaseURL     string
	AllowedOrigins  []string
	ResendAPIKey    string
	ContactFrom     string
	ContactTo       string
	RateLimitPerHr  int
	RunMigrations   bool
	SeedOnStart     bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getenv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AllowedOrigins: splitCSV(getenv("ALLOWED_ORIGINS", "http://localhost:3000")),
		ResendAPIKey:   os.Getenv("RESEND_API_KEY"),
		ContactFrom:    getenv("CONTACT_FROM", "no-reply@thurakkal.com"),
		ContactTo:      getenv("CONTACT_TO", "nabeel.thurakkal92@gmail.com"),
		RateLimitPerHr: getenvInt("CONTACT_RATE_LIMIT_PER_HOUR", 5),
		RunMigrations:  getenvBool("RUN_MIGRATIONS", true),
		SeedOnStart:    getenvBool("SEED_ON_START", true),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return cfg, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getenvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getenvBool(k string, def bool) bool {
	v := strings.ToLower(os.Getenv(k))
	switch v {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
