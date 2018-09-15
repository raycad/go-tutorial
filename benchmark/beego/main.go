package main

import (
	"./controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func initData() {
	// Remove console log
	// beego.BeeLogger.DelLogger("console")
	// Write logs to file
	beego.SetLogger(logs.AdapterFile, `{
			"filename":"logs/server.log",
			"level":7,
			"maxlines":0,
			"maxsize":0,
			"daily":true,
			"maxdays":10
		}`)
}

func main() {
	// Initialize data
	initData()

	v1 := beego.NewNamespace("/api/v1.0")
	v1.Router("/tasks", &controllers.TaskController{}, "get:ListTasks")

	beego.AddNamespace(v1)
	beego.Run(":9000")
}
