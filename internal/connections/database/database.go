package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ledger-app/config"
)

var Db *gorm.DB

func Connect() {
	cfg := config.LoadEnvironment()

	var err error

	Db, err = gorm.Open(mysql.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Error connecting to the database:", err)
	}

	logrus.Info("Connected to the database with GORM")
}
