package database

import (
	"database/sql"
	"fmt"
	"log"
	"varaden/server/config"
)

func InitDatabase(cfg config.DBConfig) (*sql.DB, error) {
	DB, err := connectPostgresql(cfg)

	if err != nil {
		fmt.Println("error on database")
		return nil, err
	}

	return DB, nil
}

func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing database connection: %v", err)
	} else {
		log.Fatalf("Database connection closed successfully")
	}
}
