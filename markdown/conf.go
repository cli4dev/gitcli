package markdown

import (
	"fmt"
	"path"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/app"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/ui"
	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/micro-plat/lib4go/security/md5"
	"github.com/urfave/cli"
)

//createVueRouter 创建vue路由
func createVueRouter() func(c *cli.Context) (err error) {
	return createConf("vue.router")
}

func createGORouter() func(c *cli.Context) (err error) {
	return createGo("conf.go")
}

func createConf(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}

		_, projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		webPath, webSrcPath := utils.GetWebSrcPath(projectPath)
		confPath := ""
		if webPath != "" {
			confPath = path.Join(utils.GetGitcliHomePath(), fmt.Sprintf("web/web_%s.json", md5.Encrypt(webPath)))
		}

		confPath = tmpl.GetVueConfPath(root)

		//读取文件
		template := confMap[tp]

		if webSrcPath == "" {
			return
		}

		path := path.Join(webSrcPath, "router/index.js")
		logs.Log.Infof("写入文件%s/router/index.js", webSrcPath)

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
		root := ""
		if c.NArg() > 1 {
			root = c.Args().Get(1)
		}
		projectName, projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		if projectPath == "" {
			return
		}

		confPath := ""

		if projectPath != "" {
			confPath = path.Join(utils.GetGitcliHomePath(), fmt.Sprintf("server/%s_%s.json", projectName, md5.Encrypt(projectPath)))
		}

		//读取文件
		template := confMap[tp]

		path := path.Join(projectPath, "init.go")
		logs.Log.Infof("写入文件%s", path)

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

var confMap = map[string]string{
	"vue.router": ui.SnippetSrcRouterIndexJS,
	"conf.go":    app.SnippetTmplConfGo,
}
