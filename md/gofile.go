package md

import (
	"fmt"
	"path/filepath"
	"strings"

	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/md/db"
	"github.com/urfave/cli"
)

//createSQL 生成SQL语句
func createGoFile(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定输出路径")
	}

	//读取文件
	tb, err := db.Markdown2DB(c.Args().First())
	if err != nil {
		return err
	}

	//转换为SQL语句
	names := strings.Split(strings.Trim(c.Args().Get(1), "/"), "/")
	pkgName := names[len(names)-1]
	files, err := db.GetSQL(tb.Tables, c.Args().Get(1), pkgName)
	if err != nil {
		return err
	}

	//生成文件
	for path, content := range files {
		fs, err := createPath(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		if _, err := fs.Write([]byte(content)); err != nil {
			return err
		}
	}

	installFile := filepath.Join(c.Args().Get(1), "install.go")
	fs, err := createPath(installFile, c.Bool("cover"))
	if err != nil {
		return err
	}
	fs.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	fs.WriteString(`import (
	"github.com/micro-plat/hydra"
)
		
func init() {
	//注册服务包
	hydra.DBCli.OnStarting(func(c hydra.ICli) error {
		hydra.Installer.DB.AddSQL(`)
	for _, t := range tb.Tables {
		fs.Write([]byte(fmt.Sprintf("\n			%s,", t.Name)))
	}
	fs.WriteString("\n")
	fs.WriteString(`			)
		return nil
	})

}`)

	return nil
}
