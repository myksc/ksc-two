package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"ksc/entity"
	"math/rand"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

const URL = "http://www.enjoybar.com"

var (
	db *gorm.DB
)

func init(){
	connstr := "root:RraDEZgfhY@tcp(127.0.0.1:3306)/ksc?charset=utf8&parseTime=true"
	driver, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err)
	}
	db = driver
}

func main(){
	fmt.Println("start")


	//爬虫
	agent := colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")
	depth := colly.MaxDepth(1)
	index := colly.NewCollector(agent, depth)
	con := infoC(index.Clone())
	index = indexC(index, con)
	err := index.Visit(URL)
	if err != nil {
		panic(err)
	}

	fmt.Println("End")
}

// indexC 首页爬取
func indexC(index *colly.Collector, info *colly.Collector) *colly.Collector {
	//请求前调用
	index.OnRequest(func(r *colly.Request) {
		fmt.Println("index爬取：", r.URL)
	})

	//请求发生错误时调用
	index.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//采集器，获取首页文章列表
	index.OnHTML("div[class='container-bd']", func(e *colly.HTMLElement) {
		csign := 0
		e.ForEach("div[class='wrap']", func(ci int, cv *colly.HTMLElement){
			cName := cv.ChildText("dt[class='title'] h3 span a")
			if cName != "" {
				csign++

				//todo 标签（插入sql）
				tagMsg := fmt.Sprintf("捕获标签成功：sign:%d, tagName:%s", csign, cName)
				fmt.Println(tagMsg)

				cv.ForEach("li", func(itemI int, itemV *colly.HTMLElement) {
					itemHref := itemV.ChildAttr("a[class='meiwen']", "href")
					itemName := itemV.ChildText("a[class='meiwen']")
					if itemName != "" && itemHref != ""{
						conUrl := fmt.Sprintf("%s%s", URL, itemHref)
						ctx := colly.NewContext()
						ctx.Put("csign", csign)
						ctx.Put("tagname", cName)
						info.Request("GET", conUrl, nil, ctx, nil)
					}
				})
			}
		})
	})

	return index
}

// infoC 内容爬取
func infoC(info *colly.Collector) *colly.Collector{

	//限速
	info.Limit(&colly.LimitRule{
		DomainRegexp: "",
		DomainGlob:   "*.enjoybar.com/*",
		Delay:        2 * time.Second,
		RandomDelay:  0,
		Parallelism:  1,
	})

	//(内容)请求前调用
	info.OnRequest(func(r *colly.Request) {
		//fmt.Println("con爬取：", r.URL)
	})

	//(内容)错误请求调用
	info.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	//(内容)
	info.OnHTML("div[class='container']", func(art *colly.HTMLElement){

		tagName  := art.Request.Ctx.Get("tagname")
		tagSign  := art.Request.Ctx.Get("csign")
		//conDate  := art.ChildText("li[class='pubdate'] span")
		//conClick := art.ChildText("li[class='click'] span")
		conArtic := art.ChildText("div[class='text'] p")
		conTitle := art.ChildText("div[class='article'] h1")

		var conImages [4]string
		var conImageInt int = 0

		art.ForEach("ul[class='picture-list'] li", func(conInt int, conImage *colly.HTMLElement){
			imageUrl := conImage.ChildAttr("a[class='meiwen'] img", "src")
			if imageUrl != "" && conImageInt <= 3 {
				conImages[conImageInt] = imageUrl
				conImageInt++
			}
		})

		imagesStr, err := json.Marshal(conImages)
		if err != nil {
			fmt.Println(err)
			return
		}

		//todo 捕捉文章成功
		conMsg := fmt.Sprintf("捕捉文章：标题：%s, 标签：%d, 内容: %s, 图片：%s, 标签名：%s", conTitle, tagSign, conArtic, imagesStr, tagName)
		fmt.Println(conMsg)
	})

	return info
}

// insertArticle 插入article表
func insertArticle(name string, sign int, content string, imgs string, tagname string){
	randNum := rand.Intn(100)
	currtime := int(time.Now().Unix())
	data := entity.Article{
		Name : name,
		Like : randNum,
		Collection : randNum,
		CreateTime : currtime,
		UpdateTime : currtime,
		Status : 1,
		Content: content,
		Imgs : imgs,
		TagId : sign,
		TagName: tagname,
	}

	//创建数据
	if err := db.Create(&data).Error; err != nil {
		panic(err)
	}else{
		msg := fmt.Sprintf("插入文章：%s", name)
		fmt.Println(msg)
	}
}

// insertTag 插入tag表
func insertTag(tagname string, tagSign int){
	data := entity.Tag{
		ID: tagSign,
		TagName: tagname,
	}

	if err := db.Create(&data).Error; err != nil {
		panic(err)
	}else{
		msg := fmt.Sprintf("插入标签：%s", tagname)
		fmt.Println(msg)
	}
}
