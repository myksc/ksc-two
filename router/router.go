package router

import (
	"github.com/dvwright/xss-mw"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"ksc/common"
	"ksc/middleware"
	"ksc/util"
	"net/http"
	"time"
)

const (
	WEB_PORT         = ":80"
	WEB_READTIMEOUT  = 3 * time.Second
	WEB_WRITETIMEOUT = 3 * time.Second
	WEB_MAXBYTES     = 1 << 20
)

//探针
type Probe struct {
	health *gin.HandlerFunc
	ready  *gin.HandlerFunc
}

var p Probe

func RoutersInit(currDir string) {

	//设置gin模式
	if viper.GetString("site.env") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}else{
		gin.SetMode(gin.DebugMode)
	}

	//启动带有中间件的路由：Logger、Recovery 中间件
	//gin.Default()
	r := gin.New()

	//logger
	r.Use(middleware.AccessLog())

	//recover
	r.Use(middleware.Recovery())

	//静态资源访问路径
	staticDir := util.StringBuilder(util.StringBuilder(currDir, "public"))
	r.Static("/public", staticDir)

	//性能分析工具
	common.Register(r)

	//开启CORS 避免跨越请求（不使用JSONP的原因，JSONP只支持GET请求）
	config := cors.DefaultConfig()
	site := viper.GetString("site.domain")
	config.AllowOrigins = []string{site}
	r.Use(cors.New(config))

	//XSS过滤
	var xssMdlwr xss.XssMw
	r.Use(xssMdlwr.RemoveXss())

	r.GET("/health", HealthProbe())
	r.GET("/ready", ReadyProbe())

	initMoudle(r)

	//创建HTTP服务
	server := &http.Server{
		Addr:           WEB_PORT,
		Handler:        r,
		ReadTimeout:    WEB_READTIMEOUT,
		WriteTimeout:   WEB_WRITETIMEOUT,
		MaxHeaderBytes: WEB_MAXBYTES,
	}

	//启动server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

	//监听端口默认为8080
	//r.Run(":8000")
}

//初始化路由
func initMoudle(app *gin.Engine){
	userModule(app)
	articleModule(app)
}

// HealthProbe 存活探针
func HealthProbe() gin.HandlerFunc {
	if p.health == nil {
		return func(c *gin.Context) {
			c.String(200, "succ")
		}
	}
	return *p.health
}

// ReadyProbe 就绪探针
func ReadyProbe() gin.HandlerFunc {
	if p.ready == nil {
		return func(c *gin.Context) {
			c.String(200, "succ")
		}
	}
	return *p.ready
}