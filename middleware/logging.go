package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("Method %s, URI: %s", c.Request().Method, c.Request().RequestURI)
		return next(c)
	}
}
