package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"os"
)

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal("Error opening log file", err)
	}

	logrus.SetOutput(file)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return func(c echo.Context) error {
		logrus.WithFields(logrus.Fields{
			"Method": c.Request().Method,
			"Url":    c.Request().URL.String(),
		}).Info("Incoming request")

		return next(c)
	}
}
