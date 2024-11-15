package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadEnvironment() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading not found .env file")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	return Config{
		Port:  getEnv("PORT", "3000"),
		DBUrl: getEnv("DB_URL", "root:password@tcp(localhost:3306)"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
