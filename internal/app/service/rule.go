package service

import (
	"replite_web/internal/app/dao"
	"replite_web/internal/app/utils"
)

type QueryView struct {
	Data map[string]map[string]any `json:"data"`
}

func QueryRuleAuthorization(role string) (response utils.Response) {
	allRule := dao.GetAuthority(role)
	// classifyAuthorization
	var result = make(map[string]map[string]any)
	for i := 0; i < len(allRule); i++ {
		var rule = allRule[i].(*dao.Rule)
		if result[rule.Type] == nil {
			result[rule.Type] = make(map[string]any)
		}
		result[rule.Type][rule.Authority] = nil
	}
	return utils.NewSuccessResponse(result)
}
