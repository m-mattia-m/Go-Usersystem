package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var users []User

func main() {
	// users = append(users, newUser("Admin", "Admin", "admin", "admin@mattiamueggler.ch", "asdfasdf", "admin"))

	r := gin.Default()
	r.POST("/registration", registration)
	r.GET("/getUsers", basicAuth, getUsers)
	r.GET("/deleteUser/:id", basicAuth, getUser)
	r.GET("/deleteUser", basicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /deleteUser/id")
	})
	r.GET("/getUser/:id", basicAuth, getUser)
	r.GET("/getUser", basicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
	r.POST("/editUser/:id", basicAuth, editUser)
	r.POST("/editUser", basicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
	r.GET("/test", test)
	r.Run(":3000")
}

func basicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	fmt.Println(user, password, hasAuth)
	fmt.Println(users)
	if hasAuth {
		successLogin := false
		for _, currentUser := range users {
			fmt.Println(currentUser)
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
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

func deleteUser(c *gin.Context) {
	id := c.Param("id")
	checkUser := false

	for i, user := range users {
		if user.Id == id {
			users = append(users[:i], users[i+1:]...)
			checkUser = true
			break
		} else {
			checkUser = false
		}
	}

	if checkUser {
		c.JSON(200, gin.H{"message": "delete user with the id: " + id})
	} else {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
}

func getUser(c *gin.Context) {
	id := c.Param("id")
	checkUser := false

	for _, user := range users {
		if user.Id == id {
			c.JSON(200, user)
			checkUser = true
			break
		} else {
			checkUser = false
		}
	}

	if !checkUser {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
}

func editUser(c *gin.Context) {
	id := c.Param("id")
	checkUser := false
	for i, user := range users {
		if user.Id == id {
			users[i].Firstname = c.PostForm("firstname")
			users[i].Lastname = c.PostForm("lastname")
			users[i].Username = c.PostForm("username")
			users[i].Email = c.PostForm("email")
			users[i].Password, _ = HashPassword(c.PostForm("password"))
			users[i].Role = c.PostForm("role")
			c.JSON(200, users[i])
			checkUser = true
			break
		} else {
			checkUser = false
		}
	}
	if !checkUser {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
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
