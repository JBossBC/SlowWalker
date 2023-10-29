package main

import (
	"fmt"
	"net/http"
	"replite_web/internal/app/config"
	"replite_web/internal/app/controller"
	"replite_web/internal/app/dao"
	"replite_web/internal/app/middleware"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// init the config file
	// config.Init()
	// init go-lock config
	utils.AssemblyMutex(utils.WithStorageClient(dao.GetRedisClient()))
	engine := gin.New()
	engine.Use(middleware.ProxyMiddleware)
	engine.Use(middleware.CORS)
	// engine.Use(middleware.BeforeHandler)
	// engine.Use(middleware.Auth)
	// engine.Use(middleware.RBACMiddleware)
	// defend the dangerous ip to access system
	engine.Use(middleware.IPLimiter)
	userRoute(engine)
	mobileRoute(engine)
	auditRoute(engine)
	ruleRoute(engine)
	funcRoute(engine)
	departmentRoute(engine)
	metricsRoute(engine)
	searchRoute(engine)
	engine.Run(fmt.Sprintf(":%s", config.GetServerConfig().Port))
}
func mobileRoute(engine *gin.Engine) { //发送验证码路由
	group := engine.Group("/phone")
	group.Handle(http.MethodGet, "/send", controller.GetMobileController().SendMessage)
}

func userRoute(engine *gin.Engine) { //登录注册路由
	group := engine.Group("/user")
	group.Handle(http.MethodGet, "/login", controller.GetUserController().Login)
	group.Handle(http.MethodPost, "/register", controller.GetUserController().Register)
	queryRoute := group.Handle(http.MethodPost, "/query", controller.GetUserController().QueryUsers)
	// route.Use(middleware.BeforeHandler)
	queryRoute.Use(middleware.Auth)
	queryRoute.Use(middleware.RBACMiddleware)
	filterRoute := group.Handle(http.MethodGet, "/filter", controller.GetUserController().FilterUsers)
	filterRoute.Use(middleware.Auth)
	filterRoute.Use(middleware.RBACMiddleware)
}

func metricsRoute(engine *gin.Engine) {
	engine.Handle(http.MethodGet, "/metrics", gin.WrapH(promhttp.Handler()))

}
func auditRoute(engine *gin.Engine) { //日志查询路由
	group := engine.Group("/log")
	// group.Use(middleware.BeforeHandler)
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/query", controller.GetLogController().QueryAuditLogs)
	group.Handle(http.MethodPost, "/remove", controller.GetLogController().RemoveAuditLogs)
}

// TODO will be destroy,because the architecture is rebuilded
func ruleRoute(engine *gin.Engine) {
	group := engine.Group("/rule")
	group.Use(middleware.Auth)
	// group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/query", controller.GetRuleController().QueryRuleAuthorization)
}

func funcRoute(engine *gin.Engine) {
	group := engine.Group("/func")
	// group.Use(middleware.BeforeHandler)

	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/execute", controller.GetTaskController().ExecTask)
}

func searchRoute(engine *gin.Engine) { //搜索功能路由

	group := engine.Group("search")
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/function", controller.GetSearchController().SearchFunctions)

}
func departmentRoute(engine *gin.Engine) {
	group := engine.Group("/department")
	// group.Use(middleware.BeforeHandler)
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/querys", controller.GetDepartmentController().QueryAllDepartments)

}
