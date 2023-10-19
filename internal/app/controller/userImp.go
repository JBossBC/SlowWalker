package controller

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

var (
	userController *UserController
	userOnce       sync.Once
)

func getUserController() *UserController {
	userOnce.Do(func() {
		userController = new(UserController)
	})
	return userController
}
func (userController *UserController) Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if username == "" || password == "" {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	user := &dao.User{
		Username: username,
		Password: password,
		IP:       ctx.RemoteIP(),
	}
	result, jwtStr := service.GetUserService().LoginAccount(user)
	resultByte := result.Serialize()
	// if err != nil {
	// 	_, err := ctx.Writer.Write(utils.HELPLESS_ERROR_RESPONSE)
	// 	if err != nil {
	// 		dao.Panicf(ctx.GetString("role"), username, "登录时向response写入数据出错", err.Error())
	// 	}
	// }
	//create the login state to keep communicate with client
	ctx.Writer.Header().Add("Authorization", fmt.Sprintf("Bearer %s", jwtStr))
	_, err := ctx.Writer.Write(resultByte)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}

const DEFUALT_AUTHORITY_LEVEL = "member"

func (userController *UserController) Register(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	code := ctx.PostForm("code")
	user := &dao.User{
		Username:    username,
		Password:    password,
		PhoneNumber: phone,
		Code:        code,
		Authority:   DEFUALT_AUTHORITY_LEVEL,
		IP:          ctx.RemoteIP(),
	}
	_, err := ctx.Writer.Write(service.GetUserService().CreateAccount(user).Serialize())
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
