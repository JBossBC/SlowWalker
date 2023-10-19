package controller

import (
	"github.com/gin-gonic/gin"
)

type Log interface {
	RemoveAuditLogs(ctx *gin.Context)
	QueryAuditLogs(ctx *gin.Context)
}

func GetLogController() Log {
	return getLogController()
}
