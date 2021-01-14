package app

import (
	"path/filepath"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
)

var tmptls = map[string]string{
	"main.go": tmplMainGo,
	"conf.go": tmplConfGo,
}

//CreateApp 创建web项目
func CreateApp(name string) error {
	for path, content := range tmptls {
		fs, err := tmpl.Create(filepath.Join(".", name, path), true)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		fs.WriteString(content)
		fs.Close()
		logs.Log.Info("生成文件:", path)
	}
	return nil
}
