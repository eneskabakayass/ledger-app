package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadEnvironment() Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		logrus.Error("DB_URL is required")
	}

	return Config{
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
