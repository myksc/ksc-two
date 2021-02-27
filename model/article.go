package model

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"ksc/common"
	"ksc/entity"
	"ksc/schema"
	"ksc/util"
	"strconv"
)

type Article struct {
	Db *gorm.DB
}

// List 列表
func (a *Article) List(page int, limit int) (schema []schema.ArticleListRes) {
	a.Db = common.GetDb()
	offset := page * limit
	var data []entity.Article
	a.Db.Where("status = ?", 1).Order("update_time DESC").Offset(offset).Limit(limit).Find(&data)
	for k, v := range data {
		//加密sourceId
		data := util.EncryptAES([]byte(strconv.Itoa(v.ID)))
		sourceId := string(data[:])

		//处理图片
		var images []string
		err := json.Unmarshal([]byte(v.Imgs), &images)
		fmt.Println(err, images, sourceId, k)


		//schema[k].SourceId = sourceId
		//schema[k].TagName  = v.TagName
		//schema[k].Name	   = v.Name
		//schema[k].TagSign  = v.TagId
	}
	return schema
}

// Insert 插入数据
func (a *Article) Insert(data *entity.Article){
	a.Db = common.GetDb()
	success := a.Db.Create(&data)
	fmt.Println(success)
}
