package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/ui"
	"github.com/urfave/cli"
)

var uiType = "list"

func createList() func(c *cli.Context) (err error) {
	uiType = "list"
	return createDetail
}
func createQuery() func(c *cli.Context) (err error) {
	uiType = "query"
	return createDetail
}
func createDetail(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}

	//读取文件
	dbtp := tmpl.MYSQL
	tpName := uiMap[uiType]
	tb, err := tmpl.Markdown2DB(c.Args().First())
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

var uiMap = map[string]string{
	"list":  ui.TmplList,
	"query": ui.TmplQueryVue,
}
