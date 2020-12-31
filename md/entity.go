package md

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/md/db"
	"github.com/urfave/cli"
)

//createSQL 生成SQL语句
func showEntity(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}

	//读取文件
	tb, err := db.Markdown2DB(c.Args().First())
	if err != nil {
		return err
	}

	for _, t := range tb.Tables {
		if t.Name != c.String("table") {
			continue
		}

		c, err := t.GetEntity()
		if err != nil {
			return err
		}
		logs.Log.Info(c)
	}
	return nil

}
