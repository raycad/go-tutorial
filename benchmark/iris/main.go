package main

import (
	"os"

	"github.com/kataras/iris"
)

// https://iris-go.com/v10/recipe
// https://github.com/kataras/iris/tree/master/_examples

func newLogFile() *os.File {
	filename := "server.log"
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func main() {
	app := iris.New()

	f := newLogFile()
	defer f.Close()

	app.Logger().SetOutput(f)

	v1 := app.Party("/api/v1.0")
	{
		v1.Get("/tasks", listTask)
	}

	app.Run(iris.Addr(":9001"))
}

// Handler
func listTask(c iris.Context) {
	ret := map[string]interface{}{"id": 1, "title": "IRIS0", "done": false}

	c.StatusCode(200)
	c.JSON(ret)

	// c.Application().Logger().Info(ret)
}
