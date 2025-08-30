package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"RSSHub/internal/domain"
)

func Load() domain.Config {
	// Пробуем загрузить .env вручную
	loadDotEnv(".env")

	db := domain.DB{
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresUser: getEnv("POSTGRES_USER", "postgres"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "changem"),
		PostgresName: getEnv("POSTGRES_DB", "rsshub"),
	}

	return domain.Config{
		DB:            db,
		WorkerCount:   getEnvInt("CLI_APP_WORKERS_COUNT", 3),
		TimerInterval: getEnv("CLI_APP_TIMER_INTERVAL", "3m"),
	}
}

func loadDotEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		// если файла нет — просто пропускаем
		log.Printf("⚠️ %s not found, using system environment variables\n", filename)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		os.Setenv(key, val)
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int32) int32 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return int32(num)
}
