package infrastructure

import "replite_web/internal/app/utils"

type MeiliSearch interface {
	SearchFunctions(label []string, description string, index string) (response utils.Response)
}

func GetMeiliSearchProvider() MeiliSearch {
	return getMeiliSearchProvider()
}
