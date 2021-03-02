package schema

type ArticleBase struct {
	SourceId string `json:"source_id"`
	Name 	 string `json:"name"`
	TagSign  int `json:"tag_sign"`
	TagName  string `json:"tag_name"`
	Images   []string `json:"images"`
	CreateTime string `json:"create_time"`
}

// 文章列表(返回内容)
type ArticleListRes struct {
	ArticleBase
}

// 内容详情
type ArticleInfoRes struct {
	ArticleBase
	Content string `json:"content"`
}
