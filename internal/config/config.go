package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DBPath        string
	CheckInterval int
	CheckTimeout  int
	WebhookURL    string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBPath:        getEnv("DB_PATH", "/data/uptime.db"),
		CheckInterval: getEnvInt("CHECK_INTERVAL", 60),
		CheckTimeout:  getEnvInt("CHECK_TIMEOUT", 10),
		WebhookURL:    getEnv("WEBHOOK_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
