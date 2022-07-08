package users

import (
	"fmt"
	"log"
	"sort"

	"usersystem/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// var users []User = getUsersFromDB()

func BasicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	var users []User = getUsersFromDB()
	user, password, hasAuth := c.Request.BasicAuth()
	fmt.Println(user, password, hasAuth)
	fmt.Println(users)
	if hasAuth {
		i := sort.Search(len(users), func(i int) bool { return user <= users[i].Username })
		if i < len(users) && users[i].Username == user {
			if checkPasswordHash(password, users[i].Password) {
				// c.JSON(200, gin.H{"message": "user found"})
				fmt.Println("successfully")
			} else {
				// c.JSON(400, gin.H{"error": "password is not correct"})
				c.Abort()
				c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				c.JSON(401, gin.H{"error": "unauthorized"})
			}
		} else {
			c.Abort()
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.JSON(400, gin.H{"error": "user not found"})
		}
	}
}

func Registration(c *gin.Context) {
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password, _ := hashPassword(c.PostForm("password"))
	role := c.PostForm("role")
	currentUser, msg := newUser(firstname, lastname, username, email, password, role)
	if msg == "" {
		// users = append(users, currentUser)
		saveUserOnDB(currentUser)
		c.JSON(200, currentUser)
	} else {
		c.JSON(400, gin.H{"error": msg})
	}

}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUsers(c *gin.Context) {
	var users []User = getUsersFromDB()
	if users != nil {
		c.JSON(200, users)
	} else {
		c.JSON(400, gin.H{"error": "No users found"})
	}
}

func newUser(firstname string, lastname string, username string, email string, password string, role string) (User, string) {
	var users []User = getUsersFromDB()
	i := sort.Search(len(users), func(i int) bool { return username <= users[i].Username })
	if i < len(users) && users[i].Username == username {
		user := User{}
		return user, "username already exists"
	}
	i = sort.Search(len(users), func(i int) bool { return email <= users[i].Email })
	if i < len(users) && users[i].Email == email {
		user := User{}
		return user, "email already exists"
	}

	var currentUser = new(User)
	currentUser.Id = uuid.New().String()
	currentUser.Firstname = firstname
	currentUser.Lastname = lastname
	currentUser.Username = username
	currentUser.Email = email
	currentUser.Password = password
	currentUser.Role = role

	return *currentUser, ""
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User = getUserFromDBById(id)

	if user.Id == id {
		deleteUserFromDB(user)
		c.JSON(200, gin.H{"message": "delete user with the id: " + id})
	} else {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user User = getUserFromDBById(id)

	if user.Id == id {
		c.JSON(200, user)
	} else {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
}

func EditUser(c *gin.Context) {
	id := c.Param("id")
	var user User = getUserFromDBById(id)

	if user.Id == id {
		user.Firstname = c.PostForm("firstname")
		user.Lastname = c.PostForm("lastname")
		user.Username = c.PostForm("username")
		user.Email = c.PostForm("email")
		user.Password, _ = hashPassword(c.PostForm("password"))
		user.Role = c.PostForm("role")
		updateUserOnDB(user)
		c.JSON(200, user)
	} else {
		c.JSON(400, gin.H{"error": "No user found with the id: " + id})
	}
}

func Login(c *gin.Context) {
	fmt.Println("login")
	// user, password, hasAuth := c.Request.BasicAuth()
}

// func generateJwtToken(c *gin.Context) {
// 	fmt.Println("jwtToken")
// 	user, password, hasAuth := c.Request.BasicAuth()
// 	jwt.Auth(SECRET)
// }

func getUsersFromDB() []User {
	var users []User
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role` FROM `users`")
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Email, &user.Password, &user.Role)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", user)
		users = append(users, user)
	}
	rows.Close()
	return users
}

func getUserFromDBById(id string) User {
	var users []User
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role` FROM `users` WHERE Id=`" + id + "`")
	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Email, &user.Password, &user.Role)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", user)
		users = append(users, user)
	}
	rows.Close()
	return users[0]
}

func saveUserOnDB(user User) {
	var query string = "INSERT INTO users (`Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role`) VALUES ('" + user.Id + "', '" + user.Firstname + "', '" + user.Lastname + "', '" + user.Username + "', '" + user.Email + "', '" + user.Password + "', '" + user.Role + "');"
	db.RunSqlQueryWithoutReturn(query)
}

func updateUserOnDB(user User) {
	var query string = "UPDATE users SET `Firstname`='" + user.Firstname + "', `Lastname`='" + user.Lastname + "', `Username`='" + user.Username + "', `Email`='" + user.Email + "', `Password`='" + user.Password + "', `Role`='" + user.Role + "' WHERE Id='" + user.Id + "';"
	db.RunSqlQueryWithoutReturn(query)
}

func deleteUserFromDB(user User) {
	var query string = "DELETE FROM `users` WHERE Id=" + user.Id + ";"
	db.RunSqlQueryWithoutReturn(query)
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
