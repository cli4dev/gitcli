package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/urfave/cli"
)

func showSelect() func(c *cli.Context) (err error) {
	return showSQL("select")
}
func showUpdate() func(c *cli.Context) (err error) {
	return showSQL("update")
}

func showInsert() func(c *cli.Context) (err error) {
	return showSQL("insert")
}

//showSQL 生成SQL语句
func showSQL(sqlType string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//读取文件
		dbtp := tmpl.MYSQL
		tpName := sqlMap[sqlType]
		tb, err := tmpl.Markdown2DB(c.Args().Get(1))
		if err != nil {
			return err
		}

		//过滤数据表
		tb.FilterByKW(c.String("table"))

		for _, tb := range tb.Tbs {

			//根据关键字过滤
			tb.FilterRowByKW(c.String("kw"))

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
}

var sqlMap = map[string]string{
	"insert": tmpl.InsertSingle,
	"update": tmpl.UpdateSingle,
	"select": tmpl.SelectSingle,
}
