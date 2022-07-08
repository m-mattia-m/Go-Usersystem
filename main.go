package main

import (
	"usersystem/db"
	"usersystem/users"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db.InitDB()

	r := gin.Default()
	r.POST("/login", users.BasicAuth, users.Login)
	r.POST("/registration", users.Registration)
	r.GET("/getUsers", users.BasicAuth, users.GetUsers)
	r.GET("/deleteUser/:id", users.BasicAuth, users.GetUser)
	r.GET("/deleteUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /deleteUser/id")
	})
	r.GET("/getUser/:id", users.BasicAuth, users.GetUser)
	r.GET("/getUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
	r.POST("/editUser/:id", users.BasicAuth, users.EditUser)
	r.POST("/editUser", users.BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
	r.Run(":3000")
}
