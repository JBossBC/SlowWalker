package controller

import (
	"log"
	"replite_web/internal/app/service"
	"sync"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
}

var (
	departmentController     *DepartmentController
	departmentControllerOnce sync.Once
)

func getDepartmentController() *DepartmentController {
	departmentControllerOnce.Do(func() {
		departmentController = new(DepartmentController)
	})
	return departmentController
}

func (departmentController *DepartmentController) QueryAllDepartments(ctx *gin.Context) {
	role, _ := ctx.Get("role")
	department, _ := ctx.Get("department")
	username, _ := ctx.Get("username")
	result := service.GetDepartmentService().QueryAllDepartments(username.(string), ctx.RemoteIP(), role.(string), department.(string)).Serialize()
	_, err := ctx.Writer.Write(result)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
