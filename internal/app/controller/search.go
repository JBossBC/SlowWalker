package controller

import "github.com/gin-gonic/gin"

type MeiliSearch interface { //定义一个接口，名叫MeiliSearch
	MeiliSearchFunctions(ctx *gin.Context)
}

func GetMeiliSearchController() MeiliSearch {
	return getMeiliSearchController()
}
