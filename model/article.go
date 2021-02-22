package model

import "ksc/entity"

const (
	SIZE = 20
)

type Article struct {
	Model
}

func (a *Article) List(page int) (data *[]entity.Article){
	maps := map[string]string{"status":"1"}
	a.Model.SearchByPage(page, SIZE, maps, &data)
	return data
}
