package handlers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"ledger-app/internal/connections/database"
	"ledger-app/internal/validation"
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

	if err := validation.ValidateStruct().Struct(user); err != nil {
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
	}).Info("Created new user")
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

	if len(users) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "No users found"})
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

	if err := validation.ValidateStruct().Struct(creditReq); err != nil {
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

	logger.Logger.Infof("User ID %d has total balance of %v", userId, totalBalance)
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

	if len(users) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "No users found"})
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

func TransferCredit(c echo.Context) error {
	senderID, err := strconv.Atoi(c.Param("sender_id"))
	if err != nil {
		logger.Logger.Error("Failed to convert sender ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid sender ID format"})
	}

	receiverID, err := strconv.Atoi(c.Param("receiver_id"))
	if err != nil {
		logger.Logger.Error("Failed to convert receiver ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid receiver ID format"})
	}

	creditReq := new(models.CreditRequest)
	if err := c.Bind(creditReq); err != nil {
		logger.Logger.Error("Invalid input: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := validation.ValidateStruct().Struct(creditReq); err != nil {
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

	var sender models.User
	if err := database.Db.Preload("Credits").First(&sender, senderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Sender not found with ID: ", strconv.Itoa(senderID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Sender not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	var receiver models.User
	if err := database.Db.Preload("Credits").First(&receiver, receiverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("Receiver not found with ID: ", strconv.Itoa(receiverID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Receiver not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	totalBalance := 0.0
	for _, credit := range sender.Credits {
		totalBalance += credit.Amount
	}

	if totalBalance < creditReq.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient balance"})
	}

	tx := database.Db.Begin()

	transaction := models.Transaction{
		UserID:          uint(senderID),
		Amount:          -creditReq.Amount,
		TransactionTime: time.Now().UTC(),
		SenderID:        uintPointer(uint(senderID)),
		ReceiverID:      uintPointer(uint(receiverID)),
	}

	receiverTransaction := models.Transaction{
		UserID:          uint(receiverID),
		Amount:          creditReq.Amount,
		TransactionTime: time.Now().UTC(),
		SenderID:        uintPointer(uint(senderID)),
		ReceiverID:      uintPointer(uint(receiverID)),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Logger.Error("Failed to add credit: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add credit"})
	}

	if err := tx.Create(&receiverTransaction).Error; err != nil {
		tx.Rollback()
		logger.Logger.Error("Failed to add credit: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add credit"})
	}

	if err := tx.Commit().Error; err != nil {
		logger.Logger.Error("Failed to commit transaction: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	logger.Logger.Infof("Credit of %v transferred from User ID %d to User ID %d", creditReq.Amount, senderID, receiverID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Credit transferred successfully"})
}

func uintPointer(val uint) *uint {
	return &val
}

func UserWithdrawsCredit(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))

	fmt.Println(userID)

	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	creditReq := new(models.CreditRequest)
	if err := c.Bind(creditReq); err != nil {
		logger.Logger.Error("Invalid input: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if err := validation.ValidateStruct().Struct(creditReq); err != nil {
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

	var user models.User
	if err := database.Db.Preload("Credits").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found with ID: ", strconv.Itoa(userID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	totalBalance := 0.0
	for _, credit := range user.Credits {
		totalBalance += credit.Amount
	}

	if totalBalance < creditReq.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient balance"})
	}

	tx := database.Db.Begin()

	transaction := &models.Transaction{
		UserID:          uint(userID),
		Amount:          -creditReq.Amount,
		TransactionTime: time.Now().UTC(),
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		logger.Logger.Error("Failed to add credit: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add credit"})
	}

	if err := tx.Commit().Error; err != nil {
		logger.Logger.Error("Failed to commit transaction: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
	}

	logger.Logger.Infof("Credit of %v withdrawn from User ID %d", creditReq.Amount, userID)
	return c.JSON(http.StatusOK, map[string]string{"message": "Credit withdrawn successfully"})
}

func GetUserBalanceAtTime(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	var user models.User
	if err := database.Db.Preload("Credits").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found with ID: ", strconv.Itoa(userID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	transactionTime, err := time.Parse(time.RFC3339, c.QueryParam("time"))
	if err != nil {
		logger.Logger.Error("Failed to parse time: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid time format"})
	}

	totalBalance := 0.0
	for _, credit := range user.Credits {
		if credit.TransactionTime.Before(transactionTime) {
			totalBalance += credit.Amount
		}
	}

	logger.Logger.Infof("User ID %d has total balance of %v at time %v", userID, totalBalance, transactionTime)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":       user.ID,
		"user_name":     user.Name,
		"total_balance": totalBalance,
	})
}
