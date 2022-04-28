package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var users []User
var accounts = make(gin.Accounts)

func main() {
	users = append(users, newUser("Admin", "Admin", "admin", "admin@mattiamueggler.ch", "asdfasdf", "admin"))

	r := gin.Default()
	r.POST("/registration", registration)
	r.GET("/getUsers", basicAuth, getUsers)
	r.GET("/test", test)
	r.Run(":3000")
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()

	if hasAuth {
		successLogin := false
		for _, currentUser := range users {
			if user == currentUser.Username && CheckPasswordHash(password, currentUser.Password) {
				// c.JSON(200, gin.H{"message": "You are authenticated"})
				fmt.Println("User authenticated")
				successLogin = true
				break
			} else {
				successLogin = false
			}
		}
		if !successLogin {
			c.Abort()
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.JSON(401, gin.H{"error": "unauthorized"})
		}
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(401, gin.H{"error": "has no login"})
	}
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello World"})
}

func registration(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password, _ := HashPassword(c.PostForm("password"))
	role := c.PostForm("role")

	currentUser := newUser(firstname, lastname, username, email, password, role)
	users = append(users, currentUser)
	c.JSON(200, currentUser)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUsers(c *gin.Context) {
	if users != nil {
		c.JSON(200, users)
	} else {
		c.JSON(400, gin.H{"error": "No users found"})
	}
}

func newUser(firstname string, lastname string, username string, email string, password string, role string) User {
	var currentUser = new(User)
	currentUser.Id = uuid.New().String()
	currentUser.Firstname = firstname
	currentUser.Lastname = lastname
	currentUser.Username = username
	currentUser.Email = email
	currentUser.Password = password
	currentUser.Role = role

	return *currentUser
}

type User struct {
	Id        string
	Firstname string
	Lastname  string
	Username  string
	Email     string
	Password  string
	Role      string
}
