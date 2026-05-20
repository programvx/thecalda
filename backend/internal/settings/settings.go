// Package settings loads runtime configuration from the environment.
package settings

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Settings holds all runtime configuration for the API.
type Settings struct {
	Environment      string   // "dev" or "prod"
	Port             string   // HTTP port the API listens on
	DatabaseURL      string   // Postgres connection string
	SupabaseURL      string   // Base URL of the Supabase stack (used for JWKS)
	CORSAllowOrigins []string // Origins permitted by the CORS policy
}

// NewSettings reads configuration from a .env file (if present) and the
// environment.
func NewSettings() *Settings {
	_ = godotenv.Load()

	return &Settings{
		Environment:      getEnv("ENVIRONMENT", "dev"),
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgresql://postgres:postgres@127.0.0.1:54322/postgres"),
		SupabaseURL:      getEnv("SUPABASE_URL", "http://127.0.0.1:54321"),
		CORSAllowOrigins: splitEnv("CORS_ALLOW_ORIGINS", ",", "http://localhost:3000", "http://127.0.0.1:3000"),
	}
}

// IsDev reports whether the API is running outside production.
func (s *Settings) IsDev() bool {
	return s.Environment != "prod"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// splitEnv splits a delimited env var into a slice, or returns fallback when unset.
func splitEnv(key, sep string, fallback ...string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	parts := strings.Split(raw, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
