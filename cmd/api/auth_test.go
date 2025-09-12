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
	"github.com/hamorrar/bookstore/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRegister_One_Customer(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`

	resp, err := client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1,"email":"user1@gmail.com","role":"Customer"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestRegister_One_Admin(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	payload := `{"email":"user2@gmail.com", "password":"password2", "role":"Admin"}`

	resp, err := client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1,"email":"user2@gmail.com","role":"Admin"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestRegister_Two_Customers_Same(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`

	resp, err := client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1,"email":"user1@gmail.com","role":"Customer"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// ----

	payload = `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`

	resp, err = client.Post(ts.URL+"/api/v1/auth/register", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected = `{"error": "could not create registered user"}`
	got = testutils.StringToJSON(string(bodyBytes))
	want = testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestLogin(t *testing.T) {

	app := SetupTest()
	router := gin.Default()
	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.RegisterAdmin(client, ts.URL+"/api/v1")

	// Login customer
	payload := `{"email":"user1@gmail.com", "password":"password1"}`

	resp, err := client.Post(ts.URL+"/api/v1/auth/login", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"userId":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Login admin
	payload = `{"email":"user2@gmail.com", "password":"password2"}`

	resp, err = client.Post(ts.URL+"/api/v1/auth/login", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected = `{"userId":2}`
	got = testutils.StringToJSON(string(bodyBytes))
	want = testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLogin_Wrong_Password(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	ts := httptest.NewServer(router)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.RegisterAdmin(client, ts.URL+"/api/v1")

	// Login customer
	payload := `{"email":"user1@gmail.com", "password":"password111"}`

	resp, err := client.Post(ts.URL+"/api/v1/auth/login", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"error": "Invalid email or password"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Login admin
	payload = `{"email":"user2@gmail.com", "password":"password222"}`
	resp, err = client.Post(ts.URL+"/api/v1/auth/login", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	expected = `{"error": "Invalid email or password"}`
	got = testutils.StringToJSON(string(bodyBytes))
	want = testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
