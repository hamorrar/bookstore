package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hamorrar/bookstore/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"`
}

// RegisterUser registers a new user
// @Summary		registers a new user
// @Description	registers a new user
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		user body registerRequest true "user registration info"
// @Success		201	{object} database.User
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 500 {object} gin.H "error generating password"
// @Failure 500 {object} gin.H "error creating user"
// @Router		/api/v1/auth/register [post]
func (app *application) registerUser(c *gin.Context) {
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate password"})
		return
	}

	register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: register.Password,
		Role:     register.Role,
	}

	err = app.models.Users.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create registered user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// Login logins in a user
// @Summary		logins a user
// @Description	logins a user
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		user body loginRequest true "user login info"
// @Success		200	{object} gin.H "Successfully logged in user"
// @Failure 400 {object} gin.H "error binding JSON"
// @Failure 401 {object} gin.H "error finding user with email"
// @Failure 500 {object} gin.H "error getting user by email"
// @Failure 401 {object} gin.H "invalid email or password"
// @Failure 500 {object} gin.H "error generating token"
// @Router			/api/v1/auth/login [post]
func (app *application) login(c *gin.Context) {
	var auth loginRequest
	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetUserByEmail(auth.Email)
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found with this email"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user to login"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": existingUser.Id,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.SetCookie("auth_token", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"userId": existingUser.Id})
}
