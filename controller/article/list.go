package article

import (
	"github.com/gin-gonic/gin"
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
		"list":list,
	}, "success")
}
