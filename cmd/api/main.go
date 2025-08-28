package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"

	"github.com/hamorrar/bookstore/internal/database"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {

	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME := initEnv()

	db := connectDB(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	defer db.Close()

	server_Port, _ := strconv.Atoi(os.Getenv("PORT"))

	models := database.NewModels(db)
	app := &application{
		port:      server_Port,
		jwtSecret: os.Getenv("SECRET_KEY"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

func initEnv() (string, string, string, string, string) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	return DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
}

func connectDB(DB_HOST string, DB_PORT string, DB_USER string, DB_PASSWORD string, DB_NAME string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
