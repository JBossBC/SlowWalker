package controller

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
}

var (
	taskController *TaskController
	taskOnce       sync.Once
)

func getTaskController() *TaskController {
	taskOnce.Do(func() {
		taskController = new(TaskController)
	})
	return taskController
}

func (taskController *TaskController) ExecTask(ctx *gin.Context) {
	// service.ExecTask(ctx.GetString("username"), ctx.RemoteIP(), , )
}
