package testutils

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
)

func SetupDB() database.Models {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	DB_DSN := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../migrate/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	models := database.NewModels(db)
	return models
}
func StringToJSON(str string) map[string]interface{} {
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(str), &res); err != nil {
		log.Fatalf("could not unmarshal expected: %v", err.Error())
	}
	return res
}

func RegisterCustomer(router *gin.Engine, url string) *httptest.ResponseRecorder {
	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func RegisterAdmin(router *gin.Engine, url string) *httptest.ResponseRecorder {
	payload := `{"email":"user2@gmail.com", "password":"password2", "role":"Admin"}`
	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}
