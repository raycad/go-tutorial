package main

import (
	"log"
	"time"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

var (
	superSecretPassword = "raycad"
)

type SdtClaims struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
	jwt_lib.StandardClaims
}

func main() {
	r := gin.Default()
	login := r.Group("/login")
	login.POST("/", func(c *gin.Context) {
		name := c.PostForm("name")
		scope := "ROLE_ADMIN, ROLE_MODERATOR"

		claims := SdtClaims{
			name,
			scope,
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

		c.JSON(200, gin.H{"token": tokenString})
	})

	apiV1 := r.Group("api/v1")
	apiV1.Use(jwt.Auth(superSecretPassword))

	apiV1.GET("/listContact", func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		claims := SdtClaims{}
		token, err := jwt_lib.ParseWithClaims(tokenString, &claims, func(token *jwt_lib.Token) (interface{}, error) {
			return []byte(superSecretPassword), nil
		})

		if token.Valid == true || err != nil {
			c.JSON(401, gin.H{"message": "Wrong token"})
			log.Println(token.Valid, err)
			return
		}

		if claims.Name == "raycad" {
			c.JSON(200, gin.H{"message": "Hello " + claims.Name + "! You have permisson!"})
		} else {
			c.JSON(401, gin.H{"message": "Sorry " + claims.Name + "! You don't have permisson!"})
		}
	})

	r.Run("127.0.0.1:9000")
}
