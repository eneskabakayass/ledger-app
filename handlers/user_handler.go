package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/internal/connections/database"
	"ledger-app/models"
	"net/http"
)

func CreateUser(c echo.Context) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := user.Validate(); err != nil {
		var errorMessage []string

		for _, err := range err.(validator.ValidationErrors) {
			errorMessage = append(errorMessage, fmt.Sprintf("Field %s failed validation: %s parameter: %s", err.Field(), err.Tag(), err.Param()))
		}

		logrus.WithFields(logrus.Fields{
			"details": errorMessage,
		}).Error("Validation failed")

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": errorMessage,
		})
	}

	if err := database.Db.Create(user).Error; err != nil {
		logrus.Error(fmt.Sprintf("Failed to create user: %s", err.Error()))

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	logrus.Info(fmt.Sprintf("Status:%d, User Name:%s, User ID:%d", http.StatusCreated, user.Name, user.ID))

	return c.JSON(http.StatusCreated, user)
}

func GetAllUser(c echo.Context) error {
	var users []models.User

	if err := database.Db.Find(&users).Error; err != nil {
		logrus.Error(fmt.Sprintf("Failed to fetch user: %s", err.Error()))

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	logrus.Info(fmt.Sprintf("Status:%d, User Count:%d", http.StatusOK, len(users)))

	return c.JSON(http.StatusOK, users)
}
