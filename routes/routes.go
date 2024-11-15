package routes

import (
	"github.com/labstack/echo/v4"
	"ledger-app/controllers"
	"ledger-app/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.LogRequest)
	e.GET("/users", controllers.GetAllUser)
	e.POST("/createUser", controllers.CreateUser)
}
