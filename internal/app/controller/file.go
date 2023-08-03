package controller

import (
	"replite_web/internal/app/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

//Authorize

func AuthorizeUpload(ctx *gin.Context) {
	fileStr := ctx.Query("files")
	if fileStr == "" {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	files := strings.Split(fileStr, ",")
}
