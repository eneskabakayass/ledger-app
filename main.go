package main

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/config"
	"ledger-app/db"
	"ledger-app/routes"
)

func main() {
	cfg := config.LoadConfig()
	e := echo.New()

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	db.Connect()
	defer db.Db.Close()

	routes.RegisterRoutes(e)

	logrus.Infof("Starting server at port %s", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
