package markdown

import (
	"fmt"
	"path"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/app"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/ui"
	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/urfave/cli"
)

//createVueRouter 创建vue路由
func createVueRouter() func(c *cli.Context) (err error) {
	return createConf("vue.router")
}

//createVueMenus 创建vue菜单
func createVueMenus() func(c *cli.Context) (err error) {
	return createConf("vue.menus")
}

func createGORouter() func(c *cli.Context) (err error) {
	return createGo("conf.go")
}

func createConf(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		root := c.Args().Get(1)

		projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		webPath, webSrcPath := utils.GetWebSrcPath(projectPath)
		confPath := tmpl.GetVueConfPath(root)
		if confPath == "" {
			return
		}

		//读取文件
		template := confMap[tp]

		if webSrcPath == "" {
			return
		}
		path := path.Join(webPath, confPathMap[tp])

		confs, err := tmpl.GetSnippetConf(confPath)
		if err != nil {
			return err
		}

		content, err := tmpl.Translate(template, "", confs)
		if err != nil {
			return err
		}

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

func createGo(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		root := c.Args().Get(1)
		projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		basePath, err := utils.GetProjectBasePath(root)
		if err != nil {
			return err
		}
		confPath := tmpl.GetGoConfPath(root)
		if confPath == "" {
			return
		}

		//读取文件
		template := confMap[tp]
		path := path.Join(projectPath, "init.go")

		confs, err := tmpl.GetSnippetConf(confPath)
		if err != nil {
			return err
		}

		content, err := tmpl.Translate(template, "", map[string]interface{}{
			"BasePath":    basePath,
			"ProjectPath": projectPath,
			"Confs":       confs,
		})
		if err != nil {
			return err
		}

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

var confMap = map[string]string{
	"vue.router": ui.SnippetSrcRouterIndexJS,
	"vue.menus":  ui.SrcMenusJson,
	"conf.go":    app.SnippetTmplConfGo,
}

var confPathMap = map[string]string{
	"vue.router": "src/router/index.js",
	"vue.menus":  "public/menus.json",
}
