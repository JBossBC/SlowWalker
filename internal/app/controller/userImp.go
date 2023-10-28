package controller

import (
	"fmt"
	"log"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"strconv"
	"sync"
	"time"

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
	user := &dao.UserInfo{
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
	user := &dao.UserInfo{
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
func (userController *UserController) QueryUsers(ctx *gin.Context) {

}

func (userController *UserController) FilterUser(ctx *gin.Context) {
	operate, _ := ctx.Get("username")
	opDepartment, _ := ctx.Get("department")
	template := dao.UserFilterTemplate{}
	template.Username = ctx.PostForm("username")
	template.PhoneNumber = ctx.PostForm("phone")
	template.Authority = ctx.PostForm("authority")
	template.RealName = ctx.PostForm("realName")
	template.Department = ctx.PostForm("department")
	var err error
	if startInt := ctx.PostForm("start"); startInt != "" {
		template.Start, err = strconv.ParseInt(startInt, 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为", operate.(string), opDepartment.(string))
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}
	if endInt := ctx.PostForm("end"); endInt != "" {
		template.End, err = strconv.ParseInt(ctx.PostForm("end"), 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为", operate.(string), opDepartment.(string))
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}
	if pageInt := ctx.PostForm("page"); pageInt != "" {
		template.Page, err = strconv.ParseInt(ctx.PostForm("page"), 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为", operate.(string), opDepartment.(string))
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}

	if pageNumberInt := ctx.PostForm("page"); pageNumberInt != "" {
		template.PageNumber, err = strconv.ParseInt(ctx.PostForm("pageNumber"), 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为", operate.(string), opDepartment.(string))
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}
	if template.End > time.Now().Unix() {
		dao.GetLogDao().Panicf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users的行为存在风险", operate.(string), opDepartment.(string))
	}
	service.GetUserService().FilterUsers(template)
}
