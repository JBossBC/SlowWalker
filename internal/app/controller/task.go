package controller

import "github.com/gin-gonic/gin"

type Task interface {
	ExecTask(ctx *gin.Context)
}

func GetTaskController() Task {
	return getTaskController()
}
