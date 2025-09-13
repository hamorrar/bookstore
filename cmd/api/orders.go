package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

// createOrder creates an order
// @Summary		creates an order
// @Description	creates an order
// @Tags		order
// @Accept		json
// @Produce		json
// @Param		order body database.Order true "new order to add to db"
// @Success		201	{object} database.Order "successfully created an order"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "error creating order"
// @Router			/api/v1/orders [post]
// @Security CookieAuth
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

// getPageOfOrders gets a page of orders
// @Summary		gets a page of orders
// @Description	gets a page of orders
// @Tags		order
// @Accept		json
// @Produce		json
// @Param		page query int false "page number to request"
// @Param limit query int false "max number of orders to return per page"
// @Success		200	{array} database.Order "successfully got a page of orders"
// @Failure 500 {object} gin.H "error getting a page"
// @Failure 403 {object} gin.H "wrong role"
// @Router			/api/v1/orders [get]
// @Security CookieAuth
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

// getAllOrders gets all orders
// @Summary		gets all orders
// @Description	gets all orders
// @Tags		order
// @Accept		json
// @Produce		json
// @Success		200	{array} database.Order "successfully got all Orders"
// @Failure 500 {object} gin.H "error getting all Orders"
// @Failure 403 {object} gin.H "wrong role"
// @Router			/api/v2/orders/all [get]
// @Security CookieAuth
func (app *application) getAllOrders(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get all orders"})
		return
	}

	limit := 3
	page := 1
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

// getOrder get one order
// @Summary		get one order
// @Description	get one order by id
// @Tags		order
// @Accept		json
// @Produce		json
// @Param		id query int true "id of order to get"
// @Success		200	{object} database.Order "successfully got an order"
// @Failure 400 {object} gin.H "invalid order id"
// @Failure 404 {object} gin.H "order not found with this id"
// @Failure 500 {object} gin.H "error getting order"
// @Failure 403 {object} gin.H "wrong role/unauthorized"
// @Router			/api/v1/orders/:id [get]
// @Security CookieAuth
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

// deleteOrder delete an order
// @Summary		delete order
// @Description	delete an order by id
// @Tags		order
// @Accept		json
// @Produce		json
// @Param		id query int true "id of order to delete"
// @Success		204	"successfully deleted"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error deleting order"
// @Router			/api/v1/orders/:id [delete]
// @Security CookieAuth
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

// updateOrder updates an order
// @Summary		update an order
// @Description	update an order by id
// @Tags		order
// @Accept		json
// @Produce		json
// @Param		id query int true "id of order to update"
// @Param order body database.Order true "updated order data"
// @Success 200	{object} database.Order "successfully updated a order"
// @Failure 403 {object} gin.H "wrong role/unauthorized"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error getting order"
// @Failure 404 {object} gin.H "order to update not found"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "failed to update order"
// @Router			/api/v1/orders/:id [put]
// @Security CookieAuth
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
