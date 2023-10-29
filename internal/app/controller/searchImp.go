package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"replite_web/internal/app/infrastructure"
)

type MeiliSearchController struct {
}

func getMeiliSearchClient() *MeiliSearchController {
	return new(MeiliSearchController)
}

func (searchController *MeiliSearchController) SearchFunctions(ctx *gin.Context) {

	label := ctx.Query("labels[]")
	description := ctx.Query("descriptions")
	var labels []string
	labels = append(labels, label)
	//之后对接service层
	bytes := infrastructure.GetMeiliSearchClient().SearchDocuments(labels, description).Serialize()
	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		log.Printf("写入response信息失败:%s", err.Error())
	}

}
