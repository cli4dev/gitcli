package cli

import (
	"os"
	"path/filepath"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/urfave/cli"
)

//VERSION 版本号
var VERSION = "0.0.1"

//App  cli app
type App struct {
	app *cli.App
	log *logs.Logger
	*option
}

//Start 启动应用程序
func (a *App) Start() {
	if err := a.app.Run(os.Args); err != nil {
		a.log.Error(err)
	}
}

//New 创建app
func (a *App) New(opts ...Option) *App {
	app := &App{log: logs.New(), option: &option{version: VERSION}}

	for _, opt := range opts {
		opt(app.option)
	}

	app.app = cli.NewApp()
	app.app.Name = filepath.Base(os.Args[0])
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help,h",
		Usage: "查看帮助信息",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,v",
		Usage: "查看版本信息",
	}
	app.app.Commands = cmds.GetCmds()
	return app
}
