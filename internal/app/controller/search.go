package controller

import "github.com/gin-gonic/gin"

type Search interface {
	SearchFunctions(ctx *gin.Context)
}

func GetSearchController() Search {
	return getMeiliSearchClient()
}
