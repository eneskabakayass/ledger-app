package handlers

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"ledger-app/internal/auth"
	"ledger-app/internal/connections/database"
	"ledger-app/internal/validation"
	"ledger-app/logger"
	"ledger-app/models"
	"net/http"
	"strconv"
	"time"
)

func GetAllUser(c echo.Context) error {
	var users []models.User

	if err := database.Db.Preload("Credits").Find(&users).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to fetch user: %s", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to fetch users"})
	}

	if len(users) == 0 {
		logger.Logger.Error(fmt.Sprintf("No users found"))
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
	tokenUserID := c.Get("userID")
	if tokenUserID == nil {
		logger.Logger.Error("Failed to retrieve user ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	requestUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	role, _ := c.Get("role").(string)

	if role != "admin" && uint(tokenUserID.(float64)) != uint(requestUserID) {
		logger.Logger.Warnf("User ID %d attempted to access balance of User ID %d", tokenUserID, requestUserID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
	}

	var user models.User
	if err := database.Db.Preload("Credits").First(&user, requestUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("User not found with ID: ", strconv.Itoa(requestUserID))
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		logger.Logger.Error("Database error: ", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	totalBalance := 0.0
	for _, credit := range user.Credits {
		totalBalance += credit.Amount
	}

	logger.Logger.Infof("User ID %d has total balance of %v", requestUserID, totalBalance)
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
		logger.Logger.Error("No users found")
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
	tokenUserID := c.Get("userID")
	if tokenUserID == nil {
		logger.Logger.Error("Failed to retrieve user ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

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

	role, _ := c.Get("role").(string)
	if role != "admin" && uint(tokenUserID.(float64)) != uint(senderID) {
		logger.Logger.Warnf("User ID %d attempted to access balance of User ID %d", tokenUserID, senderID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
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

func UserWithdrawsCredit(c echo.Context) error {
	tokenUserID := c.Get("userID")
	if tokenUserID == nil {
		logger.Logger.Error("Failed to retrieve user ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	role, _ := c.Get("role").(string)
	if role != "admin" && uint(tokenUserID.(float64)) != uint(userID) {
		logger.Logger.Warnf("User ID %d attempted to access balance of User ID %d", tokenUserID, tokenUserID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
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
	tokenUserID := c.Get("userID")
	if tokenUserID == nil {
		logger.Logger.Error("Failed to retrieve user ID from token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Logger.Error("Failed to convert user ID: ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	role, _ := c.Get("role").(string)
	if role != "admin" && uint(tokenUserID.(float64)) != uint(userID) {
		logger.Logger.Warnf("User ID %d attempted to access balance of User ID %d", tokenUserID, userID)
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied"})
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

func RegisterUser(c echo.Context) error {
	registerRoutes := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})

	if err := c.Bind(registerRoutes); err != nil {
		logger.Logger.Error("Invalid Input")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	var existingUser models.User
	if err := database.Db.Where("name = ?", registerRoutes.Username).First(&existingUser).Error; err == nil {
		logger.Logger.Error("Username already taken")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username already taken"})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerRoutes.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Error("Error hashing password")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error processing password"})
	}

	newUser := models.User{
		Name:         registerRoutes.Username,
		PasswordHash: string(passwordHash),
		IsAdmin:      false,
	}

	if err := database.Db.Create(&newUser).Error; err != nil {
		logger.Logger.Error(fmt.Sprintf("Failed to create user: %s", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	token, err := auth.GenerateToken(newUser.ID, "user")
	if err != nil {
		logger.Logger.Error("Error generating token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"user":  newUser,
		"token": token,
	}).Info("User registered successfully")
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user":    newUser,
		"token":   token,
	})
}

func LoginUser(c echo.Context) error {
	loginPayload := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})

	if err := c.Bind(loginPayload); err != nil {
		logger.Logger.Error("Invalid input")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	var user models.User
	if err := database.Db.Where("name = ?", loginPayload.Username).First(&user).Error; err != nil {
		logger.Logger.Warn("Invalid username or password")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginPayload.Password)); err != nil {
		logger.Logger.Warn("Invalid username or password")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	role := "user"
	if user.IsAdmin {
		role = "admin"
	}

	token, err := auth.GenerateToken(user.ID, role)
	if err != nil {
		logger.Logger.Error("Error generating token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating token"})
	}

	logger.Logger.WithFields(map[string]interface{}{
		"userID":   user.ID,
		"username": user.Name,
		"role":     role,
	}).Info("User logged in successfully")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   token,
	})
}

func uintPointer(val uint) *uint {
	return &val
}
