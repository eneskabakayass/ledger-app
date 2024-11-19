package database

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ledger-app/config"
	"ledger-app/logger"
	"ledger-app/models"
)

var Db *gorm.DB

func Connect() {
	cfg := config.LoadEnvironment()

	var err error

	Db, err = gorm.Open(mysql.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		logger.Logger.Fatal("Error connecting to the database:", err)
	}

	err = Db.AutoMigrate(&models.User{}, &models.Transaction{})
	if err != nil {
		logger.Logger.Fatal("Error migrate the database", err)
	}

	logger.Logger.Infof("Connected to the database with GORM")
}
