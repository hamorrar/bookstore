package main

import (
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

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	app := setupApp()

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

func setupApp() *application {
	server_Port, _ := strconv.Atoi(os.Getenv("PORT"))
	DB_DSN := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

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

	return app
}
