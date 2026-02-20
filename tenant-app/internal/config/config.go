package config

import (
	"os"
	"strings"
)

type Config struct {
	HTTPPort       string
	MetaDBConn     string
	KafkaBrokers   []string
	KafkaTopic     string
	JWTSigningKey  string
}

func Load() *Config {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		MetaDBConn:    getEnv("META_DB_CONN", "postgres://postgres:postgres@postgres:5432/tenant_meta?sslmode=disable"),
		KafkaBrokers:  getEnvSlice("KAFKA_BROKERS", []string{"kafka:29092"}),
		KafkaTopic:    getEnv("KAFKA_TOPIC", "tenant.signups"),
		JWTSigningKey: getEnv("JWT_SIGNING_KEY", "change-me-in-production"),
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
