package controllers

import (
	"fmt"
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

	fmt.Println(user.Name)

	if user.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User name is required"})
	}

	query := "INSERT INTO users (name) VALUES (?)"
	result, err := db.Db.Exec(query, user.Name)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert user"})
	}

	userId, err := result.LastInsertId()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve last insert ID"})
	}

	user.ID = int(userId)

	return c.JSON(http.StatusCreated, user)
}

func GetAllUser(c echo.Context) error {
	if c.QueryParam("error") == "true" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	return c.String(http.StatusOK, "Get User List")
}
