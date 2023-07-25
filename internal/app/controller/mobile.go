package controller

import (
	"log"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

func SendMessage(ctx *gin.Context) {
	phone := ctx.Query("phoneNumber")
	if !utils.IsValidPhoneNumber(phone) {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	bytes := service.SendMessage(phone, ctx.RemoteIP()).SerializeJSON()
	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
