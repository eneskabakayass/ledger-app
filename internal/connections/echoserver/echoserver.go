package echoserver

import "github.com/labstack/echo/v4"

func GetInstance() *echo.Echo {
	e := echo.New()

	return e
}
