package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

type User interface {
	LoginAccount(user *dao.UserInfo) (response utils.Response, jwtStr string)
	CreateAccount(user *dao.UserInfo) (response utils.Response)
	UpdateInfo(user *dao.UserInfo) (response utils.Response)
	// QueryUserByUsername(user *dao.UserInfo) (response utils.Response)
	FilterUsers(user *dao.UserInfo) (response utils.Response)
	DeleteUser(user *dao.UserInfo) (response utils.Response)
	QueryUsers(filter dao.UserInfo) (response utils.Response)
	// QueryUsersByDepartment(userInfo *dao.UserInfo) (response utils.Response)
}

func GetUserService() User {
	return getUserService()
}
