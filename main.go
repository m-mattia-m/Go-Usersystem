package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var users []User
var accounts = make(gin.Accounts)

func main() {
	users = append(users, newUser("Admin", "Admin", "admin", "admin@mattiamueggler.ch", "Mattia12345!", "admin"))

	// authorized := r.Group("/admin", gin.BasicAuth(accounts))
	r := gin.Default()
	// authorized := initUsers(r)
	// authorized := authorizeRequest(r)

	r.POST("/registration", registration)
	r.GET("/getUsers", basicAuth, getUsers)
	r.GET("/test", test)
	r.Run(":3000")
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == "testuser" && password == "testpass" {
		// log.WithFields(log.Fields{
		// 	"user": user,
		// }).Info("User authenticated")
		fmt.Println("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		return
	}
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Hello World"})
}

// func authorizeRequest(r *gin.Engine) *gin.RouterGroup {
// 	authorized := r.Group("/admin", gin.BasicAuth(accounts), func(c *gin.Context) {
// 		fmt.Println(accounts)
// 		// username := c.MustGet(gin.AuthUserKey).(string)
// 		// fmt.Println(username)
// 	})
// 	return authorized
// }

func initUsers(r *gin.Engine) { // *gin.RouterGroup
	// accounts = make(gin.Accounts)
	for _, user := range users {
		accounts[user.Username] = user.Password
	}
	fmt.Println(accounts)
	gin.BasicAuth(accounts)
	// return authorizeRequest(r)
}

func registration(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	role := c.PostForm("role")

	currentUser := newUser(firstname, lastname, username, email, password, role)
	users = append(users, currentUser)
	c.JSON(200, currentUser)
}

func getUsers(c *gin.Context) {
	// username := c.MustGet(gin.AuthUserKey).(string)
	// password := c.MustGet(gin.AuthUserKey).(string)
	// fmt.Println(username + ", " + password)

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
