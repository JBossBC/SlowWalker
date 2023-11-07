package controller

import (
	"log"
	"replite_web/internal/app/infrastructure"
	"replite_web/internal/app/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type MeiliSearchController struct {
}

var (
	meiliSearchController *MeiliSearchController
	meiliSearchOnce       sync.Once
)

func getMeiliSearchController() *MeiliSearchController {
	meiliSearchOnce.Do(func() {
		meiliSearchController = new(MeiliSearchController)
	})
	return meiliSearchController
}

const FunctionIndex = "SlowWalker"

func (meiliSearchController *MeiliSearchController) MeiliSearchFunctions(ctx *gin.Context) {
	label := ctx.Query("labels")
	description := ctx.Query("descriptions")
	var labels []string
	labels = append(labels, label)
	labels = utils.ParseLabel(labels)
	bytes := infrastructure.GetMeiliSearchProvider().SearchFunctions(labels, description, FunctionIndex).Serialize()
	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}

}
