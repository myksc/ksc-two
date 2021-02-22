package article

import (
	"github.com/gin-gonic/gin"
	"ksc/controller"
	//"ksc/model/article"

	//"ksc/model/article"
	"ksc/model"
)

func List(c *gin.Context){
	a := new(model.Article)
	list := a.List(0)
	//c.JSON(200, info)
	controller.Success(c, gin.H{
		"list":list,
	}, "success")
}
