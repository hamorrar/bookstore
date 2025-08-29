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
	psqlInfo := os.Getenv("PSQL_INFO")
	psqlURL := os.Getenv("DB_URL")

	fmt.Println("in migrate main")
	fmt.Println("psql info: ", psqlInfo)
	fmt.Println("psql url: ", psqlURL)

	// psqlURL2 := "postgres://localhost:5432/database?sslmode=disable"
	// db, err := sql.Open("postgres", psqlInfo)
	db, err := sql.Open("postgres", psqlURL)
	// db, err := sql.Open("postgres", psqlURL2)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("--PING ERROR")
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// fSrc, err := (&file.File{}).Open("cmd/migrate/migrations")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("--fsrc:", fSrc)

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)
	fmt.Println("--m:", m)

	if err != nil {
		log.Fatal(err)
	}

	// m.Up()

	direction := os.Args[1]
	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("UP ERR")
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction. Use 'up' or 'down'.")
	}
	// defer db.Close()

}

func initDBEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME := getDBEnv()

	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD)
	os.Setenv("PSQL_INFO", psqlInfo)

	os.Setenv("DB_URL", "postgres://postgres:postgres@localhost:5432/bookstore?sslmode=disable")
}

func getDBEnv() (string, string, string, string, string) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	return DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

}
