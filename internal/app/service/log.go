package service

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"strconv"
)

// func QueryLogs(page int, pageNumber int) (response utils.Response) {
// 	result, err := dao.QueryLogs(page, pageNumber)
// 	if err != nil {
// 		log.Printf("查询日志失败:%s", err.Error())
// 		response = utils.NewFailedResponse("查询失败")
// 		return
// 	}
// 	return utils.NewSuccessResponse(result)
// }

func FilterLogs(l *dao.Log, page int, pageNumber int) (response utils.Response) {
	result, logcount, err := dao.FilterLogs(l, page, pageNumber)
	//传入初始化好的日志结构体，其中包含日志层级，操作人员，ip地址
	//传入要查询第几页以及一页要显示多少条日志

	//max return page
	if err != nil {
		log.Printf("查询日志:%v,page:%d,pageNumber:%d失败:%s", l, page, pageNumber, err.Error())
		response = utils.NewFailedResponse("系统出错")
		return
	}

	// if result.IsEmpty() {
	// 	response = utils.NewFailedResponse("日志不存在")
	// 	return
	// }

	return utils.NewSuccessLogResponse(result, strconv.Itoa(logcount)) //把查询结果和总的日志条数返回到前端
}
