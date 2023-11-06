package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"replite_web/internal/app/infrastructure"
	"replite_web/internal/app/utils"
	"sync"
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

func (meiliSearchController *MeiliSearchController) MeiliSearchFunctions(ctx *gin.Context) {
	label := ctx.Query("labels")
	description := ctx.Query("descriptions")
	var labels []string
	labels = append(labels, label)
	labels = utils.ParseLabel(labels)
	bytes := infrastructure.GetMeiliSearchProvider().SearchFunctions(labels, description, "WebFFF").Serialize()

	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}

}
