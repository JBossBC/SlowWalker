package service

import "replite_web/internal/app/utils"

type Department interface {
	QueryAllDepartments(username string, ip string, role string, department string) utils.Response
}

func GetDepartmentService() Department {
	return getDepartmentService()
}
