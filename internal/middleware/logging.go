package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/logger"
)

func LogRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Logger.WithFields(logrus.Fields{
			"Method": c.Request().Method,
			"Url":    c.Request().URL.String(),
		}).Info("Incoming request")

		return next(c)
	}
}
