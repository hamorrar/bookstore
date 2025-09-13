package main

import (
	"encoding/json"
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
	"github.com/hamorrar/bookstore/internal/database"
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

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

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
	v1.GET("/books/:id", app.getBook)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// only admin can make a book
	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")
	testutils.MakeABook(client, ts.URL+"/api/v1")

	// customer gets the book
	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")
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

func TestUpdateBook(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)
	authGroup.PUT("/books/:id", app.updateBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")
	testutils.MakeABook(client, ts.URL+"/api/v1")

	payload := `{"id":1, "title":"Title11", "author":"First","price":1}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/books/1", strings.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected := `{"id":1, "title":"Title11", "author":"First","price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateBook_Wrong_Role(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)
	authGroup.PUT("/books/:id", app.updateBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")
	testutils.MakeABook(client, ts.URL+"/api/v1")

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"id":1, "title":"Title11", "author":"First","price":1}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/books/1", strings.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected := `{"error":"Unauthorized to update book"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestGetPageofBooks(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)
	v1.GET("/books", app.getPageOfBooks)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// only admin can make a book
	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	payload := `{"title":"Title1", "author":"First","price":1}`
	_, err := client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title2", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title3", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title4", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Get(ts.URL + "/api/v1/books/")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var got []database.Book
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	expected := `[{"id":1, "title":"Title1", "author":"First","price":1}, {"id":2, "title":"Title2", "author":"First","price":1}]`

	var want []database.Book
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting all", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetPage_Params(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)
	v1.GET("/books", app.getPageOfBooks)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// only admin can make a book
	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	payload := `{"title":"Title1", "author":"First","price":1}`
	_, err := client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title2", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title3", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title4", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Get(ts.URL + "/api/v1/books/?page=2&limit=2")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var got []database.Book
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	expected := `[{"id":3, "title":"Title3", "author":"First","price":1}, {"id":4, "title":"Title4", "author":"First","price":1}]`

	var want []database.Book
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting all", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteBook(t *testing.T) {

	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)
	authGroup.DELETE("/books/:id", app.deleteBook)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")
	testutils.MakeABook(client, ts.URL+"/api/v1")

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/books/1", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	got := string(bodyBytes)
	want := ""

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

}

func TestGetAllBooks(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/books", app.createBook)

	v2 := router.Group("/api/v2")
	authGroup = v2.Group("/")
	authGroup.Use(app.AuthMiddleware())

	authGroup.GET("/books/all", app.getAllBooks)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	payload := `{"title":"Title1", "author":"First","price":1}`
	_, err := client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title2", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title3", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"title":"Title4", "author":"First","price":1}`
	_, err = client.Post(ts.URL+"/api/v1/books", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Get(ts.URL + "/api/v2/books/all")
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	var got []database.Book
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting all books", err.Error())
	}

	expected := `[{"Id":1,"title":"Title1", "author":"First","price":1},{"Id":2,"title":"Title2", "author":"First","price":1},{"Id":3,"title":"Title3", "author":"First","price":1},{"Id":4,"title":"Title4", "author":"First","price":1}]`
	var want []database.Book
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting all books", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
