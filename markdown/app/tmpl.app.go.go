package app

const tmplAppGo = `package main
import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	//设置配置参数
	install()
}

`
