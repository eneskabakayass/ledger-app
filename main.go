package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/config"
	"ledger-app/internal/connections/database"
	"ledger-app/routes"
)

func recoverPanic() {
	if r := recover(); r != nil {
		logrus.WithFields(logrus.Fields{
			"error": r,
		}).Error("Recovered from panic")
	}
}

func main() {
	defer recoverPanic()

	cfg := config.LoadEnvironment()
	e := echo.New()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	database.Connect()

	sqlDB, err := database.Db.DB()
	if err != nil {
		logrus.Fatal("Error getting DB connection", err)
	}

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			logrus.Error("Error closing DB connection", err)
		}
	}()

	routes.RegisterRoutes(e)

	logrus.Infof("Starting server at port %s", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		logrus.Fatal("Error starting server", err)
	}
}
