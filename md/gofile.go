package md

import (
	"fmt"

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

	for _, t := range tb.Tables {
		path, content, err := t.GetGoFile(c.Args().Get(1))
		if err != nil {
			return err
		}
		fs, err := createPath(path, c.Bool("cover"))
		if err != nil {
			return err
		}
		logs.Log.Info("生成文件:", path)
		if _, err := fs.Write([]byte(content)); err != nil {
			return err
		}
		fs.Close()
	}

	path, content, err := tb.GetDBInstallFile(c.Args().Get(1))
	if err != nil {
		return err
	}
	fs, err := createPath(path, c.Bool("cover"))
	if err != nil {
		return err
	}
	fs.WriteString(content)
	fs.Close()
	return nil
}
