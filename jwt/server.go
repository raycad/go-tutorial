package main

import (
	"log"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

// User defines user information
type User struct {
	Name     string
	Password string
	Role     string
	Token    string
}

// SdtClaims defines the custom claims
type SdtClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt_lib.StandardClaims
}

var (
	superSecretPassword = "raycad"
)

var userList []User

func initData() {
	userList = []User{}
	userList = append(userList,
		User{Name: "admin", Password: "admin", Role: "ADMIN, MODERATOR", Token: ""},
		User{Name: "raycad", Password: "123", Role: "MODERATOR", Token: ""},
		User{Name: "seedotech", Password: "123456", Role: "USER", Token: ""})
}

func getUserByAccount(name string, password string) (*User, int) {
	for i, user := range userList {
		if user.Name == name && user.Password == password {
			return &user, i
		}
	}

	return nil, -1
}

func getUserByToken(token string) *User {
	for _, user := range userList {
		if user.Token == token {
			return &user
		}
	}

	return nil
}

func main() {
	r := gin.Default()

	// Initialize data
	initData()

	login := r.Group("/login")
	login.POST("/", func(c *gin.Context) {
		name := c.PostForm("name")
		password := c.PostForm("password")

		user, i := getUserByAccount(name, password)
		if user == nil {
			c.JSON(401, gin.H{"message": "Wrong user information!"})
			return
		}

		log.Println("name = " + name)

		claims := SdtClaims{
			user.Name,
			user.Role,
			jwt_lib.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
				Issuer:    "seedotech",
			},
		}

		token := jwt_lib.NewWithClaims(jwt_lib.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(superSecretPassword))

		if err != nil {
			c.JSON(500, gin.H{"message": "Could not generate token"})
			return
		}

		// Update new token string
		userList[i].Token = tokenString

		c.JSON(200, gin.H{"token": tokenString})
	})

	apiV1 := r.Group("api/v1")
	apiV1.Use(jwt.Auth(superSecretPassword))

	apiV1.GET("/listContact", func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")

		user := getUserByToken(tokenString)
		if user == nil {
			c.JSON(401, gin.H{"message": "Wrong token"})
			return
		}

		if user.Name == "admin" {
			c.JSON(200, gin.H{"message": "Hello " + user.Name + "! You are admin user"})
		} else {
			c.JSON(200, gin.H{"message": "Hello " + user.Name + "! You are normal user"})
		}
		/*claims := SdtClaims{}
		token, err := jwt_lib.ParseWithClaims(tokenString, &claims, func(token *jwt_lib.Token) (interface{}, error) {
			return []byte(superSecretPassword), nil
		})

		log.Println("Name "+claims.Name, token.Valid, err)
		if token.Valid == false || err != nil {
			c.JSON(401, gin.H{"message": "Wrong token"})
			return
		}

		if claims.Name == "raycad" {
			c.JSON(200, gin.H{"message": "Hello " + claims.Name + "! You have permisson!"})
		} else {
			c.JSON(401, gin.H{"message": "Sorry " + claims.Name + "! You don't have permisson!"})
		}*/
	})

	r.Run(":9000")
}
