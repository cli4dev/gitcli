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
					Name:   "sql",
					Usage:  "创建mysql文件,gitcli create sql  db.md  ../modules/const/sql/mysql ",
					Action: createSQL,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "cover,v",
							Usage: `-文件已存在时自动覆盖`,
						},
					},
				},
				{
					Name:   "gofile",
					Usage:  "创建go文件,gitcli create gofile db.md  ../modules/const/sql/mysql ",
					Action: createGoFile,
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
			},
		})
}
