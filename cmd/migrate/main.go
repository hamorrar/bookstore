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

	initDBEnv()
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

func initDBEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME := getDBEnv()

	DB_URL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	os.Setenv("DB_URL", DB_URL)
}

func getDBEnv() (string, string, string, string, string) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	return DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

}
