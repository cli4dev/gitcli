package md

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "md",
			Usage: "SQL语句",
			Subcommands: []cli.Command{
				{
					Name:   "db",
					Usage:  "创建数据库结构文件",
					Action: createScheme,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "gofile,g",
							Usage: `-生成到gofile中`,
						},
						cli.BoolFlag{
							Name:  "cover,v",
							Usage: `-文件已存在时自动覆盖`,
						},
					},
				},
				{
					Name:   "entity",
					Usage:  "显示实体信息 ",
					Action: showEntity,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "table,t",
							Required: true,
							Usage:    `-表名称`,
						},
					},
				},
				{
					Name:   "select",
					Usage:  "获取查询语句",
					Action: showSelect,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "table,t",
							Required: true,
							Usage:    `-表名称`,
						},
					},
				},
			},
		})
}
