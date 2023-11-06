package controller

import (
	"github.com/gin-gonic/gin"
)

type MeiliSearch interface {
	MeiliSearchFunctions(ctx *gin.Context)
}

func GetMeiliSearchController() MeiliSearch {
	return getMeiliSearchController()
}
