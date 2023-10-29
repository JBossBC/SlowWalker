package controller

import "github.com/gin-gonic/gin"

type Department interface {
	QueryAllDepartments(ctx *gin.Context)
}

func GetDepartmentController() Department {
	return getDepartmentController()
}
