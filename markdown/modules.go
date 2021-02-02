package markdown

import (
	"fmt"
	"path"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/utils"
	"github.com/urfave/cli"
)

func createModulesSeq() func(c *cli.Context) (err error) {
	return createModules("seq")
}

func createModules(tp string) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		if len(c.Args()) == 0 {
			return fmt.Errorf("未指定markdown文件")
		}
		root := c.Args().Get(1)
		projectPath, err := utils.GetProjectPath(root)
		if err != nil {
			return err
		}

		basePath, err := utils.GetProjectBasePath(projectPath)
		if err != nil {
			return err
		}
		confPath := tmpl.GetGoConfPath(root)
		if confPath == "" {
			return
		}
		//读取文件
		template := modulesMap[tp]
		path := path.Join(projectPath, "modules/db/mysql.seq.info.go")
		if tmpl.PathExists(path) {
			return
		}

		content, err := tmpl.Translate(template, "", map[string]interface{}{
			"BasePath":    basePath,
			"ProjectPath": projectPath,
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

var modulesMap = map[string]string{
	"seq": tmpl.ModulesDBSeqTmpl,
}

// var confPathMap = map[string]string{
// 	"vue.router": "src/router/index.js",
// 	"vue.menus":  "public/menus.json",
// }
