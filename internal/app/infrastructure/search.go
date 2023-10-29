package infrastructure

import "replite_web/internal/app/utils"

type Search interface {
	SearchDocuments(label []string, description string) (response utils.Response)
}

func GetMeiliSearchClient() Search {
	return getMeiliSearchClient() //use in service
}
