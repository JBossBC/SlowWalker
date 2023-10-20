package service

import (
	"fmt"
	"replite_web/internal/app/dao"
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
	Data map[string]map[string]any `json:"data"`
}

func (ruleService *RuleService) QueryRuleAuthorization(role string) (response utils.Response) {
	allRule := dao.GetRuleDao().GetAuthority(role)
	// classifyAuthorization
	var result = make(map[string]map[string]any)
	for i := 0; i < len(allRule); i++ {
		var rule = allRule[i].(*dao.RuleInfo)
		if result[rule.Type] == nil {
			result[rule.Type] = make(map[string]any)
		}
		result[rule.Type][rule.Authority] = nil
	}
	fmt.Println(utils.NewSuccessResponse(result))
	return utils.NewSuccessResponse(result)
}
