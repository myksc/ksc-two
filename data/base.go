package data

import (
	"github.com/jinzhu/gorm"
)

type BaseData struct {
	Db *gorm.DB
}

// SearchByPage 分页检索视频
func (d *BaseData) SearchByPage(page int , size int, maps interface{}) (data interface{}){
	offset := page * size
	d.Db.Where("status = ?", 0).Order("update_time DESC").Offset(offset).Limit(size).Find(&data)
	return data
}


