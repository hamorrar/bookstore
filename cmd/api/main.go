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

	initDBEnv()
	psqlInfo := os.Getenv("PSQL_INFO")
	fmt.Println("main api psqlinfo: ", psqlInfo)

	server_Port, _ := strconv.Atoi(os.Getenv("PORT"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

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

func initDBEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME := getDBEnv()

	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD)
	os.Setenv("PSQL_INFO", psqlInfo)
}

func getDBEnv() (string, string, string, string, string) {
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	return DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

}
