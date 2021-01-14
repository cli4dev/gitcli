package ui

import (
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/markdown/tmpl"
)

var tmptls = map[string]string{
	"src/App.vue":                srcAppVUE,
	"src/main.js":                srcMainJS,
	"src/pages/system/menus.vue": srcPagesSystemMenus,
	"src/pages/system/index.vue": srcPagesSystemIndex,
	"src/router/index.js":        srcRouterIndexJS,
	"src/store/index.js":         srcStoreIndexJS,
	"src/utility/http.js":        srcUtilityHttpJS,
	"src/utility/enums.js":       srcUtilityEnumJS,
	"index.html":                 indexHTML,
	"package.json":               packageJSON,
	"babel.config.js":            babelConfigJS,
	".gitignore":                 gitignore,
	".env.dev":                   srcEnvDev,
	".env.prod":                  srcEnvProd,
	"vue.config.js":              vueConfigJS,
}

//CreateWeb 创建web项目
func CreateWeb(name string) error {
	err := createFiles(name)
	if err != nil {
		return err
	}
	return Clear(name)

}

//Clear 清理缓存
func Clear(dir string) error {
	if err := run(dir, "npm", "install", "--no-optional", "--verbose"); err != nil {
		return err
	}
	if err := run(dir, "npm", "cache", "clear", "--force"); err != nil {
		return err
	}
	return run(dir, "npm", "install")
}

func run(dir string, name string, args ...interface{}) error {
	session := sh.InteractiveSession()
	session.SetDir(filepath.Join("./", dir))
	logs.Log.Info(append([]interface{}{name}, args...)...)
	session.Command(name, args...)
	if err := session.Run(); err != nil {
		return err
	}
	return nil
}

//createFiles 创建文件
func createFiles(name string) error {
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
