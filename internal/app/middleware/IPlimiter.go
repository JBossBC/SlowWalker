package middleware

import (
	"replite_web/internal/app/service"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

// be known as the protection of  important interface
func IPLimiter(ctx *gin.Context) {
	if service.IsWarnIP(ctx.RemoteIP()) {
		ctx.AbortWithStatus(utils.ForbiddenAccess)
		return
	}
}
