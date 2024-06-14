package database

import (
	"database/sql"
	"fmt"
	"os"

	"log"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func getDatabaseConnectionString() string {
	connString := "%s?_foreign_keys=on"
	return fmt.Sprintf(connString, os.Getenv("SQLITE_DATABASE_CONNECTION_STRING"))
}

func NewDatabaseConnection() *Database {
	d := &Database{}

	db, err := sql.Open("sqlite3", getDatabaseConnectionString())

	if err != nil {
		log.Fatalf("Error opening database connection: %s", err)
	}

	d.db = db

	return d
}

func (d *Database) GetConnection() *sql.DB {
	return d.db
}

func (d *Database) CloseConnection() {
	d.db.Close()
}

func (d *Database) Ping() {
	err := d.db.Ping()

	if err != nil {
		log.Fatalf("Error pinging database connection: %s", err)
	}

	slog.Info("Database connection pinged successfully")
}
