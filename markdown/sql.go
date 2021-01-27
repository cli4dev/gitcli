package markdown

import (
	"fmt"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/urfave/cli"
)

func createCurd() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		err = showSQL("curd")(c)
		if err != nil {
			return err
		}
		return createImport("driver")(c)
	}
}

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
		tb, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return err
		}
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}

		_, projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		//过滤数据表
		tb.FilterByKW(c.String("table"))

		for _, tb := range tb.Tbs {
			path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/sql", projectPath), tb.Name, "go")
			//根据关键字过滤
			tb.FilterRowByKW(c.String("kw"))
			tb.DBType = dbtp
			tb.SetPkg(path)

			//翻译文件
			content, err := tmpl.Translate(tpName, dbtp, tb)
			if err != nil {
				return err
			}
			if !c.Bool("w2f") {
				logs.Log.Info(content)
				return nil
			}
			//生成文件
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

func createImport(sqlType string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {

		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}

		//读取文件
		dbtp := tmpl.MYSQL
		tpName := sqlMap[sqlType]
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}

		_, projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		path := tmpl.GetFileName(fmt.Sprintf("%s/modules/const/sql", projectPath), dbtp, "go")

		//翻译文件
		content, err := tmpl.Translate(tpName, dbtp, nil)
		if err != nil {
			return err
		}
		if !c.Bool("w2f") {
			logs.Log.Info(content)
			return nil
		}
		//生成文件
		fs, err := tmpl.Create(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		fs.WriteString(content)
		fs.Close()
		return nil
	}
}

var sqlMap = map[string]string{
	"insert": tmpl.InsertSingle,
	"update": tmpl.UpdateSingle,
	"select": tmpl.SelectSingle,
	"curd":   tmpl.MarkdownCurdSql,
	"driver": tmpl.MarkdownCurdDriverSql,
}
