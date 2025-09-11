package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	router := gin.Default()
	router.GET("/ping", ping)
	req, _ := http.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	expectedResp := `{"status":200,"version":"v1"}`

	got := testutils.StringToJSON(w.Body.String())
	want := testutils.StringToJSON(expectedResp)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, w.Code)
}
