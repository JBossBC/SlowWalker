package controller

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

const DEFAULT_PAGE_NUMBER = 10
const DEFAULT_PAGE = 1

// func queryAuditLogs(ctx *gin.Context) {
// 	pageStr := ctx.Query("page")
// 	page, err := strconv.ParseUint(pageStr, 10, 64)
// 	if pageStr == "" || err != nil {
// 		// ctx.AbortWithStatus(utils.BadReqest)
// 		// return
// 		page = DEFAULT_PAGE
// 	}
// 	pageNumberStr := ctx.Query("pageNumber")
// 	pageNumber, err := strconv.ParseUint(pageNumberStr, 10, 64)
// 	if pageNumberStr == "" || err != nil {
// 		pageNumber = DEFAULT_PAGE_NUMBER
// 	}
// 	_, err = ctx.Writer.Write(service.QueryLogs(int(page), int(pageNumber)).Serialize())
// 	if err != nil {
// 		log.Printf("写入response信息失败:%s", err.Error())
// 	}
// }

func QueryAuditLogs(ctx *gin.Context) {
	//TODO according to the authriorty or IP to query the log
	level := ctx.Query("level")
	operator := ctx.Query("operator")
	pageStr := ctx.Query("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if pageStr == "" || err != nil {
		// ctx.AbortWithStatus(utils.BadReqest)
		// return
		page = DEFAULT_PAGE
	}
	pageNumberStr := ctx.Query("pageNumber")
	//20
	//TODO add the max pageNumber limit
	pageNumber, err := strconv.ParseUint(pageNumberStr, 10, 64)
	if pageNumberStr == "" || err != nil {
		pageNumber = DEFAULT_PAGE_NUMBER
	}
	ip := ctx.Query("ip")
	var result []byte
	// if level == "" && operator == "" && ip == "" {
	// 	result = service.QueryLogs(int(page), int(pageNumber)).Serialize()
	// } else {
	l := &dao.Log{
		Level:    dao.LogLevel(level),
		Operator: operator,
		IP:       ip,
	}
	result = service.FilterLogs(l, int(page), int(pageNumber)).Serialize()
	// }
	_, err = ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
