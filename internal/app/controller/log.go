package controller

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"

	"github.com/gin-gonic/gin"
)

type Log interface {
}

func GetLogController() Log {
	return getLogController()
}

func RemoveAuditLogs(ctx *gin.Context) {
	var result []byte
	var filters []dao.LogInfo
	err := ctx.ShouldBind(&filters)

	result = service.GetLogService().RemoveLogs(filters).Serialize()
	_, err = ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
