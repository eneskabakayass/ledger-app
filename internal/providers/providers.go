package providers

import (
	"github.com/labstack/echo/v4"
	"ledger-app/config"
	"ledger-app/internal/connections/database"
	"ledger-app/internal/middleware"
	"ledger-app/logger"
	"ledger-app/routes"
)

func InitLogger() {
	logger.InitLogger()
}

func InitDatabase() {
	database.Connect()
}

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middleware.LogRequest)
	routes.RegisterRoutes(e)
}

func StartServer(e *echo.Echo, cfg *config.Config) {
	logger.Logger.Infof("Starting server at port %s", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Logger.Fatal("Error starting server", err)
	}
}
