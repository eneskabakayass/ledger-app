package handlers

import (
	"github.com/labstack/echo/v4"
	"ledger-app/internal/connections/database"
	"ledger-app/logger"
	"ledger-app/models"
	"net/http"
	"strconv"
)

type RoleUpdatePayload struct {
	Role string `json:"role"`
}

func IsAdmin(userID uint) (bool, error) {
	var user models.User
	if err := database.Db.First(&user, userID).Error; err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}
func UpdateUserRole(c echo.Context) error {
	adminUserId := c.Get("userID")
	isAdmin, err := IsAdmin(uint(adminUserId.(float64)))

	if err != nil && !isAdmin {
		logger.Logger.Error("Error checking if user is admin: ", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	userId := c.Param("userID")

	targetUserID, err := strconv.Atoi(userId)
	if err != nil {
		logger.Logger.Error("Error converting user ID to integer: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var payload RoleUpdatePayload
	if err := c.Bind(&payload); err != nil || payload.Role == "" {
		logger.Logger.Error("Error binding payload: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid payload"})
	}

	var targetUser models.User
	if err := database.Db.First(&targetUser, targetUserID).Error; err != nil {
		logger.Logger.Error("Error fetching user: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
	}

	if targetUser.ID == uint(adminUserId.(float64)) {
		logger.Logger.Warn("Admin attempting to change own role")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot change your own role"})
	}

	if payload.Role == "admin" {
		targetUser.IsAdmin = true
	} else {
		targetUser.IsAdmin = false
	}

	if err := database.Db.Save(&targetUser).Error; err != nil {
		logger.Logger.Error("Error updating user role: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating user role"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Role updated successfully"})
}
