package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

func (app *application) createOrder(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to create an order"})
		return
	}

	var order database.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := app.models.Orders.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (app *application) getPageOfOrders(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get a page of orders"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	orders, err := app.models.Orders.GetPageOfOrders(limit, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get all orders page by page."})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (app *application) getAllOrders(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get all orders"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	var allOrders []*database.Order

	for {
		orders, err := app.models.Orders.GetPageOfOrders(limit, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get all orders page by page."})
			return
		}

		allOrders = append(allOrders, orders...)
		if len(orders) < limit {
			break
		}
		page++

	}

	c.JSON(http.StatusOK, allOrders)
}

func (app *application) getOrder(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get an order"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// user = app.GetUserFromContext(c)
	order, err := app.models.Orders.GetOrder(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get this user's order"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found to get"})
		return
	}

	if order.User_Id != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get this order"})
		return
	}

	c.JSON(http.StatusOK, order)

}

func (app *application) deleteOrder(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to create an order"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// user := app.GetUserFromContext(c)
	existingOrder, err := app.models.Orders.GetOrder(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order to delete"})
		return
	}

	if existingOrder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found to delete"})
		return
	}

	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete this order"})
		return
	}

	if err := app.models.Orders.DeleteOrder(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (app *application) updateOrder(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Customer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update an order"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	// user := app.GetUserFromContext(c)
	existingOrder, err := app.models.Orders.GetOrder(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order to update"})
		return
	}

	if existingOrder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found to update"})
		return
	}

	if existingOrder.User_Id != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update this order"})
		return
	}

	updatedOrder := &database.Order{}
	if err := c.ShouldBindJSON(updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrder.Id = id

	if err := app.models.Orders.UpdateOrder(updatedOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedOrder)
}
