package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/hamorrar/bookstore/internal/testutils"
	"github.com/joho/godotenv"
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

func TestMain(m *testing.M) {
	rootPath, _ := os.Getwd()
	for rootPath != "/" {
		envPath := filepath.Join(rootPath, ".env")
		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				log.Fatalf("Error loading .env file: %v", err)
			}
			break
		}
		rootPath = filepath.Dir(rootPath)
	}
	code := m.Run()
	os.Exit(code)
}
