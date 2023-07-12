package middleware

import (
	"net/http"
	"replite_web/internal/app/utils"

	"github.com/gin-gonic/gin"
)

func ProxyMiddleware(c *gin.Context) {
	// 在这里进行代理验证逻辑
	// 获取代理值
	proxyValue := c.GetHeader("Proxy-Header")
     //cant allow the proxy 
	// 对代理值进行验证
	if proxyValue == "" {
		// 代理值有效，继续处理请求
		c.Next()
	} else {
		// 代理值无效，返回错误响应
		c.JSON(http.StatusForbidden, utils.NewFailedResponse("Invalid proxy"))
		c.Abort()
	}
}
