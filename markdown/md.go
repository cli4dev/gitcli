package markdown

import (
	"github.com/lib4dev/cli/cmds"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "md2",
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
						cli.StringFlag{
							Name:  "table,t",
							Usage: `-表名称`,
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
					Usage:  "select语句",
					Action: showSelect,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "table,t",
							Required: true,
							Usage:    `-表名称`,
						},
					},
				}, {
					Name:   "update",
					Usage:  "update语句",
					Action: showUpdate,
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
