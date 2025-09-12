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

func TestGetUser(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	authGroup := router.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.GET("/api/v1/users/:id", app.getUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v1/users/1")
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected := `{"id":1, "email":"user2@gmail.com", "role":"Admin"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAllUsers(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	v2 := router.Group("/api/v2")
	authGroup := v2.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.GET("/users/all", app.getAllUsers)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	// Make a few test customers to get
	payload := `{"email":"user3@gmail.com", "password":"password1", "role":"Customer"}`
	client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))

	payload = `{"email":"user4@gmail.com", "password":"password1", "role":"Customer"}`
	client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))

	payload = `{"email":"user5@gmail.com", "password":"password1", "role":"Customer"}`
	client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))

	resp, err := client.Get(ts.URL + "/api/v2/users/all")
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	var got []database.User
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		log.Fatalf("pain: %v", err.Error())
	}

	expected := `[{"id":1,"email":"user2@gmail.com","role":"Admin"},{"id":2,"email":"user3@gmail.com","role":"Customer"},{"id":3,"email":"user4@gmail.com","role":"Customer"},{"id":4,"email":"user5@gmail.com","role":"Customer"}]`
	var want []database.User
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		log.Fatalf("pain 2: %v", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetUser_Wrong_Role(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	authGroup := router.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.GET("/api/v1/users/:id", app.getUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v1/users/1")
	if err != nil {
		log.Fatalln("fatal request", err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected := `{"error":"Unauthorized to get users"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestUpdateUser(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	authGroup := router.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.PUT("/api/v1/users/:id", app.updateUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	payload := `{"email":"user22@gmail.com", "password":"password2", "role":"Admin"}`

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/users/1", strings.NewReader(payload))
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

	expected := `{"id":1, "email":"user22@gmail.com", "role":"Admin"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateUser_Wrong_Role(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	authGroup := router.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.PUT("/api/v1/users/:id", app.updateUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"email":"user22@gmail.com", "password":"password2", "role":"Customer"}`

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/users/1", strings.NewReader(payload))
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

	expected := `{"error":"Unauthorized to update user"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestDeleteUser(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	authGroup := router.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.DELETE("/api/v1/users/:id", app.deleteUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/users/1", nil)
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
