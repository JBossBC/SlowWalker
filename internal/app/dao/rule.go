package dao

type Rule interface {
	QueryAllRules() (rules []*RuleInfo, err error)
	QueryRules(page int, pageNumber int) (rules []*RuleInfo, err error)
	GetRule(owner string, authority string) (value any, ok bool)
	GetAuthority(owner string) []any
}

func GetRuleDao() Rule {
	return getRuleDao()
}
