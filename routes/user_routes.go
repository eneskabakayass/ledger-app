package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/handlers"
	"ledger-app/internal/middleware"
)

func RegisterUsersRoutes(e *echo.Echo) {
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.LoginUser)

	adminGroup := e.Group("/admin", middleware.JWTMiddleware, middleware.AdminMiddleware)
	adminGroup.GET("/users", handlers.GetAllUser)
	adminGroup.GET("/balances", handlers.GetAllUsersTotalBalance)
	adminGroup.POST("/users/:id/credit", handlers.AddCreditToUser)
	adminGroup.PUT("/users/:userID/role", handlers.UpdateUserRole)
	
	userGroup := e.Group("/users", middleware.JWTMiddleware)
	userGroup.GET("/:id/balance", handlers.GetUserBalance)
	userGroup.GET("/:id/time/balance", handlers.GetUserBalanceAtTime)
	userGroup.POST("/:sender_id/transfer/:receiver_id", handlers.TransferCredit)
	userGroup.POST("/:id/debit", handlers.UserWithdrawsCredit)
}
