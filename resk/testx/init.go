package testx

import (
	"github.com/memo012/red-packet/resk/infra"
	"github.com/memo012/red-packet/resk/infra/base"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func init() {
	// 获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("../example/config.ini", 1)
	// 加载和解析配置文件
	conf := ini.NewIniFileConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.GinServerStarter{})

	app := infra.New(conf)
	app.Start()
}
