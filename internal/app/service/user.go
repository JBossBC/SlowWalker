package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

type User interface {
	LoginAccount(user *dao.User) (response utils.Response, jwtStr string)
	CreateAccount(user *dao.User) (response utils.Response)
	UpdateInfo(user *dao.User) (response utils.Response)
	QueryUser(user *dao.User) (response utils.Response)
	QueryUsers(page int, pageNumber int) (response utils.Response)
	DeleteUser(user *dao.User) (response utils.Response)
}

func GetUserService() User {
	return getUserService()
}
