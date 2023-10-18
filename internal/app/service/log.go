package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

type Log interface {
	FilterLogs(l *dao.Log, page int, pageNumber int) (response utils.Response)
}

func GetLogService() Log {
	return getLogService()
}

func RemoveLogs(filters []dao.Log) (response utils.Response) {
	err := dao.RemoveLogs(filters)
	if err != nil {
		response = utils.NewFailedResponse("删除失败")
		return
	}

	return utils.NewSuccessResponse("删除成功")

}
