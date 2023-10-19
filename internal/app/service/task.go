package service

import "replite_web/internal/app/utils"

type Task interface {
	ExecTask(operate string, ip string, function string, params map[string]string, isLocal bool) (response utils.Response)
}

func GetTaskService() Task {
	return getTaskService()
}
