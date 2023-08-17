package controller

import (
	"log"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

func QueryRuleAuthorization(ctx *gin.Context) {
	role, ok := ctx.Get("role")
	if !ok {
		ctx.AbortWithStatus(utils.SystemError)
		return
	}
	_, err := ctx.Writer.Write(service.GetRuleService().QueryRuleAuthorization(role.(string)).Serialize())
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
