package service

import (
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

type user struct {
}

// var (
// 	userService *user
// 	once        sync.Once
// )

//	func GetUserService() *user {
//		once.Do(func() {
//			userService = &user{}
//		})
//		return userService
//	}
func LoginAccount(user *dao.User, remoteIP string) (response utils.Response, jwtStr string) {
	single, err := dao.QueryUser(&dao.User{
		Username: user.Username,
	})
	if err != nil {
		log.Printf("查询user %v 出错\r\n", user)
		response = utils.NewFailedResponse("系统错误")
		return
	}
	if single.IsEmpty() || utils.Encrypt(user.Password) != single.Password {
		response = utils.NewFailedResponse("登录失败")
		return
	}
	cliams := utils.JwtClaims{Username: single.Username, Role: single.Authority}
	jwtStr, err = utils.CreateJWT(config.ServerConf.Secret, cliams)
	if err != nil {
		log.Printf("创建JWT(%v)异常:%s", cliams, err.Error())
		response = utils.NewFailedResponse("系统错误")
		return
	}
	dao.Printf(single.Authority, remoteIP, "%s 登录成功,操作IP地址为:%s", user.Username, remoteIP)
	return utils.NewSuccessResponse("登录成功"), jwtStr

}
func CreateAccount(user *dao.User, remoteIP string) (response utils.Response) {
	// whether the username  exists
	if !IsValidRegisterInternal(remoteIP) {
		response = utils.NewFailedResponse("注册次数太多,请等一会再试")
		return
	}
	single, err := dao.QueryUser(&dao.User{
		Username: user.Username,
	})
	if err != nil {
		log.Printf("查询user %v 出错:%s\r\n", user, err.Error())
		response = utils.NewFailedResponse("系统错误")
		return
	}
	if !single.IsEmpty() {
		response = utils.NewFailedResponse("用户名已存在")
		return
	}
	if !IsMatching(user.PhoneNumber, user.Code) {
		response = utils.NewFailedResponse("验证码出错")
		return
	}
	// encrypt the  password before keep data to database
	user.Password = utils.Encrypt(user.Password)
	// insert
	err = dao.CreateUser(user)
	if err != nil {
		log.Printf("mongoDB 创建 user  document(%v) 失败:%s \r\n", user, err.Error())
		response = utils.NewFailedResponse("创建失败")
		return
	}
	//TODO the code cant match the real result will write redis to defend error
	RegisterSuccessHook(remoteIP)
	//delete code redis cache
	DeleteCode(user.PhoneNumber)
	// insert the remote IP to defind the numerous invalid operations
	//TODO the user.Authority is operator not  user which be created
	dao.Printf(remoteIP, remoteIP, "成功创建用户%s,操作IP地址:%s", user.Username, remoteIP)
	return utils.NewSuccessResponse(nil)
}

/*
对于分布式条件下应该加分布式锁
*/
func UpdateInfo(user *dao.User, remoteIP string) (response utils.Response) {
	single, err := dao.QueryUser(&dao.User{
		Username: user.Username,
	})
	if err != nil {
		response = utils.NewFailedResponse("系统错误")
		return
	}
	if single.IsEmpty() {
		response = utils.NewFailedResponse("修改的用户信息不存在")
		return
	}
	err = dao.UpdateUser(user)
	if err != nil {
		log.Printf("mongoDB 修改 user  document(%v) 失败:%s \r\n", single, err.Error())
		response = utils.NewFailedResponse("创建失败")
		return
	}
	dao.Printf(single.Authority, remoteIP, "成功修改%v信息为%v,操作IP地址为:%s", single, user, remoteIP)
	return utils.NewSuccessResponse("修改成功")
}

func QueryUser(user *dao.User) (response utils.Response) {
	single, err := dao.QueryUser(user)
	if err != nil {
		response = utils.NewFailedResponse("系统错误")
		return
	}
	return utils.NewSuccessResponse(single)
}

func QueryUsers(page int, pageNumber int) (response utils.Response) {
	all, err := dao.QueryUsers(page, pageNumber)
	if err != nil {
		response = utils.NewFailedResponse("查询失败")
		return
	}
	return utils.NewSuccessResponse(all)
}

func DeleteUser(user *dao.User) (response utils.Response) {
	err := dao.DeleteUser(user)
	if err != nil {
		response = utils.NewFailedResponse("删除失败")
		return
	}
	return utils.NewSuccessResponse(nil)
}
