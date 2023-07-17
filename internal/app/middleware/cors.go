package middleware

import (
	"net/http"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

func CORS(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE,OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(utils.ProtocolConvert)
		return
	}

	ctx.Next()

}
