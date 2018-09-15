package controllers

import (
	"github.com/astaxie/beego"
)

// TaskController defines controller data itself
type TaskController struct {
	beego.Controller
}

// Example:
//
//   req: GET /task/
//   res: 200 {"Tasks": [
//          {"ID": 1, "Title": "Learn Go", "Done": false},
//          {"ID": 2, "Title": "Buy bread", "Done": true}
//        ]}
func (tc *TaskController) ListTasks() {
	ret := map[string]interface{}{"id": 1, "title": "BEEGO", "done": false}

	tc.Ctx.Output.SetStatus(200)
	tc.Ctx.Output.JSON(ret, false, false)

	// beego.Info(ret)
}
