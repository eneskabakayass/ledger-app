package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"ledger-app/internal/auth"
	"ledger-app/logger"
	"net/http"
	"strings"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			logger.Logger.Error("Missing Authorization header")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			logger.Logger.Error("Invalid token")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Logger.Error("Invalid token claims")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}

		fmt.Println(claims["userID"])
		c.Set("userID", claims["userID"])
		c.Set("role", claims["role"])

		return next(c)
	}
}
