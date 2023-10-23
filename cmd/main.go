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
	engine.Run(fmt.Sprintf(":%s", config.GetServerConfig().Port))
}
func mobileRoute(engine *gin.Engine) {
	group := engine.Group("/phone")
	group.Handle(http.MethodGet, "/send", controller.GetMobileController().SendMessage)
}

func userRoute(engine *gin.Engine) {
	group := engine.Group("/user")
	group.Handle(http.MethodGet, "/login", controller.GetUserController().Login)
	group.Handle(http.MethodPost, "/register", controller.GetUserController().Register)
}

func MetricsRoute(engine *gin.Engine) {
	engine.Handle("/metrics", promhttp.Handler)
}
func auditRoute(engine *gin.Engine) {
	group := engine.Group("/log")
	group.Use(middleware.BeforeHandler)
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	// group.Handle(http.MethodGet, "/query", controller.QueryAuditLogs)
	group.Handle(http.MethodGet, "/query", controller.GetLogController().QueryAuditLogs)
}

func ruleRoute(engine *gin.Engine) {
	group := engine.Group("/rule")
	group.Use(middleware.Auth)
	group.Handle(http.MethodGet, "/query", controller.GetRuleController().QueryRuleAuthorization)
}

func funcRoute(engine *gin.Engine) {
	group := engine.Group("func")
	group.Use(middleware.BeforeHandler)
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware)
	group.Handle(http.MethodGet, "/execute", controller.GetTaskController().ExecTask)
}
