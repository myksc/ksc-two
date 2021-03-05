package common

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"ksc/util"
)

const (
	WEB_CONFILE		 = "config/conf.yaml"
	WEB_CONFTYPE	 = "yaml"
)

// InitViper
func InitViper(appDir string){
	file  := util.StringBuilder(appDir, WEB_CONFILE)
	viper.SetConfigFile(file)
	viper.SetConfigType(WEB_CONFTYPE)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 监控配置文件变化
	viper.WatchConfig()

	//配置发生变化
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

}
