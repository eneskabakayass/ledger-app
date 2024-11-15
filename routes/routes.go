package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/handlers"
	"ledger-app/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.LogRequest)
	e.GET("/users", handlers.GetAllUser)
	e.POST("/createUser", handlers.CreateUser)
}
