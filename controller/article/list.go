package article

import (
	"github.com/gin-gonic/gin"
	"ksc/common"
	"ksc/controller"
	"ksc/model"
	"strconv"
)

func List(c *gin.Context){
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("num", "20"))

	article := new(model.Article)
	list := article.List(page, limit)
	controller.Success(c, gin.H{
		"list" : list,
	})
}

func Info(c *gin.Context){
	sourceId := c.DefaultQuery("sourceId", "")
	panic("error")
	if sourceId == "" {
		controller.Fail(c, common.NewError(500, "参数错误", "参数错误"))
	}
}
