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
	"time"
)

const (
	TIMELAYEROUT = "2006-01-02 15:04:05"
)

type Article struct {
	Db *gorm.DB
}

// List 列表
func (a *Article) List(page int, limit int) (schemas []schema.ArticleListRes) {
	a.Db = common.GetDb()
	offset := page * limit
	var data []entity.Article
	a.Db.Where("status = ?", 1).Order("update_time DESC").Offset(offset).Limit(limit).Find(&data)
	for _, v := range data {
		//加密sourceId
		data := util.EncryptAES([]byte(strconv.Itoa(v.ID)))
		sourceId := string(data[:])

		//创建时间
		createTime := time.Unix(int64(v.CreateTime), 0).Format(TIMELAYEROUT)

		//处理图片
		var images []string
		err := json.Unmarshal([]byte(v.Imgs), &images)
		if err != nil {
			continue
		}

		articleSchema := schema.ArticleListRes{
			SourceId: sourceId,
			Name: v.Name,
			TagSign: v.TagId,
			TagName: v.TagName,
			Images: images,
			CreateTime: createTime,
		}
		schemas = append(schemas, articleSchema)
	}
	return
}

// Insert 插入数据
func (a *Article) Insert(data *entity.Article){
	a.Db = common.GetDb()
	success := a.Db.Create(&data)
	fmt.Println(success)
}
