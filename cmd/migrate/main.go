package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: 'up' or 'down'.")
	}

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	createDB()

	psqlURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", psqlURL)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	direction := os.Args[1]
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
	defer db.Close()
}

func createDB() {
	// Open default "postgres" database to create bookstore database
	DEFAULT_PSQL_INFO := os.Getenv("DEFAULT_PSQL_INFO")

	db, err := sql.Open("postgres", DEFAULT_PSQL_INFO)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Check if the prod/test database exists or needs to be created
	DB_NAME := os.Getenv("DB_NAME")

	var exists bool
	query := fmt.Sprintf("select exists(select 1 from pg_database where datname = '%s')", DB_NAME)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking database existence: %v", err)
	}

	// Create the database if not exists
	if !exists {
		query := fmt.Sprintf("create database %s", DB_NAME)
		_, err = db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating database '%s': %v", DB_NAME, err)
		}
		fmt.Printf("Database '%s' created or already exists.\n", DB_NAME)
	}
}
