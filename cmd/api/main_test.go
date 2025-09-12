package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hamorrar/bookstore/internal/testutils"
)

var BASE_URL_v1 string = fmt.Sprintf("http://%s:%s/api/v1", os.Getenv("HOST"), os.Getenv("PORT"))

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
