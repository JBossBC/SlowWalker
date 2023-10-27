package service

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"sync"
)

type DepartmentService struct {
	ruleHandleChain ruleChainInterface
}

var (
	departmentService     *DepartmentService
	departmentServiceOnce sync.Once
)

func getDepartmentService() *DepartmentService {
	departmentServiceOnce.Do(func() {
		departmentService = new(DepartmentService)
		//create ruleChain handle,the sequence cant change
		rChain := new(memberChain)
		rChain.setNext(new(memberChain)).setNext(new(ruleChain))
		departmentService.ruleHandleChain = rChain
	})
	return departmentService
}

type ruleChainInterface interface {
	setNext(next ruleChainInterface) ruleChainInterface
	handle(username string, ip string, role string, department string) utils.Response
}

type ruleChain struct {
	next ruleChainInterface
}

func (rule *ruleChain) setNext(next ruleChainInterface) ruleChainInterface {
	rule.next = next
	return next
}

func (rule *ruleChain) handle(username string, ip string, role string, department string) utils.Response {
	dao.GetLogDao().Panicf(username, ip, "查询非定义的role权限:(role:%s,department:%s)", role, department)
	return utils.NewFailedResponse("查询失败")
}

type adminChain struct {
	ruleChain
}

func (adminChain *adminChain) handle(username string, ip string, role string, department string) utils.Response {
	if role != "admin" {
		return adminChain.next.handle(username, ip, role, department)
	}
	result, err := dao.GetDepartmentDao().QueryDepartments()
	if err != nil {
		log.Printf("调用Dao.QueryDepartments()查询departments出错:%s", err.Error())
		return utils.NewFailedResponse("查询异常")
	}
	return utils.NewSuccessResponse(result)
}

type memberChain struct {
	ruleChain
}

func (memberChain *memberChain) handle(username string, ip string, role string, department string) utils.Response {
	if role != "member" {
		return memberChain.handle(username, ip, role, department)
	}
	result, err := dao.GetDepartmentDao().QueryDepartment(&dao.DepartmentInfo{Name: department})
	if err != nil {
		log.Printf("调用Dao.QueryDepartment()查询departments出错:%s", err.Error())
		return utils.NewFailedResponse("查询异常")
	}
	return utils.NewSuccessResponse(result)
}
func (departmentService *DepartmentService) QueryAllDepartments(username string, ip string, role string, department string) utils.Response {
	return departmentService.ruleHandleChain.handle(username, ip, role, department)
}
