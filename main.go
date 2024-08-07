package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


var secretKey = []byte("your-secret-key")

func createToken(username string) (string, error) {
	
	    // Create a new JWT token with claims
		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": username,                    // Subject (user identifier)
			"iss": "todo-app",                  // Issuer
			"aud": getRole(username),           // Audience (user role)
			"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
			"iat": time.Now().Unix(),                 // Issued at
		})
	
		tokenString, err := claims.SignedString(secretKey)
		if err != nil {
			return "", err
		}
	  // Print information about the created token
		fmt.Printf("Token claims added: %+v\n", claims)
		return tokenString, nil
}

func getRole(username string) string {
	if username == "senior" {
		return "senior"
	}
	return "employee"
}

type Todo struct {
	Text string
	Done bool
}

var todos []Todo
var loggedInUser string

func main() {
	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Todos":    todos,
			"LoggedIn": loggedInUser != "",
			"Username": loggedInUser,
		})
	})

	router.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Dummy credential check
		if (username == "employee" && password == "password") || (username == "senior" && password == "password") {
			tokenString, err := createToken(username)
			if err != nil {
				c.String(http.StatusInternalServerError, "Error creating token")
				return
			}

			loggedInUser = username
			fmt.Printf("Token created: %s\n", tokenString)
			c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
			c.Redirect(http.StatusSeeOther, "/")
		} else {
			c.String(http.StatusUnauthorized, "Invalid credentials")
		}
})

	router.POST("/add", func(c *gin.Context) {
		text := c.PostForm("todo")
		todo := Todo{Text: text, Done: false}
		todos = append(todos, todo)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.POST("/toggle", func(c *gin.Context) {
		index := c.PostForm("index")
		toggleIndex(index)
		c.Redirect(http.StatusSeeOther, "/")
	})

	router.Run(":8080")
}

func toggleIndex(index string) {
	i, _ := strconv.Atoi(index)
	if i >= 0 && i < len(todos) {
		todos[i].Done = !todos[i].Done
	}
}