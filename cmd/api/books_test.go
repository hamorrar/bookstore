package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4"
	"github.com/hamorrar/bookstore/internal/testutils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"title":"Title1", "author":"First","price":1}`

	resp, err := client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1, "title":"Title1", "author":"First","price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetBook(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)
	authGroup.GET("/books/:id", app.getBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	testutils.MakeABook(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v1/books/1")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1, "title":"Title1", "author":"First","price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
