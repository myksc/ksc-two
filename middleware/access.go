package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"ksc/common"
	"ksc/util"
	"strings"
	"time"
)

type LoggerConfig struct {
	// requestBody 打印长度
	PrintRequestLen int
	// responsebody 打印长度
	PrintResponseLen int
	// mcpack数据协议的uri，请求参数打印原始二进制
	McpackReqUris []string
	// 请求参数不打印
	IgnoreReqUris []string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	s = strings.Replace(s, "\n", "", -1)
	if w.body != nil {
		w.body.WriteString(s)
	}
	return w.ResponseWriter.WriteString(s)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		//idx := len(b)
		// gin render json 后后面会多一个换行符
		//if b[idx-1] == '\n' {
		//	b = b[:idx-1]
		//}
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

const (
	printRequestLen  = 10240
	printResponseLen = 10240
)

var (
	// 暂不需要，后续考虑看是否需要支持用户配置
	mcpackReqUris []string
	ignoreReqUris []string
)

// access日志打印
func AccessLog() gin.HandlerFunc {
	// 当前模块名
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 请求url
		path := c.Request.URL.Path
		// 请求报文
		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = c.GetRawData()
			if err != nil {
				//zlog.Warnf(c, "get http request body error: %s", err.Error())
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		blw := new(bodyLogWriter)
		if printResponseLen <= 0 {
			blw = &bodyLogWriter{body: nil, ResponseWriter: c.Writer}
		} else {
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		}
		c.Writer = blw

		c.Set("handler", c.HandlerName())
		logID := common.GetLogID(c)

		// 处理请求
		c.Next()

		response := ""
		if blw.body != nil {
			if len(blw.body.String()) <= printResponseLen {
				response = blw.body.String()
			} else {
				response = blw.body.String()[:printResponseLen]
			}
		}

		bodyStr := ""
		flag := false
		// macpack的请求，以二进制输出日志
		for _, val := range mcpackReqUris {
			if strings.Contains(path, val) {
				bodyStr = fmt.Sprintf("%v", requestBody)
				flag = true
				break
			}
		}
		if !flag {
			// 不打印RequestBody的请求
			for _, val := range ignoreReqUris {
				if strings.Contains(path, val) {
					bodyStr = ""
					flag = true
					break
				}
			}
		}
		if !flag {
			bodyStr = string(requestBody)
		}

		if c.Request.URL.RawQuery != "" {
			bodyStr += "&" + c.Request.URL.RawQuery
		}

		if len(bodyStr) > printRequestLen {
			bodyStr = bodyStr[:printRequestLen]
		}

		// 结束时间
		end := time.Now()

		// 固定notice
		commonFields := []zap.Field{
			zap.String("logId", logID),
			zap.String("requestId", common.GetRequestID(c)),
			zap.String("localIp", util.GetLocalIp()),
			zap.String("cuid", getReqValueByKey(c, "cuid")),
			zap.String("device", getReqValueByKey(c, "device")),
			zap.String("channel", getReqValueByKey(c, "channel")),
			zap.String("os", getReqValueByKey(c, "os")),
			zap.String("vc", getReqValueByKey(c, "vc")),
			zap.String("vcname", getReqValueByKey(c, "vcname")),
			zap.String("userid", getReqValueByKey(c, "userid")),
			zap.String("uri", c.Request.RequestURI),
			zap.String("host", c.Request.Host),
			zap.String("method", c.Request.Method),
			zap.String("httpProto", c.Request.Proto),
			zap.String("handle", c.HandlerName()),
			zap.String("userAgent", c.Request.UserAgent()),
			zap.String("refer", c.Request.Referer()),
			zap.String("clientIp", c.ClientIP()),
			zap.String("cookie", getCookie(c)),
			zap.String("requestStartTime", util.GetFormatRequestTime(start)),
			zap.String("requestEndTime", util.GetFormatRequestTime(end)),
			zap.Float64("cost", util.GetRequestCost(start, end)),
			zap.String("requestParam", bodyStr),
			zap.Int("responseStatus", c.Writer.Status()),
			zap.String("response", response),
		}
		common.GetAccessLogger().With(commonFields...).Info("notice")
	}
}

// 从request body中解析特定字段作为notice key打印
func getReqValueByKey(ctx *gin.Context, k string) string {
	if vs, exist := ctx.Request.Form[k]; exist && len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func getCookie(ctx *gin.Context) string {
	cStr := ""
	for _, c := range ctx.Request.Cookies() {
		cStr += fmt.Sprintf("%s=%s&", c.Name, c.Value)
	}
	return strings.TrimRight(cStr, "&")
}
