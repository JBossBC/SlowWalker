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
