package service

import (
	"fmt"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"strings"
	"sync"
)

type TaskService struct {
}

var (
	taskService *TaskService
	taskOnce    sync.Once
)

func getTaskService() *TaskService {
	taskOnce.Do(func() {
		taskService = new(TaskService)
	})
	return taskService
}

// the params resprent the user params view
func (taskService *TaskService) ExecTask(operate string, ip string, function string, params map[string]string, isLocal bool) (response utils.Response) {
	views, err := dao.GetFuncViewDao().GetFuncViews(function)
	if err != nil {
		return utils.NewFailedResponse("系统错误")
	}
	var mediumParams = make([]string, 0, len(params))
	// params pre-inspect
	for i := 0; i < len(views); i++ {
		var tmp = views[i]
		if tmp.Sign {
			if _, ok := params[tmp.View]; !ok {
				return utils.NewFailedResponse(fmt.Sprintf("%s 为必填项", tmp.View))
			}
		}
	}
	for key, value := range params {
		var valid = false
		for j := 0; j < len(views); j++ {
			var tmp = views[j]
			if strings.Compare(tmp.View, key) == 0 {
				if tmp.IsMedium {
					mediumParams = append(mediumParams, value)
				}
				valid = true
				break
			}
		}
		if !valid {
			dao.GetLogDao().Errorf(operate, ip, "操作者使用了无效的参数:%s", function)
			return utils.NewFailedResponse("参数错误")
		}
	}
	// init operate
	funcmap := dao.GetFunctionDao().GetFuncMap(function)
	if funcmap == nil {
		dao.GetLogDao().Errorf(operate, ip, "操作者正在使用未知功能:%s", function)
		return utils.NewFailedResponse("没有上传对应函数功能")
	}
	var completeParams = make([]string, 0, len(params))
	//the params has included in views  for this stage
	for key, value := range params {
		for i := 0; i < len(views); i++ {
			var tmp = views[i]
			if strings.Compare(key, tmp.View) == 0 {
				completeParams = append(completeParams, tmp.Params, value)
			}
		}
	}
	err = dao.GetLinuxPlatform().PushTask(dao.NewOperate(operate, function, dao.WithParams(completeParams), dao.WithMedium(mediumParams)))
	if err != nil {
		return utils.NewFailedResponse("任务发送失败")
	}
	return utils.NewSuccessResponse("创建任务成功")

}

type PlatformChain interface {
	setNext(PlatformChain) PlatformChain
	Handle()
}

type PlatformChainImpl struct {
	next PlatformChain
}

func (chain *PlatformChainImpl) setNext(platform PlatformChain) PlatformChain {
	chain.next = platform
	return platform
}

func (chain *PlatformChainImpl) Handle() {

}

// TODO how to design the flexiable function to add
func CreateTask() {}
