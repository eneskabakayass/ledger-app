package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/controllers"
	"ledger-app/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.LoggingMiddleware)
	e.GET("/api", controllers.GetAllUser)
	e.POST("/create/new-user", controllers.CreateUser)
}
