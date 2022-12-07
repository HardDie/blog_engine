package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/HardDie/blog_engine/internal/logger"
)

type Config struct {
	DBPath         string
	Port           string
	PwdMaxAttempts int
	PwdBlockTime   int
}

func Get() *Config {
	if err := godotenv.Load(); err != nil {
		if check := os.IsNotExist(err); !check {
			logger.Error.Printf("failed to load env vars: %s", err)
		}
	}

	return &Config{
		DBPath:         getEnv("DB_PATH", "blog.db"),
		Port:           getEnv("PORT", ":8080"),
		PwdMaxAttempts: getEnvAsInt("PWD_MAX_ATTEMPTS", 5),
		PwdBlockTime:   getEnvAsInt("PWD_BLOCK_TIME", 24),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
func getEnvAsInt(key string, defaultValue int) int {
	value := getEnv(key, "")
	if v, e := strconv.Atoi(value); e == nil {
		return v
	}
	return defaultValue
}
