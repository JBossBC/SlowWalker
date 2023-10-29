package test

import (
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"testing"
)

func TestLogin(t *testing.T) {
	user := &dao.UserInfo{
		Username: "member",
		Password: "member",
	}
	response, jwtStr := service.GetUserService().LoginAccount(user)
	expectResponse := utils.NewSuccessResponse("登录成功")
	if len(jwtStr) == 0 {
		t.Errorf("返回的jwtStr为空")
	}
	if response != expectResponse {
		t.Fatalf("登录失败")
	}
}

func TestCreateAccount(t *testing.T) {
	User := &dao.UserInfo{
		Username:    "member",
		Password:    "password123",
		PhoneNumber: "1234567890",
		Code:        "123",
	}
	response := service.GetUserService().CreateAccount(User)
	Response1 := utils.NewFailedResponse("注册次数太多,请等一会再试")
	Response2 := utils.NewFailedResponse("系统错误")
	Response3 := utils.NewFailedResponse("用户名已存在")
	Response4 := utils.NewFailedResponse("验证码出错")
	Response5 := utils.NewFailedResponse("创建失败")
	// base, ok := response.(utils.BaseResponse)
	// if !ok || base.State != true {
	// 	panic("测试出错")
	// }

	if response == Response1 {
		t.Fatalf("注册次数太多,请等一会再试")
	} else if response == Response2 {
		t.Fatalf("系统错误")
	} else if response == Response3 {
		t.Fatalf("用户名已存在")
	} else if response == Response4 {
		t.Fatalf("验证码错误")
	} else if response == Response5 {
		t.Fatalf("创建失败")
	}
}

func TestUpdateInfo(t *testing.T) {
	User := &dao.UserInfo{
		Username:    "member",
		Password:    "member",
		PhoneNumber: "15182242848",
	}
	response := service.GetUserService().UpdateInfo(User)
	Response1 := utils.NewFailedResponse("系统错误")
	Response2 := utils.NewFailedResponse("修改的用户信息不存在")
	Response3 := utils.NewFailedResponse("创建失败")
	Response4 := utils.NewFailedResponse("修改成功")
	if response == Response1 {
		t.Fatalf("系统错误")
	} else if response == Response2 {
		t.Fatalf("修改的用户信息不存在")
	} else if response == Response3 {
		t.Fatalf("创建失败")
	} else if response == Response4 {
		t.Fatalf("修改成功")
	}
}
func TestQueryUser(t *testing.T) {
	User := &dao.UserInfo{
		Username: "member",
	}
	response := service.GetUserService().QueryUsers(User)
	response1 := utils.NewFailedResponse("系统错误")
	if response == response1 {
		t.Fatalf("系统错误")
	}
	log.Println(response)
}
func TestQueryUsers(t *testing.T) {
	response := service.GetUserService().QueryUsers(1, 1)
	response1 := utils.NewFailedResponse("查询失败")
	if response == response1 {
		t.Fatalf("查询失败")
	}
	log.Println(response)
}
func TestDeleteUser(t *testing.T) {
	user := &dao.UserInfo{
		Username: "memberafafds",
	}
	response := service.GetUserService().DeleteUser(user)
	response1 := utils.NewFailedResponse("删除失败")
	if response == response1 {
		t.Fatalf("删除失败")
	}
	log.Println(response)
}
