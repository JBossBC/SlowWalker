package controller

import (
	"log"
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type MobileController struct {
}

var (
	mobileController *MobileController
	mobileOnce       sync.Once
)

func getMobileController() *MobileController {
	mobileOnce.Do(func() {
		mobileController = new(MobileController)
	})
	return mobileController
}

func (mobileController *MobileController) SendMessage(ctx *gin.Context) {
	phone := ctx.Query("phoneNumber")
	if !utils.IsValidPhoneNumber(phone) {
		ctx.AbortWithStatus(utils.BadReqest)
		return
	}
	bytes := service.GetMobileService().SendMessage(phone, ctx.RemoteIP()).Serialize()
	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		log.Printf("[mobileController][SendMessage]写入response信息失败:%s", err.Error())
	}
}
