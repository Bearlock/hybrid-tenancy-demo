package config

import (
	"os"
	"strings"
)

const AppName = "org-app"

type Config struct {
	HTTPPort       string
	TenantDBConn   string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	KafkaBrokers   []string
	KafkaTopic     string
}

func Load() *Config {
	return &Config{
		HTTPPort:     getEnv("HTTP_PORT", "8002"),
		TenantDBConn: getEnv("TENANT_DB_CONN", "postgres://postgres:postgres@postgres:5432/org_app_tenant_registry?sslmode=disable"),
		DBHost:       getEnv("DB_HOST", "postgres"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "postgres"),
		KafkaBrokers: getEnvSlice("KAFKA_BROKERS", []string{"kafka:29092"}),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "tenant.signups"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvSlice(key string, defaultVal []string) []string {
	if v := os.Getenv(key); v != "" {
		return strings.Split(v, ",")
	}
	return defaultVal
}
