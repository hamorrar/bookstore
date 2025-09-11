package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/testutils"
	"github.com/stretchr/testify/assert"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
)

func TestRegister(t *testing.T) {
	app := SetupTest()
	router := gin.Default()
	router.POST("/api/v1/auth/register", app.registerUser)

	url := BASE_URL + "/api/v1"

	payload := `{"email":"user1@gmail.com", "password":"password1", "role":"Customer"}`
	req, _ := http.NewRequest("POST", url+"/auth/register", strings.NewReader(payload))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedUser1Resp := `{"id":1,"email":"user1@gmail.com","role":"Customer"}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expectedUser1Resp)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, w.Code)
}
