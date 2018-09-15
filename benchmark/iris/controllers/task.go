package controllers

import "github.com/kataras/iris"

// TaskController is the equivalent
type TaskController struct {
}

func (tc *TaskController) ListTasks(c iris.Context) {
	ret := map[string]interface{}{"id": 1, "title": "test1", "done": false}

	c.StatusCode(200)
	c.JSON(ret)

	c.Application().Logger().Debug(ret)
}
