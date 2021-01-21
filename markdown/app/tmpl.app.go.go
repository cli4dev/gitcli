package app

const tmplAppGo = `package main
import (
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
	{{- if .router}}
	//import.router#//
	//#import.router//
	{{- end}}
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	//设置配置参数
	install()

	//启动时参数配置检查
	App.OnStarting(func(appConf app.IAPPConf) error {
		return nil
	})

	{{- if .router}}
	//service.router#//
	//#service.router//
	{{end -}}
}

`
