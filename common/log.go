package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

const (
	ContextKeyLogID     = "logID"				//每次生成的业务日志id
	ContextKeyRequestID = "requestId"			//请求Id
)

// web 请求 兼容odp生成logid方式
func GetLogID(ctx *gin.Context) string {
	if ctx != nil {
		if logID := ctx.GetString(ContextKeyLogID); logID != "" {
			return logID
		}
		if ctx.Request != nil {
			if logID := ctx.GetHeader("X_BD_LOGID"); strings.TrimSpace(logID) != "" {
				ctx.Set(ContextKeyLogID, logID)
				return logID
			}
		}
	}

	usec := uint64(time.Now().UnixNano())
	logID := strconv.FormatUint(usec&0x7FFFFFFF|0x80000000, 10)

	// 这里有map并发写不安全问题，业务在job使用的时候规范ctx传参可避免，暂时不做加锁处理
	if ctx != nil {
		ctx.Set(ContextKeyLogID, logID)
	}

	return logID
}

func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return ""
	}
	requestId, exist := ctx.Get(ContextKeyRequestID)
	if exist {
		return requestId.(string)
	}
	return ""
}
