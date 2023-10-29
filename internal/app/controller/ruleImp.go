package controller

import (
	"log"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type RuleController struct {
}

var (
	ruleController *RuleController
	ruleOnce       sync.Once
)

func getRuleController() *RuleController {
	ruleOnce.Do(func() {
		ruleController = new(RuleController)
	})
	return ruleController
}

func (ruleController *RuleController) QueryRuleAuthorization(ctx *gin.Context) {
	role, ok := ctx.Get("role")
	if !ok {
		ctx.AbortWithStatus(utils.SystemError)
		return
	}
	_, err := ctx.Writer.Write(service.GetRuleService().QueryRuleAuthorization(role.(string)).Serialize())
	if err != nil {
		log.Printf("[ruleContrller][QueryRuleAuthorization]写入response信息失败:%s", err.Error())
	}
}
