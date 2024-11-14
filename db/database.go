package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"ledger-app/config"
	"log"
)

var Db *sql.DB

func Connect() {
	cfg := config.LoadConfig()
	var err error
	Db, err = sql.Open("mysql", cfg.DBUrl)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	log.Println("Connected to the database")
}
