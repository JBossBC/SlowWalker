package controller

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	username := ctx.Query("username")
	password := ctx.Query("password")
	if username == "" || password == "" {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	user := &dao.User{
		Username: username,
		Password: password,
	}
	result, jwtStr := service.LoginAccount(user, ctx.RemoteIP())
	resultByte := result.SerializeJSON()
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

func Register(ctx *gin.Context) {
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
	}
	_, err := ctx.Writer.Write(service.CreateAccount(user, ctx.RemoteIP()).SerializeJSON())
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}
}
