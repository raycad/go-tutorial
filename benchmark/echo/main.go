package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	v1 := e.Group("/api/v1.0")

	// Routes
	v1.GET("/tasks", listTask)

	// Start server
	e.Logger.Fatal(e.Start(":9002"))
}

// Handler
func listTask(c echo.Context) error {
	ret := map[string]interface{}{"id": 1, "title": "ECHO1", "done": false}

	return c.JSON(http.StatusOK, ret)
}
