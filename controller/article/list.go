package article

import (
	"github.com/gin-gonic/gin"
	"ksc/controller"
	"ksc/model"
	"strconv"
)

var limit int = 20

func List(c *gin.Context){
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))

	if err != nil {
		controller.Fail(c, gin.H{}, "error")
		return
	}

	article := new(model.Article)
	list := article.List(page, limit)
	controller.Success(c, gin.H{
		"list":list,
	}, "success")
}
