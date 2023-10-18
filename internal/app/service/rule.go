package service

import "replite_web/internal/app/utils"

type Rule interface {
	QueryRuleAuthorization(role string) (response utils.Response)
}

func GetRuleService() Rule {
	return getRuleService()
}
