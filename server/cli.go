package server

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "server",
			Usage: "运行后端应用程序",
			Subcommands: cli.Commands{
				{
					Name:   "run",
					Usage:  "启动服务",
					Action: runServer(),
				},
			},
		})
}
