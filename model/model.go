package model

import (
	"github.com/jinzhu/gorm"
	"ksc/common"
)

type Model struct {
	Db *gorm.DB
}

// SearchByPage 带分页搜索
func (model *Model) SearchByPage(page int, size int, maps interface{}, out interface{}) {
	model.Db = common.GetDb()
	offset := page * size
	model.Db.Where("status = ?", 0).Offset(offset).Limit(size).Order("update_time DESC").Find(&out)
}


// SearchByMaps 不带分页的搜索
func (model *Model) SearchByMaps(maps interface{}, out interface{})  {
	model.Db.Where(maps).Find(&out)
}
