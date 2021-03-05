package controller


import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"ksc/common"
	"net/http"
	"time"
	"fmt"
	"strings"
)

// default response
type DefaultResponse struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

// Success 成功返回
func Success(c *gin.Context, data gin.H) {
	renderJson := DefaultResponse{0, "succ", data}
	c.JSON(http.StatusOK, renderJson)
	return
}

// Fail 失败返回
func Fail(c *gin.Context, err error) {
	//Response(c, http.StatusOK, 500, data, msg)
	var renderJson DefaultResponse
	switch errors.Cause(err).(type) {
	case common.Error:
		renderJson.ErrNo = errors.Cause(err).(common.Error).ErrNo
		renderJson.ErrMsg = errors.Cause(err).(common.Error).ErrMsg
		renderJson.Data = gin.H{}
	default:
		renderJson.ErrNo = -1
		renderJson.ErrMsg = errors.Cause(err).Error()
		renderJson.Data = gin.H{}
	}

	c.JSON(http.StatusOK, renderJson)
	// 打印错误栈
	StackLogger(c, err)
	return
}

// StackLogger 打印错误栈
func StackLogger(c *gin.Context, err error) {
	if !strings.Contains(fmt.Sprintf("%+v", err), "\n") {
		return
	}

	var info []byte
	if c != nil {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack", "requestId": common.GetLogID(c)})
	} else {
		info, _ = json.Marshal(map[string]interface{}{"time": time.Now().Format("2006-01-02 15:04:05"), "level": "error", "module": "errorstack"})
	}

	fmt.Printf("%s\n-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", string(info), err)
}



