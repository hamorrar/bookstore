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

func TestCreateOrder(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Pending","total_price":1}`

	resp, err := client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1,"user_id":1, "status":"Pending","total_price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetOrder(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.GET("/orders/:id", app.getOrder)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	// only customer can make an order
	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")
	testutils.MakeAnOrder(client, ts.URL+"/api/v1")

	// customer gets the order
	resp, err := client.Get(ts.URL + "/api/v1/orders/1")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	expected := `{"id":1,"user_id":1, "status":"Pending","total_price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateOrder(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.PUT("/orders/:id", app.updateOrder)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")
	testutils.MakeAnOrder(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Sold","total_price":1}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/orders/1", strings.NewReader(payload))
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

	expected := `{"id":1,"user_id":1, "status":"Sold","total_price":1}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateOrder_Wrong_Role(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.PUT("/orders/:id", app.updateOrder)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")
	testutils.MakeAnOrder(client, ts.URL+"/api/v1")

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Sold","total_price":1}`
	req, err := http.NewRequest(http.MethodPut, ts.URL+"/api/v1/orders/1", strings.NewReader(payload))
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

	expected := `{"error":"Unauthorized to update an order"}`
	got := testutils.StringToJSON(string(bodyBytes))
	want := testutils.StringToJSON(expected)

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestGetPageofOrders(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.GET("/orders", app.getPageOfOrders)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Pending","total_price":1}`
	_, err := client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":2}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":3}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":4}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v1/orders/")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var got []database.Order
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	expected := `[{"id":1, "user_id":1, "status":"Pending","total_price":1}, {"id":2, "user_id":1, "status":"Pending","total_price":2}]`

	var want []database.Order
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetOrder_Params(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.GET("/orders", app.getPageOfOrders)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Pending","total_price":1}`
	_, err := client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":2}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":3}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":4}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v1/orders/?page=2&limit=2")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var got []database.Order
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	expected := `[{"id":3, "user_id":1, "status":"Pending","total_price":3}, {"id":4, "user_id":1, "status":"Pending","total_price":4}]`

	var want []database.Order
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteOrder(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)
	authGroup.DELETE("/orders/:id", app.deleteOrder)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")
	testutils.MakeAnOrder(client, ts.URL+"/api/v1")

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/orders/1", nil)
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

func TestGetAllOrders(t *testing.T) {
	app := SetupTest()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	v1.POST("/auth/register", app.registerUser)
	v1.POST("/auth/login", app.login)

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.POST("/orders", app.createOrder)

	v2 := router.Group("/api/v2")
	authGroup = v2.Group("/")
	authGroup.Use(app.AuthMiddleware())
	authGroup.GET("/orders/all", app.getAllOrders)

	ts := httptest.NewServer(router)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	testutils.RegisterCustomer(client, ts.URL+"/api/v1")
	testutils.LoginCustomer(client, ts.URL+"/api/v1")

	payload := `{"user_id":1, "status":"Pending","total_price":1}`
	_, err := client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":2}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":3}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	payload = `{"user_id":1, "status":"Pending","total_price":4}`
	_, err = client.Post(ts.URL+"/api/v1/orders", "application/json", strings.NewReader(payload))
	if err != nil {
		log.Fatal(err.Error())
	}

	testutils.RegisterAdmin(client, ts.URL+"/api/v1")
	testutils.LoginAdmin(client, ts.URL+"/api/v1")

	resp, err := client.Get(ts.URL + "/api/v2/orders/all")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ts.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close()

	var got []database.Order
	if err := json.Unmarshal(bodyBytes, &got); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	expected := `[{"id":1, "user_id":1, "status":"Pending","total_price":1}, {"id":2, "user_id":1, "status":"Pending","total_price":2}, {"id":3, "user_id":1, "status":"Pending","total_price":3}, {"id":4, "user_id":1, "status":"Pending","total_price":4}]`

	var want []database.Order
	if err := json.Unmarshal([]byte(expected), &want); err != nil {
		fmt.Println("unmarshalling error while test getting page", err.Error())
	}

	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
