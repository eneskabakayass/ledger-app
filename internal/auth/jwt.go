package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(GetSecretKey()), nil
	})

}

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signInToken, err := token.SignedString([]byte(GetSecretKey()))

	if err != nil {
		return "", nil
	}

	return signInToken, nil
}
