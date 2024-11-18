package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/config"
	"ledger-app/internal/connections/database"
	"ledger-app/logger"
	"ledger-app/routes"
)

func recoverPanic() {
	if r := recover(); r != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": r,
		}).Error("Recovered from panic")
	}
}

func main() {
	defer recoverPanic()

	logger.InitLogger()

	cfg := config.LoadEnvironment()
	e := echo.New()

	database.Connect()

	sqlDB, err := database.Db.DB()
	if err != nil {
		logger.Logger.Fatal("Error getting DB connection", err)
	}

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			logger.Logger.Error("Error closing DB connection", err)
		}
	}()

	routes.RegisterRoutes(e)

	logger.Logger.Infof("Starting server at port %s", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Logger.Fatal("Error starting server", err)
	}
}
