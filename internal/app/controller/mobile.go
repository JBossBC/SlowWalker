package controller

import "github.com/gin-gonic/gin"

type Mobile interface {
	SendMessage(ctx *gin.Context)
}

func GetMobileController() Mobile {
	return getMobileController()
}
