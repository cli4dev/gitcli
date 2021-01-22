package markdown

import (
	"fmt"
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/app"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/urfave/cli"
)

func createApp(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定项目名称")
	}
	//创建项目
	err = app.CreateApp(c.Args().First())
	if err != nil {
		return err
	}
	return nil
}

func createServiceBlock() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if err := createBlockCode("service")(c); err != nil {
			return err
		}
		if err := showSQL("curd")(c); err != nil {
			return err
		}
		if !c.Bool("field") {
			return nil
		}
		if err := showCode("field")(c); err != nil {
			return err
		}
		return nil
	}
}

func createBlockCode(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		//读取文件
		dbtp := tmpl.MYSQL
		template := appCodeMap[tp]
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}

		_, projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}
		basePath, err := utils.GetProjectBasePath(projectPath)
		if err != nil {
			return err
		}

		tbs, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}

		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		for _, tb := range tbs.Tbs {

			//根据关键字过滤
			path := tmpl.GetFilePath(fmt.Sprintf("%s/services", projectPath), tb.Name, "go")
			tb.FilterRowByKW(c.String("kw"))
			tb.SetPkg(path)
			tb.SetBasePath(basePath)

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

func createEnums() func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		//读取文件
		dbtp := tmpl.MYSQL
		tbs, err := tmpl.Markdown2DB(c.Args().First())
		if err != nil {
			return fmt.Errorf("处理markdown文件表格出错:%+v", err)
		}

		//过滤数据表
		tbs.FilterByKW(c.String("table"))

		//根据关键字过滤
		path := filepath.Join(".", "system", "enums.go")
		tbs.SetPkg(path)

		//翻译文件
		content, err := tmpl.Translate(app.TmplEnumsHandler, dbtp, tbs)
		if err != nil {
			return fmt.Errorf("翻译%s模板出错:%+v", "enums", err)
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

var appCodeMap = map[string]string{
	"service": app.TmplServiceHandler,
}
