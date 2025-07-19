package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stealcash/AgentFlow/app/globals"
	"log"
)

var DB *sql.DB

func MainConnection() error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		globals.Config.MainDatabase.Host,
		globals.Config.MainDatabase.Port,
		globals.Config.MainDatabase.User,
		globals.Config.MainDatabase.Password,
		globals.Config.MainDatabase.Name,
		globals.Config.MainDatabase.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	return nil
}

func CloseMainDb() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf(" Error closing main DB: %v", err)
		} else {
			log.Println(" Main DB connection closed")
		}
	}
}
