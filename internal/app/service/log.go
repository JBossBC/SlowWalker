package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
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

type FilterLogView struct {
	Data  any   `json:"data"`
	Total int32 `json:"total"`
}

func FilterLogs(l *dao.Log, page int, pageNumber int) (response utils.Response) {
	result, err := dao.FilterLogs(l, page, pageNumber)
	//max return page
	if err != nil {
		// log.Printf("查询日志:%v,page:%d,pageNumber:%d失败:%s", l, page, pageNumber, err.Error())
		response = utils.NewFailedResponse("查询失败")
		return
	}
	total, err := dao.AggregateLogSum()
	if err != nil {
		response = utils.NewFailedResponse("查询失败")
		return
	}
	view := &FilterLogView{
		Total: total,
		Data:  result,
	}
	// if result.IsEmpty() {
	// 	response = utils.NewFailedResponse("日志不存在")
	// 	return
	// }
	return utils.NewSuccessResponse(view)
}


