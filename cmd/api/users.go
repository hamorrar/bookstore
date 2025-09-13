package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

// getAllUsers gets all users
// @Summary		gets all users
// @Description	gets all users
// @Tags		user
// @Accept		json
// @Produce		json
// @Success		200	{array} database.User "successfully got all users"
// @Failure 500 {object} gin.H "error getting all users"
// @Failure 403 {object} gin.H "wrong role"
// @Router			/api/v2/users/all [get]
// @Security CookieAuth
func (app *application) getAllUsers(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get all users"})
		return
	}

	limit := 3
	page := 1
	var allUsers []*database.User
	for {
		users, err := app.models.Users.GetPageOfUsers(limit, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get all users page by page."})
			return
		}

		allUsers = append(allUsers, users...)
		if len(users) < limit {
			break
		}
		page++
	}

	c.JSON(http.StatusOK, allUsers)
}

// getUser get one user
// @Summary		get one user
// @Description	get one user by id
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id query int true "id of user to get"
// @Success		200	{object} database.User "successfully got a user"
// @Failure 400 {object} gin.H "invalid user id"
// @Failure 404 {object} gin.H "user not found with this id"
// @Failure 500 {object} gin.H "error getting user"
// @Failure 403 {object} gin.H "wrong role/unauthorized"
// @Router			/api/v1/users/:id [get]
// @Security CookieAuth
func (app *application) getUser(c *gin.Context) {
	userCtx := app.GetUserFromContext(c)
	if userCtx.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get users"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := app.models.Users.GetUserById(id)
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// deleteUser delete a user
// @Summary		delete user
// @Description	delete a user by id
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id query int true "id of user to delete"
// @Success		204	"successfully deleted"
// @Failure 403 {object} gin.H "wrong role"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error deleting user"
// @Router			/api/v1/users/:id [delete]
// @Security CookieAuth
func (app *application) deleteUser(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to delete user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := app.models.Users.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// updateUser updates a user
// @Summary		update a user
// @Description	update a user by id
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		id query int true "id of user to update"
// @Param user body database.User true "updated user data"
// @Success 200	{object} database.User "successfully updated a user"
// @Failure 403 {object} gin.H "wrong role/unauthorized"
// @Failure 400 {object} gin.H "invalid id"
// @Failure 500 {object} gin.H "error getting user"
// @Failure 404 {object} gin.H "user to update not found"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "failed to update user"
// @Router			/api/v1/users/:id [put]
// @Security CookieAuth
func (app *application) updateUser(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to update user"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	existingUser, err := app.models.Users.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retreive user"})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updatedUser := &database.User{}
	if err := c.ShouldBindJSON(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser.Id = id

	if err := app.models.Users.UpdateUser(updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
