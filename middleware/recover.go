package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"ksc/common"
	"ksc/util"
	"time"
)

func Recovery() gin.HandlerFunc {
	return CatchRecover()
}

func CatchRecover() gin.HandlerFunc {
	// panic 捕获
	return func(c *gin.Context) {
		if err := recover(); err != nil {
			// 请求报文，包括了请求参数
			var requestBody []byte
			if c.Request.Body != nil {
				requestBody, _ = c.GetRawData()
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
			}
			bodyStr := string(requestBody)

			//当前时间
			ctime := time.Now()

			//转换成error
			error := err.(error)

			// 固定notice
			commonFields := []zap.Field{
				zap.String("localIp", util.GetLocalIp()),
				zap.String("uri", c.Request.RequestURI),
				zap.String("host", c.Request.Host),
				zap.String("method", c.Request.Method),
				zap.String("handle", c.HandlerName()),
				zap.String("time", util.GetFormatRequestTime(ctime)),
				zap.String("requestParam", bodyStr),
				zap.Int("responseStatus", c.Writer.Status()),
				zap.String("error", error.Error()),
			}

			common.PanicLogger(c, "panic", commonFields...)
		}

		c.Next()
	}
}
