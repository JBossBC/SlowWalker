package service

import (
	"log"
	"replite_web/internal/app/config"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
	"sync"
	"time"
)

type userService struct {
}

var (
	userSvc  *userService
	userOnce sync.Once
)

func GetUserService() *userService {
	userOnce.Do(func() {
		userSvc = new(userService)
	})
	return userSvc
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
func (userSvc *userService) LoginAccount(user *dao.User) (response utils.Response, jwtStr string) {
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
	expirationTime := time.Now().Add(time.Hour * 2) // 设置过期时间为当前时间加上2小时
	jwtStr, err = utils.CreateJWT(config.ServerConf.Secret, cliams, expirationTime)
	if err != nil {
		log.Printf("创建JWT(%v)异常:%s", cliams, err.Error())
		response = utils.NewFailedResponse("系统错误")
		return
	}
	dao.Printf(single.Authority, user.IP, "%s 登录成功,操作IP地址为:%s", user.Username, user.IP)
	return utils.NewSuccessResponse("登录成功"), jwtStr

}
func (userSvc *userService) CreateAccount(user *dao.User) (response utils.Response) {
	// whether the username  exists
	if !IsValidRegisterInternal(user.IP) {
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
	RegisterSuccessHook(user.IP)
	//delete code redis cache
	DeleteCode(user.PhoneNumber)
	// insert the remote IP to defind the numerous invalid operations
	//TODO the user.Authority is operator not  user which be created
	dao.Printf(user.IP, user.IP, "成功创建用户%s,操作IP地址:%s", user.Username, user.IP)
	return utils.NewSuccessResponse(nil)
}

/*
对于分布式条件下应该加分布式锁
*/
func (userSvc *userService) UpdateInfo(user *dao.User) (response utils.Response) {
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
	dao.Printf(single.Authority, user.IP, "成功修改%v信息为%v,操作IP地址为:%s", single, user, user.IP)
	return utils.NewSuccessResponse("修改成功")
}

func (userSvc *userService) QueryUser(user *dao.User) (response utils.Response) {
	single, err := dao.QueryUser(user)
	if err != nil {
		response = utils.NewFailedResponse("系统错误")
		return
	}
	return utils.NewSuccessResponse(single)
}

func (userSvc *userService) QueryUsers(page int, pageNumber int) (response utils.Response) {
	all, err := dao.QueryUsers(page, pageNumber)
	if err != nil {
		response = utils.NewFailedResponse("查询失败")
		return
	}
	return utils.NewSuccessResponse(all)
}

func (userSvc *userService) DeleteUser(user *dao.User) (response utils.Response) {
	err := dao.DeleteUser(user)
	if err != nil {
		response = utils.NewFailedResponse("删除失败")
		return
	}
	return utils.NewSuccessResponse(nil)
}
