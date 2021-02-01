package app

import (
	"fmt"
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
	"github.com/micro-plat/gitcli/markdown/utils"
)

var tmptls = map[string]string{
	"main.go": tmplMainGo,
	"conf.go": tmplConfGo,
	//"app.go":  tmplAppGo,
	"go.mod": tmplGoMod,
}

//CreateApp 创建web项目
func CreateApp(name string) error {
	projectPath, err := utils.GetProjectPath(name)
	if err != nil {
		return err
	}
	basePath, err := utils.GetProjectBasePath(projectPath)
	if err != nil {
		return err
	}
	for file, template := range tmptls {
		//翻译文件
		param := map[string]interface{}{
			"projectPath": projectPath,
			"router":      true,
			"basePath":    basePath,
		}
		content, err := tmpl.Translate(template, tmpl.MYSQL, param)
		if err != nil {
			return fmt.Errorf("翻译%s模板出错:%+v", file, err)
		}
		fs, err := tmpl.Create(filepath.Join(projectPath, file), true)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		fs.WriteString(content)
		fs.Close()
		logs.Log.Info("生成文件:", file)
	}
	return nil
}
