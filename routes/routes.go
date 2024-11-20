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
	e.GET("/user/:id/balance", handlers.GetUserBalance)
	e.GET("/users/balances", handlers.GetAllUsersTotalBalance)
	e.POST("/user/create", handlers.CreateUser)
	e.POST("/users/:id/credit", handlers.AddCreditToUser)
}
