package controller

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryAuditLogs(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if pageStr == "" || err != nil {
		ctx.AbortWithStatus(utils.BadReqest)
		return

	}
	pageNumberStr := ctx.Query("pageNumber")
	pageNumber, err := strconv.ParseUint(pageNumberStr, 10, 64)
	if pageNumberStr == "" || err != nil {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	_, err = ctx.Writer.Write(service.QueryLogs(int(page), int(pageNumber)).SerializeJSON())
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}

func QueryAuditLog(ctx *gin.Context) {
	//TODO according to the authriorty or IP to query the log
	level := ctx.Query("level")
	operator := ctx.Query("operator")
	ip := ctx.Query("ip")
	l := &dao.Log{
		Level:    dao.LogLevel(level),
		Operator: operator,
		IP:       ip,
	}
	_, err := ctx.Writer.Write(service.QueryLog(l).SerializeJSON())
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}

}
