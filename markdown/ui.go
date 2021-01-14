package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/ui"
	"github.com/urfave/cli"
)

//createUI 创建web界面
func createUI(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定项目名称")
	}
	if c.Bool("clear") {
		return ui.Clear(c.Args().First())
	}

	return ui.CreateWeb(c.Args().First())

}

//createUI 创建web界面
func clear(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return ui.Clear("")
	}
	return ui.Clear(c.Args().First())

}

//createList 创建列表页面
func createList() func(c *cli.Context) (err error) {
	return create("list")
}

//createDetail 创建详情页面
func createDetail() func(c *cli.Context) (err error) {
	return create("detail")
}
func create(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//读取文件
		dbtp := tmpl.MYSQL
		tpName := uiMap[tp]
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
}

var uiMap = map[string]string{
	"list":   ui.TmplList,
	"detail": ui.TmplDetail,
}
