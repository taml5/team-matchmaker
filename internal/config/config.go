package config

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/ncruces/go-sqlite3/driver"

	_ "github.com/ncruces/go-sqlite3/embed"
)

var (
	RiotAPIKey string
	QueueID    int = 0 // queue ID for custom games
	DB         *sql.DB
)

func LoadEnv() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	RiotAPIKey = os.Getenv("RIOT_API_KEY")
}

func LoadDB() error {
	// Connect to the database
	db, err := sql.Open("sqlite3", "./data/data.db")
	if err != nil {
		return err
	}

	// Check if the database is reachable
	if err = db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}
