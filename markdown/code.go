package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/urfave/cli"
)

var codeType string

func showEnitfy() func(c *cli.Context) (err error) {
	codeType = "entity"
	return showCode
}
func showField() func(c *cli.Context) (err error) {
	codeType = "field"
	return showCode
}

//showCode 生成代码语句
func showCode(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}

	//读取文件
	dbtp := tmpl.MYSQL
	tb, err := tmpl.Markdown2DB(c.Args().First())
	if err != nil {
		return err
	}

	//过滤数据表
	tb.FilteByKW(c.String("table"))
	script := entityMap[codeType]
	for _, tb := range tb.Tbs {

		//翻译文件
		content, err := tmpl.Translate(script, dbtp, tb)
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

var entityMap = map[string]string{
	"entity": tmpl.EntityTmpl,
	"field":  tmpl.FieldsTmpl,
}
