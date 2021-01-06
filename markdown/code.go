package markdown

import (
	"fmt"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/urfave/cli"
)

//showCode 生成代码语句
func showCode(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("生成类型，如：entity,fields")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定markdown文件")
	}
	script, ok := entityMap[strings.ToLower(c.Args().First())]
	if !ok {
		return fmt.Errorf("不支持的entity类型:%s", c.Args().First())
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
