package providers

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"ledger-app/config"
	"ledger-app/internal/connections/database"
	"ledger-app/internal/middleware"
	"ledger-app/logger"
	"ledger-app/models"
	"ledger-app/routes"
)

func InitLogger() {
	logger.InitLogger()
}

func InitDatabase() {
	database.Connect()
}

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middleware.LogRequest)
	routes.RegisterRoutes(e)
}

func InitDefaultAdmin() {
	var count int64
	if err := database.Db.Model(&models.User{}).Where("is_admin = ?", true).Count(&count).Error; err != nil {
		logger.Logger.Fatalf("Failed to query database for admin user: %v", err)
	}

	if count > 0 {
		logger.Logger.Infof("Default admin already existing. Skipping creation.")
		return
	}

	adminConfig := config.LoadEnvironment()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminConfig.DefaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Fatalf("Failed to hash default admin password: %v", err)
	}

	admin := models.User{
		Name:         adminConfig.DefaultAdminUserName,
		PasswordHash: string(hashedPassword),
		IsAdmin:      true,
	}

	if err := database.Db.Create(&admin).Error; err != nil {
		logger.Logger.Fatalf("Failed to create default admin user: %v", err)
	}

	logger.Logger.Infof("Default admin created with username: %s and password: %s\n", adminConfig.DefaultAdminUserName, adminConfig.DefaultAdminPassword)
}

func StartServer(e *echo.Echo, cfg *config.Config) {
	logger.Logger.Infof("Starting server at port %s", cfg.Port)

	if err := e.Start(":" + cfg.Port); err != nil {
		logger.Logger.Fatal("Error starting server", err)
	}
}
