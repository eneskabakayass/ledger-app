package auth

import (
	"encoding/base64"
	"ledger-app/logger"
	"math/rand"
)

var secretKey string

func init() {
	secretKey = generateRandomKey()
}

func GetSecretKey() string {
	return secretKey
}

func generateRandomKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		logger.Logger.Fatalf("Failed to generate secret key %v", err)
	}
	return base64.StdEncoding.EncodeToString(key)
}
