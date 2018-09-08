package main

import (
	"log"
	"net/http"
	"runtime"
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

// ProcessedMessage defines the message information after processing
type ProcessedMessage struct {
	Code    int
	Name    string
	Message gin.H
}

var (
	superSecretPassword = "raycad"
)

var userList []User

func initData() {
	log.Printf(">>>> GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	log.Printf(">>>> cpus = %d\n", cpus)
	log.Printf(">>>> GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))

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
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong user information!"})
			return
		}

		log.Printf("name = %s\n", name)

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
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
			return
		}

		// Update new token string
		userList[i].Token = tokenString

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})

	apiV1 := r.Group("api/v1.0")
	apiV1.Use(jwt.Auth(superSecretPassword))

	apiV1.GET("/listContact", func(c *gin.Context) {
		// Create copy to be used inside the goroutine
		cCp := c.Copy()
		result := make(chan ProcessedMessage)
		go func() {
			// time.Sleep(1 * time.Second)
			tokenString := cCp.Request.Header.Get("Authorization")

			user := getUserByToken(tokenString)
			if user == nil {
				result <- ProcessedMessage{http.StatusUnauthorized, user.Name, gin.H{"message": "Wrong token"}}
				return
			}

			if user.Name == "admin" {
				result <- ProcessedMessage{http.StatusOK, user.Name, gin.H{"message": "Hello " + user.Name + "! You are admin user"}}
			} else {
				result <- ProcessedMessage{http.StatusOK, user.Name, gin.H{"message": "Hello " + user.Name + "! You are normal user"}}
			}
		}()

		pm := ProcessedMessage(<-result)
		c.JSON(pm.Code, pm.Message)
		log.Println("1. Processed user = " + pm.Name)
		// log.Printf("2. Processed user = %s\n", pm.Name)

		/*// It's slower when benchmarking without using goroutines
		// log.Printf is slower than log.Println
		tokenString := c.Request.Header.Get("Authorization")

		user := getUserByToken(tokenString)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token"})
			return
		}

		if user.Name == "admin" {
			c.JSON(http.StatusOK, gin.H{"message": "Hello " + user.Name + "! You are admin user"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Hello " + user.Name + "! You are normal user"})
		}
		log.Println("2. Processed user = " + user.Name)
		// log.Printf("2. Processed user = %s", user.Name)*/

		/*claims := SdtClaims{}
		token, err := jwt_lib.ParseWithClaims(tokenString, &claims, func(token *jwt_lib.Token) (interface{}, error) {
			return []byte(superSecretPassword), nil
		})

		log.Printf("Name %s, %d, %s", claims.Name, token.Valid, err)
		if token.Valid == false || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Wrong token"})
			return
		}

		if claims.Name == "raycad" {
			c.JSON(http.StatusOK, gin.H{"message": "Hello " + claims.Name + "! You have permisson!"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Sorry " + claims.Name + "! You don't have permisson!"})
		}*/
	})

	r.Run(":9000")
}
