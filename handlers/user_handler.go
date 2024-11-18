package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"ledger-app/internal/connections/database"
	"ledger-app/logger"
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

		logger.Logger.WithFields(logrus.Fields{
			"details": errorMessage,
		}).Error("Validation failed")

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": errorMessage,
		})
	}

	if err := database.Db.Create(user).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to create user: %s", err.Error()))

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"Status":    http.StatusOK,
		"User Name": user.Name,
		"User ID":   user.ID,
	}).Info("Listen all users")

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func GetAllUser(c echo.Context) error {
	var users []models.User

	if err := database.Db.Find(&users).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to fetch user: %s", err.Error()))

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch user"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"Status":     http.StatusOK,
		"User Count": len(users),
	}).Info("Listen all users")

	return c.JSON(http.StatusOK, users)
}

func AddUserCredit(c echo.Context) error {
	userId := c.Param("id")
	creditReq := new(models.CreditRequest)

	if err := c.Bind(creditReq); err != nil || creditReq.Credit <= 0 {
		logger.Logger.Error("Invalid credit input")

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid credit input"})
	}

	var user models.User
	if err := database.Db.First(&user, userId).Error; err != nil {
		logger.Logger.Error("User not found")

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
	}

	oldCredit := user.Credit
	user.Credit += creditReq.Credit
	if err := database.Db.Save(&user).Error; err != nil {
		logger.Logger.Error("Failed to update credit")

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to update credit"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"User ID":      userId,
		"User Name":    user.Name,
		"Old Credit":   oldCredit,
		"New Credit":   user.Credit,
		"Added Credit": creditReq.Credit,
	}).Info("User credit updated successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Credit added successfully",
		"user":    user,
	})
}
