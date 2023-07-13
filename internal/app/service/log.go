package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

func QueryLogs(page int, pageNumber int) (response utils.Response) {
	result, err := dao.QueryLogs(page, pageNumber)
	if err != nil {
		response = utils.NewFailedResponse("查询失败")
		return
	}
	return utils.NewSuccessResponse(result)
}

func QueryLog(log *dao.Log) (response utils.Response) {
	result, err := dao.QueryLog(log)
	if err != nil {
		response = utils.NewFailedResponse("系统出错")
		return
	}
	if result.IsEmpty() {
		response = utils.NewFailedResponse("日志不存在")
		return
	}
	return utils.NewSuccessResponse(result)
}
