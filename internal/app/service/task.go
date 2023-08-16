package service

import (
	"fmt"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

// the params resprent the user params view
func ExecTask(operate string,ip string, function string, params map[string]any) (response utils.Response) {
	views, err := dao.GetFuncViews(function)
	if err != nil {
		return utils.NewFailedResponse("系统错误")
	}
	// params pre-inspect
	for i := 0; i < len(views); i++ {
		var tmp = views[i]
		if tmp.Sign {
			if _, ok := params[tmp.View]; !ok {
				return utils.NewFailedResponse(fmt.Sprintf("%s 为必填项", tmp.View))
			}
		}
	}
	for key,_:=range params{
		var tmp=params[i]
		var valid =false
		for j:=0;j<len(views);j++{
			if strings.Compare(views[j].View,key) == 0{
				valid=true
				break
			}  
		}
		if !valid{
			dao.Errorf(operate,ip,"操作者使用了无效的参数:%s",function)
			return utils.NewFailedResponse("参数错误")
		}
	}
	// init operate
	funcmap := dao.GetFuncMap(function)
	if funcmap == nil{
		dao.Errorf(operate,ip,"操作者正在使用未知功能:%s",function)
		return utils.NewFailedResponse("没有上传对应函数功能");
	}
	var completeParams:=make([]string,0,len(params))
	//the params has included in views  for this stage
	for key,_:=range params{
		views[key]
	}
	op := new(dao.BaseOperate)
	op.Function=function
	op.Params = completeParams
	

}
