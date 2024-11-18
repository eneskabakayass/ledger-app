package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/handlers"
	"ledger-app/internal/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.LogRequest)
	e.GET("/health", handlers.HealthCheck)
	e.GET("/users", handlers.GetAllUser)
	e.POST("/createUser", handlers.CreateUser)
	e.POST("/users/:id/credit", handlers.AddUserCredit)
}
