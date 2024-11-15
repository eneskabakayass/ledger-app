package controllers

import (
	"github.com/labstack/echo/v4"
	"ledger-app/db"
	"ledger-app/models"
	"net/http"
)

func CreateUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if user.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User name is required"})
	}

	if err := db.Db.Create(user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func GetAllUser(c echo.Context) error {
	return c.String(http.StatusOK, "All User List")
}
