package users

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"sort"
	"usersystem/db"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Main(r *gin.RouterGroup) {

	db.CreateUsersTable()

	r.GET("/login", Login)
	r.POST("/registration", Registration)
	r.GET("/getUsers", BasicAuth, GetUsers)
	r.GET("/deleteUser/:id", BasicAuth, GetUser)
	r.GET("/deleteUser", BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /deleteUser/id")
	})
	r.GET("/getUser/:id", BasicAuth, GetUser)
	r.GET("/getUser", BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
	r.POST("/editUser/:id", BasicAuth, EditUser)
	r.POST("/editUser", BasicAuth, func(c *gin.Context) {
		c.JSON(400, "Send an ID of a user with. Example: /getUser/id")
	})
}

func BasicAuth(c *gin.Context) {
	var users []User = getUsersFromDB()
	token := c.Request.Header.Get("token")
	userId := c.Request.Header.Get("userid")

	if len(token) > 0 && len(userId) > 0 {
		i := sort.Search(len(users), func(i int) bool { return userId <= users[i].Id })
		if i < len(users) && users[i].Id == userId {
			if users[i].Token == token {
				fmt.Println("Successfully Login with Token")
			} else {
				c.Abort()
				c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				c.JSON(401, gin.H{"error": "unauthorized"})
				c.JSON(401, gin.H{"users": users[i], "CurrentToken": token, "userToken": users[i].Token})
			}
		} else {
			c.Abort()
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.JSON(400, gin.H{"error": "user not found", "token": token, "userId": userId, "userslength": len(users), "i": i})
			c.JSON(400, gin.H{"users": users[i]})
		}
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(400, gin.H{"error": "no credentials provided"})
	}
}

func Login(c *gin.Context) {
	var users []User = getUsersFromDB()
	username, password, hasAuth := c.Request.BasicAuth()
	if hasAuth {
		i := sort.Search(len(users), func(i int) bool { return users[i].Username >= username })
		if i < len(users) && users[i].Username == username {
			if checkPasswordHash(password, users[i].Password) {
				fmt.Println("[Login]: Successfully Login with Username and Password")
				token, err := GenerateRandomStringURLSafe(128)
				if err != nil {
					c.JSON(400, gin.H{"error": "Error generating token"})
				}
				updateTokenOnDB(users[i], token)
				c.JSON(200, gin.H{"token": token, "userId": users[i].Id})
			} else {
				c.Abort()
				c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
				c.JSON(401, gin.H{"error": "unauthorized"})
			}
		} else {
			c.Abort()
			c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			c.JSON(400, gin.H{"error": "user not found"})
		}
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(400, gin.H{"error": "no credentials provided"})
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

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}
	return string(ret), nil
}

func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func getUsersFromDB() []User {
	var users []User
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role`, `Token` FROM `users`")
	if err != nil {
		fmt.Println("[DB]: Can't Select users from DB \t-->\t" + err.Error())
	}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Email, &user.Password, &user.Role, &user.Token)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to User-Object (userId: %v) \t-->\t"+err.Error(), user.Id)
		}
		users = append(users, user)
	}
	rows.Close()
	return users
}

func getUserFromDBById(id string) User {
	var users []User
	rows, err := db.RunSqlQueryWithReturn("SELECT `Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role`, `Token` FROM `users` WHERE Id=`" + id + "`")
	if err != nil {
		fmt.Println("[DB]: Can't Select user by Id from DB \t-->\t" + err.Error())
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Firstname, &user.Lastname, &user.Username, &user.Email, &user.Password, &user.Role, &user.Token)
		if err != nil {
			fmt.Print("[DB]: Can't convert DB-response to User-Object (userId: %v) \t-->\t"+err.Error(), user.Id)
		}
		users = append(users, user)
	}
	rows.Close()
	return users[0]
}

func saveUserOnDB(user User) {
	var query string = "INSERT INTO users (`Id`, `Firstname`, `Lastname`, `Username`, `Email`, `Password`, `Role`, `Token`) VALUES ('" + user.Id + "', '" + user.Firstname + "', '" + user.Lastname + "', '" + user.Username + "', '" + user.Email + "', '" + user.Password + "', '" + user.Role + "', '');"
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

func updateTokenOnDB(user User, token string) {
	var query string = "UPDATE `users` SET `Token`= '" + token + "'  WHERE Id='" + user.Id + "';"
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
	Token     string
}
