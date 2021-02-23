package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"ksc/common"
	"ksc/entity"
)

type Article struct {
	Db *gorm.DB
}

// List
func (a *Article) List(page int, limit int) (data []entity.Article) {
	a.Db = common.GetDb()
	offset := page * limit
	a.Db.Where("status = ?", 0).Order("update_time DESC").Offset(offset).Limit(limit).Find(&data)
	return data
}

// Insert 插入数据
func (a *Article) Insert(data *entity.Article){
	a.Db = common.GetDb()
	success := a.Db.Create(&data)
	fmt.Println(success)
}
