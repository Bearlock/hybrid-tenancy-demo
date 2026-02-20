package config

import (
	"os"
)

type Config struct {
	HTTPPort      string
	JWTSigningKey string
	FactAppURL    string
	OrgAppURL     string
	TodoAppURL    string
}

func Load() *Config {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8000"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY", "change-me-in-production"),
		FactAppURL:    getEnv("FACT_APP_URL", "http://localhost:8001"),
		OrgAppURL:     getEnv("ORG_APP_URL", "http://localhost:8002"),
		TodoAppURL:    getEnv("TODO_APP_URL", "http://localhost:8003"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
