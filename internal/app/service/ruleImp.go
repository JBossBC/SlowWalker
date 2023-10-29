package service

import (
	"replite_web/internal/app/utils"
	"sync"
)

type RuleService struct {
}

var (
	ruleService *RuleService
	ruleOnce    sync.Once
)

func getRuleService() *RuleService {
	ruleOnce.Do(func() {
		ruleService = &RuleService{}
	})
	return ruleService
}

type QueryView struct {
	Role string `json:"role"`
}

func (ruleService *RuleService) QueryRuleAuthorization(role string) (response utils.Response) {
	// allRule := dao.GetRuleDao().GetAuthority(role)
	// // classifyAuthorization
	// var result = make(map[string]map[string]any)
	// for i := 0; i < len(allRule); i++ {
	// 	var rule = allRule[i].(*dao.RuleInfo)
	// 	if result[rule.Type] == nil {
	// 		result[rule.Type] = make(map[string]any)
	// 	}
	// 	result[rule.Type][rule.Authority] = nil
	// }
	queryView := new(QueryView)
	queryView.Role = role
	return utils.NewSuccessResponse(queryView)
}
