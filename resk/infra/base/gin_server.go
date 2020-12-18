package base

import (
	"github.com/gin-gonic/gin"
	"github.com/memo012/red-packet/resk/infra"
	"github.com/sirupsen/logrus"
)

var ginServerStarter *gin.Engine

func Gin() *gin.Engine {
	return ginServerStarter
}

type GinServerStarter struct {
	infra.BaseStarter
}

func (g *GinServerStarter) Init(ctx infra.StarterContext) {
	// 创建gin 实例
	ginServerStarter = gin.Default()
}

func (g *GinServerStarter) Start(ctx infra.StarterContext) {
	// 把路由信息打印到控制台
	routes := ginServerStarter.Routes()
	for _, g := range routes {
		logrus.Info(g)
	}
	// 启动gin
	port := ctx.Props().GetDefault("app.server.port", "8080")
	ginServerStarter.Run(":" + port)
}

func (g *GinServerStarter) StartBlocking() bool {
	return true
}
