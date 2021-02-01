package app

const tmplConfGo = `package main

import (
	"fmt"

	"github.com/micro-plat/hydra"
	"github.com/micro-plat/hydra/conf/app"
	"github.com/micro-plat/hydra/conf/server/header"
	"github.com/micro-plat/hydra/conf/vars/db"
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	
	//设置配置参数
	hydra.Conf.Web("8089").Header(header.WithCrossDomain())
	hydra.Conf.Vars().DB().MySQL("db", "root", "rTo0CesHi2018Qx", "192.168.0.36:3306", "sms_test", db.WithConnect(20, 10, 600))

	//启动时参数配置检查
	App.OnStarting(func(appConf app.IAPPConf) error {

		if _, err := hydra.C.DB().GetDB(); err != nil {
			return fmt.Errorf("db数据库配置错误,err:%v", err)
		}

		return nil
	})
}

`

const SnippetTmplConfGo = `package {{if (hasSuffix .ProjectPath .BasePath )}}main{{else}}{{.ProjectPath|fileBasePath}}{{end}}

import (
	"github.com/micro-plat/hydra"
	{{- range $i,$v:=.Confs|importPath }}
	"{{$i}}"
	{{- end}}
)

//init 检查应用程序配置文件，并根据配置初始化服务
func init() {
	hydra.OnReady(func() {
	{{- range $i,$v:=.Confs }}
		hydra.S.Web("/{{$v.Name|rmhd|rpath}}", {{$v.Name|rmhd|parentPath|names|lastStr}}.New{{$v.Name|rmhd|varName}}Handler())
	{{- end}}
	})
}

`
