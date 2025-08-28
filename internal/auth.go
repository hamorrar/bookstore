package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func createToken(username string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "bookstore",                      // Issuer
		"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	fmt.Printf("Token claims added: %+v\n", claims)

	tokenString, err := claims.SignedString(os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getRole(username string) string {
	switch username {
	case "testuser1@gmail.com", "testuser2@gmail.com":
		return "customer"
	case "testuser3@gmail.com":
		return "admin"
	}
	return ""
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Dummy credential check
	if (username == "testuser1@gmail.com" && password == "testuser1password") || (username == "testuser3@gmail.com" && password == "testuser3password") {
		tokenString, err := createToken(username)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error creating token")
			return
		}

		fmt.Println(tokenString)

		// 	loggedInUser := username
		// 	fmt.Printf("Token created: %s\n", tokenString)
		// 	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
		// 	c.Redirect(http.StatusSeeOther, "/")
		// } else {
		// 	c.String(http.StatusUnauthorized, "Invalid credentials")
	}
}
