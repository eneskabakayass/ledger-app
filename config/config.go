package config

import (
	"github.com/joho/godotenv"
	"ledger-app/logger"
	"os"
)

type Config struct {
	Port                 string
	DBUrl                string
	DefaultAdminUserName string
	DefaultAdminPassword string
}

func LoadEnvironment() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Logger.Error("Error loading .env file")
	}

	return &Config{
		Port:                 getEnv("PORT", "80"),
		DBUrl:                getEnv("DB_URL", "root:12345@tcp(db:3306)/ledger_app"),
		DefaultAdminUserName: getEnv("DEFAULT_ADMIN_USERNAME", "admin"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
