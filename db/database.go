package db

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ledger-app/config"
	"log"
)

var Db *gorm.DB

func Connect() {
	cfg := config.LoadEnvironment()

	var err error

	Db, err = gorm.Open(mysql.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	log.Println("Connected to the database with GORM")
}
