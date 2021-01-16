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
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}
		//读取文件
		dbtp := tmpl.MYSQL
		template := uiMap[tp]

		tbs, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}

		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		for _, tb := range tbs.Tbs {

			//根据关键字过滤
			tb.FilterRowByKW(c.String("kw"))

			//翻译文件
			content, err := tmpl.Translate(template, dbtp, tb)
			if err != nil {
				return fmt.Errorf("翻译%s模板出错:%+v", tp, err)
			}
			if !c.Bool("w2f") {
				logs.Log.Info(content)
				return nil
			}

			//生成文件
			path := tmpl.GetPath(root, fmt.Sprintf("%s.%s", tb.Name, tp))
			fs, err := tmpl.Create(path, c.Bool("cover"))
			if err != nil {
				return err
			}
			logs.Log.Info("生成文件:", path)
			fs.WriteString(content)
			fs.Close()

		}
		return nil
	}
}

var uiMap = map[string]string{
	"list":   ui.TmplList,
	"detail": ui.TmplDetail,
}
