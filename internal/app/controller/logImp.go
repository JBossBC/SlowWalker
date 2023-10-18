package controller

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

const DEFAULT_PAGE_NUMBER = 10
const DEFAULT_PAGE = 1

type LogController struct {
}

var (
	logController *LogController
	logOnce       sync.Once
)

func getLogController() *LogController {
	logOnce.Do(func() {
		logController = &LogController{}
	})
	return logController
}

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

func (logController *LogController) QueryAuditLogs(ctx *gin.Context) {
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
	l := &dao.LogInfo{
		Level:    dao.LogLevel(level),
		Operator: operator,
		IP:       ip,
	}
	result = service.GetLogService().FilterLogs(l, int(page), int(pageNumber)).Serialize()
	// }
	_, err = ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}

func (logController *LogController)RemoveAuditLogs(ctx *gin.Context) {
	var result []byte
	var filters []dao.LogInfo
	err := ctx.ShouldBind(&filters)

	result = service.GetLogService().RemoveLogs(filters).Serialize()
	_, err = ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
