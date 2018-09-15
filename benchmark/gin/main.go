package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()

	// gin.SetMode(gin.ReleaseMode)
	// gin.DefaultWriter = ioutil.Discard

	// Simple group: v1
	v1 := router.Group("/api/v1.0")
	{
		v1.GET("/tasks", listTask)
	}

	router.Run(":9003")
}

// Handler
func listTask(c *gin.Context) {
	ret := map[string]interface{}{"id": 1, "title": "GIN01", "done": false}

	c.JSON(http.StatusOK, ret)
}
