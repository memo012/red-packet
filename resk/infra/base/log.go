package base

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 定义日志格式
	formatter := &log.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-01.15:04:05"
	log.SetFormatter(formatter)

	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}
	// 控制台高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false
	// 日志文件和滚动配置
	log.Info("测试")
	log.Debug("测试")
}
