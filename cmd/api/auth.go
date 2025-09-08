package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hamorrar/bookstore/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// type loginRequest struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=4"`
// }

// type loginResponse struct {
// 	Token  string `json:"token"`
// 	UserId int    `json:"userId"`
// }

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"`
}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not created registered user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// func createToken(username string) (string, error) {
// 	// Create a new JWT token with claims
// 	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"sub": username,                         // Subject (user identifier)
// 		"iss": "bookstore",                      // Issuer
// 		"aud": getRole(username),                // Audience (user role)
// 		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
// 		"iat": time.Now().Unix(),                // Issued at
// 	})

// 	fmt.Printf("Token claims added: %+v\n", claims)

// 	tokenString, err := claims.SignedString(os.Getenv("SECRET_KEY"))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func getRole(username string) string {
// 	switch username {
// 	case "testuser1@gmail.com", "testuser2@gmail.com":
// 		return "customer"
// 	case "testuser3@gmail.com":
// 		return "admin"
// 	}
// 	return ""
// }

// func Login(c *gin.Context) {
// 	username := c.PostForm("username")
// 	password := c.PostForm("password")

// 	// Dummy credential check
// 	if (username == "testuser1@gmail.com" && password == "testuser1password") || (username == "testuser3@gmail.com" && password == "testuser3password") {
// 		tokenString, err := createToken(username)
// 		if err != nil {
// 			c.String(http.StatusInternalServerError, "Error creating token")
// 			return
// 		}

// 		fmt.Println(tokenString)

// 		// 	loggedInUser := username
// 		// 	fmt.Printf("Token created: %s\n", tokenString)
// 		// 	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
// 		// 	c.Redirect(http.StatusSeeOther, "/")
// 		// } else {
// 		// 	c.String(http.StatusUnauthorized, "Invalid credentials")
// 	}
// }
