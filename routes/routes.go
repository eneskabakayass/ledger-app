package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/handlers"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", handlers.HealthCheck)
	e.GET("/users", handlers.GetAllUser)
	e.GET("/user/:id/balance", handlers.GetUserBalance)
	e.GET("/users/balances", handlers.GetAllUsersTotalBalance)
	e.POST("/user/create", handlers.CreateUser)
	e.POST("/users/:id/credit", handlers.AddCreditToUser)
	e.POST("/users/:sender_id/transfer/:receiver_id", handlers.TransferCredit)
}
