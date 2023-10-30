package controller

import "github.com/gin-gonic/gin"

type User interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	FilterUsers(ctx *gin.Context)
	QueryUsers(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
}

func GetUserController() User {
	return getUserController()
}
