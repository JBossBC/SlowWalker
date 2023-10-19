package controller

import "github.com/gin-gonic/gin"

type Rule interface {
	QueryRuleAuthorization(ctx *gin.Context)
}

func GetRuleController() Rule {
	return GetRuleController()
}
