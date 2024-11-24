package middleware

import (
	"github.com/labstack/echo/v4"
	"ledger-app/logger"
	"net/http"
)

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			logger.Logger.Error("Unauthorized access attempt to admin route")
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: Admins only"})
		}

		return next(c)
	}
}
