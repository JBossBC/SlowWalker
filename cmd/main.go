package main

import (
	"fmt"
	"net/http"
	"replite_web/internal/app/config"
	"replite_web/internal/app/controller"
	"replite_web/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// init the config file
	// config.Init()
	engine := gin.New()
	engine.Use(middleware.ProxyMiddleware)
	engine.Use(middleware.CORS)
	// engine.Use(middleware.BeforeHandler)
	// engine.Use(middleware.Auth)
	// engine.Use(middleware.RBACMiddleware)
	userRoute(engine)
	mobileRoute(engine)
	auditRoute(engine)
	ruleRoute(engine)

	gin.DisableBindValidation()

	engine.Run(fmt.Sprintf(":%s", config.ServerConf.Port))

}

func mobileRoute(engine *gin.Engine) {
	group := engine.Group("/phone")
	group.Handle(http.MethodGet, "/send", controller.SendMessage)
}

func userRoute(engine *gin.Engine) {
	group := engine.Group("/user")
	group.Handle(http.MethodGet, "/login", controller.Login)
	group.Handle(http.MethodPost, "/register", controller.Register)
}

func auditRoute(engine *gin.Engine) {
	group := engine.Group("/log")
	group.Use(middleware.BeforeHandler)
	group.Use(middleware.Auth)
	group.Use(middleware.RBACMiddleware) //使用了三个中间件
	// group.Handle(http.MethodGet, "/query", controller.QueryAuditLogs)
	group.Handle(http.MethodGet, "/query", controller.QueryAuditLogs)
}

func ruleRoute(engine *gin.Engine) {
	group := engine.Group("/rule")
	group.Use(middleware.Auth)
	group.Handle(http.MethodGet, "/query", controller.QueryRuleAuthorization)
}
