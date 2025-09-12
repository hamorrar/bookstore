package main

import (
	"os"
	"strconv"

	"github.com/hamorrar/bookstore/internal/testutils"
)

func SetupTest() *application {
	models := testutils.SetupDB()
	server_Port, _ := strconv.Atoi(os.Getenv("PORT"))
	app := &application{
		port:      server_Port,
		jwtSecret: os.Getenv("SECRET_KEY"),
		models:    models,
	}
	return app
}
