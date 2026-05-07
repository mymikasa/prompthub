package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv             string
	HTTPAddr           string
	MySQLDSN           string
	SessionSecret      string
	SecretEncryptionKey string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		HTTPAddr:           getEnv("HTTP_ADDR", ":8080"),
		MySQLDSN:           os.Getenv("MYSQL_DSN"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		SecretEncryptionKey: os.Getenv("SECRET_ENCRYPTION_KEY"),
	}

	if cfg.MySQLDSN == "" {
		return nil, fmt.Errorf("MYSQL_DSN is required")
	}
	return cfg, nil
}

func (c *Config) IsDev() bool {
	return c.AppEnv == "development"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
