package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hamorrar/bookstore/src/router"
)

func main() {

	ginRouter := gin.Default()
	router.InitRouter(ginRouter)

	DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME := initEnv()

	connectDB(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	// Run the router
	ginRouter.Run()

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

func connectDB(DB_HOST string, DB_PORT string, DB_USER string, DB_PASSWORD string, DB_NAME string) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
}
