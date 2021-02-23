package entity

type Tag struct {
	ID int `gorm:"column:id;primary_key" json:"id"`
	TagName string `gorm:"column:tag_name" json:"tag_name"`
}

//表名
func (Tag) TableName() string{
	return "tags"
}