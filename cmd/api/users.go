package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
)

func (app *application) getAllUsers(c *gin.Context) {
	user := app.GetUserFromContext(c)
	if user.Role != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized to get all users"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

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
