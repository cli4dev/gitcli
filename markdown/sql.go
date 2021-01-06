package markdown

import (
	"fmt"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/urfave/cli"
)

//showSQL 生成SQL语句
func showSQL(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定SQL生成类型")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定markdown文件")
	}
	tpName, ok := sqlMap[strings.ToLower(c.Args().First())]
	if !ok {
		return fmt.Errorf("不支持的SQL类型:%s", c.Args().First())
	}

	//读取文件
	dbtp := tmpl.MYSQL
	tb, err := tmpl.Markdown2DB(c.Args().Get(1))
	if err != nil {
		return err
	}

	//过滤数据表
	tb.FilteByKW(c.String("table"))

	for _, tb := range tb.Tbs {

		//根据关键字过滤
		tb.FilteRowByKW(c.String("kw"))

		//翻译文件
		content, err := tmpl.Translate(tpName, dbtp, tb)
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

var sqlMap = map[string]string{
	"insert": tmpl.InsertSingle,
	"update": tmpl.UpdateSingle,
	"select": tmpl.SelectSingle,
}
