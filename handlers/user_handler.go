package handlers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"ledger-app/internal/connections/database"
	"ledger-app/logger"
	"ledger-app/models"
	"net/http"
	"strconv"
	"time"
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

	if err := database.Db.Preload("Credits").Find(&user).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to fetch user: %s", err.Error()))

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to fetch user"})
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

	if err := database.Db.Preload("Credits").Find(&users).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to fetch user: %s", err.Error()))

		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to fetch users"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"Status":     http.StatusOK,
		"User Count": len(users),
	}).Info("Listen all users")

	return c.JSON(http.StatusOK, users)
}

func AddCreditToUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	creditReq := new(models.CreditRequest)
	if err := c.Bind(creditReq); err != nil {
		logger.Logger.Error("Invalid input: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if creditReq.Amount <= 0 {
		logger.Logger.Error("Invalid credit amount: must be positive")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Credit amount must be greater than zero"})
	}

	var user models.User
	if err := database.Db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found with ID: ", strconv.Itoa(userID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	credit := models.Transaction{
		UserID:          uint(userID),
		Amount:          creditReq.Amount,
		TransactionTime: time.Now().UTC(),
	}

	if err := database.Db.Create(&credit).Error; err != nil {
		logger.Logger.Error("Failed to add credit: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add credit"})
	}

	logger.Logger.Infof("Credit of %v added to User ID %d", creditReq.Amount, userID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Credit added successfully"})
}

func GetUserBalance(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	var user models.User
	if err := database.Db.Preload("Credits").First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found with ID: ", strconv.Itoa(userId))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	totalBalance := 0.0
	for _, credit := range user.Credits {
		totalBalance += credit.Amount
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":       user.ID,
		"user_name":     user.Name,
		"total_balance": totalBalance,
	})
}

func GetAllUsersTotalBalance(c echo.Context) error {
	var users []models.User

	if err := database.Db.Preload("Credits").Find(&users).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to fetch users: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}

	var userWithBalances []map[string]interface{}

	for _, user := range users {
		totalBalance := 0.0
		for _, credit := range user.Credits {
			totalBalance += credit.Amount
		}

		userWithBalances = append(userWithBalances, map[string]interface{}{
			"user_id":       user.ID,
			"user_name":     user.Name,
			"total_balance": totalBalance,
		})
	}

	logger.Logger.WithField("UserBalances", userWithBalances).Info("Listen all users with total balance")

	return c.JSON(http.StatusOK, userWithBalances)
}