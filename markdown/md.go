package markdown

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
						cli.StringFlag{
							Name:  "table,t",
							Usage: `-表名称`,
						},
						cli.BoolFlag{
							Name:  "drop,d",
							Usage: `-包含表删除语句`,
						},
						cli.BoolFlag{
							Name:  "seqfile,s",
							Usage: `-包含序列文件`,
						},
						cli.BoolFlag{
							Name:  "cover,v",
							Usage: `-文件已存在时自动覆盖`,
						},
					},
				}, {
					Name:   "sql",
					Usage:  "sql语句,如：select,update,insert",
					Action: showSQL,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "kw,k",
							Usage: `-约束字段`,
						},
						cli.StringFlag{
							Name:     "table,t",
							Required: true,
							Usage:    `-表名称`,
						},
					},
				},
				{
					Name:   "code",
					Usage:  "显示实体信息 ",
					Action: showCode,
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
