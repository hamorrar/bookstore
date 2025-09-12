package main

import (
	"net/http"
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

	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := `{"id":1,"email":"user1@gmail.com","role":"Customer"}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_One_Admin(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	payload := `{"email":"user2@gmail.com", "password":"password2", "role":"Admin"}`
	req, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := `{"id":1,"email":"user2@gmail.com","role":"Admin"}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegister_Two_Customers_Same(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	payload1 := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req1, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload1))

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	expected1 := `{"id":1,"email":"user1@gmail.com","role":"Customer"}`

	got1 := testutils.StringToJSON(w1.Body.String())
	want1 := testutils.StringToJSON(expected1)

	assert.Equal(t, want1, got1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// ----

	payload2 := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req2, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload2))

	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	expected2 := `{"error":"could not create registered user"}`

	got2 := testutils.StringToJSON(w2.Body.String())
	want2 := testutils.StringToJSON(expected2)

	assert.Equal(t, want2, got2)
	assert.Equal(t, http.StatusInternalServerError, w2.Code)

}

func TestLogin(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	// Register customer
	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Register admin
	payload = `{"email":"user2@gmail.com", "password":"password2", "role":"Admin"}`
	req, _ = http.NewRequest("POST", BASE_URL_v1+"/auth/register", strings.NewReader(payload))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Login customer
	payload = `{"email":"user1@gmail.com", "password":"password1"}`
	req, _ = http.NewRequest("POST", BASE_URL_v1+"/auth/login", strings.NewReader(payload))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := `{"userId":1}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, w.Code)

	// Login admin
	payload = `{"email":"user2@gmail.com", "password":"password2"}`
	req, _ = http.NewRequest("POST", BASE_URL_v1+"/auth/login", strings.NewReader(payload))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected = `{"userId":2}`

	got = testutils.StringToJSON(w.Body.String())
	want = testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogin_Wrong_Password(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)
	router.POST("/api/v1/auth/login", app.login)

	_ = testutils.RegisterCustomer(router, BASE_URL_v1+"/auth/register")
	_ = testutils.RegisterAdmin(router, BASE_URL_v1+"/auth/register")

	// Login customer
	payload := `{"email":"user1@gmail.com", "password":"password111"}`
	req, _ := http.NewRequest("POST", BASE_URL_v1+"/auth/login", strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected := `{"error": "Invalid email or password"}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Login admin
	payload = `{"email":"user2@gmail.com", "password":"password222"}`
	req, _ = http.NewRequest("POST", BASE_URL_v1+"/auth/login", strings.NewReader(payload))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expected = `{"error": "Invalid email or password"}`

	got = testutils.StringToJSON(w.Body.String())
	want = testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
