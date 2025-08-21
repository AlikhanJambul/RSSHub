package config

import (
	"RSSHub/internal/models"
	"os"
	"strconv"
)

func Load() models.Config {
	db := models.DB{
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "changem"),
		PostgresName: getEnv("POSTGRES_DB", "rsshub"),
	}

	return models.Config{
		DB:            db,
		WorkerCount:   getEnvInt("CLI_APP_WORKERS_COUNT", 3),
		TimerInterval: getEnv("CLI_APP_TIMER_INTERVAL", "3m"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return num
}
