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
		log.Printf("[userContrller][Login]写入response信息失败:%s", err.Error())
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
		log.Printf("[userContrller][Register]写入response信息失败:%s", err.Error())
	}
}
func (userController *UserController) FilterUsers(ctx *gin.Context) {
	filterTemplate := new(dao.UserInfo)
	filterTemplate.Department = ctx.Query("department")
	filterTemplate.Username = ctx.Query("username")
	filterTemplate.RealName = ctx.Query("realName")
	filterTemplate.Authority = ctx.Query("authority")
	filterTemplate.PhoneNumber = ctx.Query("phoneNumber")
	page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	filterTemplate.Page = int(page)
	pageNumber, err := strconv.ParseInt(ctx.Query("pageNumber"), 10, 64)
	if err != nil || pageNumber <= 0 || pageNumber > MAX_PAGE_NUMBER {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	filterTemplate.PageNumber = int(pageNumber)
	_, err = ctx.Writer.Write(service.GetUserService().FilterUsers(filterTemplate).Serialize())
	if err != nil {
		log.Printf("[userContrller][QueryUsersByDepartment]写入response信息失败:%s", err.Error())
	}
}

const MAX_PAGE_NUMBER = 100

// fuzzy query
func (userController *UserController) QueryUsers(ctx *gin.Context) {
	operate, _ := ctx.Get("username")
	opDepartment, _ := ctx.Get("department")
	template := dao.UserInfo{}
	template.Username = ctx.PostForm("username")
	template.PhoneNumber = ctx.PostForm("phone")
	template.Authority = ctx.PostForm("authority")
	template.RealName = ctx.PostForm("realName")
	template.Department = ctx.PostForm("department")
	var err error
	if startInt := ctx.PostForm("start"); startInt != "" {
		template.Start, err = strconv.ParseInt(startInt, 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为:%v", operate.(string), opDepartment.(string), template)
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}
	// else {
	// 	//must include the start
	// 	ctx.AbortWithStatus(utils.BadReqest)
	// 	return
	// }
	if endInt := ctx.PostForm("end"); endInt != "" {
		template.End, err = strconv.ParseInt(endInt, 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为:%v", operate.(string), opDepartment.(string), template)
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
	}
	//  else {
	// 	//must include the end
	// 	ctx.AbortWithStatus(utils.BadReqest)
	// 	return
	// }
	if pageInt := ctx.PostForm("page"); pageInt != "" {
		page, err := strconv.ParseInt(pageInt, 10, 64)
		if err != nil {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为:%v", operate.(string), opDepartment.(string), template)
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
		template.Page = int(page)
	} else {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}

	if pageNumberInt := ctx.PostForm("pageNumber"); pageNumberInt != "" {
		pageNumber, err := strconv.ParseInt(pageNumberInt, 10, 64)
		if err != nil || template.PageNumber > MAX_PAGE_NUMBER {
			dao.GetLogDao().Errorf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users时存在错误行为:%v", operate.(string), opDepartment.(string), template)
			ctx.AbortWithStatus(utils.BadReqest)
			return
		}
		template.PageNumber = int(pageNumber)
	} else {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	if template.End > time.Now().Unix() {
		dao.GetLogDao().Panicf(operate.(string), ctx.RemoteIP(), "%s,所属部门:%s,搜索users的行为存在风险", operate.(string), opDepartment.(string))
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	_, err = ctx.Writer.Write(service.GetUserService().QueryUsers(template).Serialize())
	if err != nil {
		log.Printf("[userContrller][FilterUser]写入response信息失败:%s", err.Error())
	}
}
