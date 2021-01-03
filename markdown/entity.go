package markdown

import (
	"fmt"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/urfave/cli"
)

//createSQL 生成SQL语句
func showEntity(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}

	//读取文件
	dbtp := tmpl.MYSQL
	tb, err := tmpl.Markdown2DB(c.Args().First())
	if err != nil {
		return err
	}

	for _, tb := range tb.Tbs {
		if !strings.Contains(tb.Name, c.String("table")) {
			continue
		}

		//翻译文件
		content, err := tmpl.Translate(tmpl.EntityTmpl, dbtp, tb)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		logs.Log.Info(content)
	}
	return nil

}
