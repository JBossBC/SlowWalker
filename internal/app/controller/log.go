package controller

import (
	"fmt"
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
// 	_, err = ctx.Writer.Write(service.QueryLogs(int(page), int(pageNumber)).SerializeJSON())
// 	if err != nil {
// 		log.Printf("写入response信息失败:%s", err.Error())
// 	}
// }

func QueryAuditLogs(ctx *gin.Context) {

	fmt.Println("开始日志查询")

	//TODO according to the authriorty or IP to query the log
	level := ctx.Query("level")                     //从前端获取日志等级
	operator := ctx.Query("operator")               //从前端获取操作人员
	pageStr := ctx.Query("page")                    //从前端获取要查询第几页的内容
	page, err := strconv.ParseUint(pageStr, 10, 64) //从前端获取要显示第几页的内容
	//将字符串 pageStr 转换为无符号整数类型 uint64，10表示10进制，64表示uint64
	if pageStr == "" || err != nil {
		// ctx.AbortWithStatus(utils.BadReqest)
		// return
		page = DEFAULT_PAGE
	}
	pageNumberStr := ctx.Query("pageNumber")

	fmt.Println("page=", page)
	fmt.Println("pageNumber=", pageNumberStr)

	//page：要显示第几页
	//pageNumberStr：一页显示多少条日志数据

	//20
	//TODO add the max pageNumber limit
	pageNumber, err := strconv.ParseUint(pageNumberStr, 10, 64) //从前端获取每页要展示的日志数量
	if pageNumberStr == "" || err != nil {
		pageNumber = DEFAULT_PAGE_NUMBER
	}

	ip := ctx.Query("ip") //从前端获取ip
	var result []byte     //定义一个字节数组
	// if level == "" && operator == "" && ip == "" {
	// 	result = service.QueryLogs(int(page), int(pageNumber)).SerializeJSON()
	// } else {

	l := &dao.Log{ //初始化并且实例化一个Log结构体对象
		Level:    dao.LogLevel(level), //定义日志等级
		Operator: operator,
		IP:       ip,
	}

	fmt.Println("这是l", l)

	result = service.FilterLogs(l, int(page), int(pageNumber)).SerializeJSON() //这里要加一个，在result中还要返回当前日志库中总的日志数量
	// }
	_, err = ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
