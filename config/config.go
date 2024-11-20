package config

import (
	"github.com/joho/godotenv"
	"ledger-app/logger"
	"os"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadEnvironment() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Error("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		logger.Logger.Error("DB_URL is required")
	}

	return &Config{
		Port:  getEnv("PORT", "3000"),
		DBUrl: getEnv("DB_URL", "root@tcp(localhost:3306)/ledger_app"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
