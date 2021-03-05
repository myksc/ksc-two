package main

import (
	"ksc/common"
	"ksc/router"
	"ksc/util"
	"os"
)

func main(){
	currDir, _ := os.Getwd()
	appDir := util.StringBuilder(currDir, "/")

	//配置文件
	common.InitViper(appDir)

	//MYSQL
	common.InitDb()

	//ZAP 日志
	common.InitZap()

	//路由
	router.RoutersInit(appDir)
}