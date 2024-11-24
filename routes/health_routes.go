package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/handlers"
)

func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/health", handlers.HealthCheck)
}
