package md

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lib4dev/cli/cmds"
	logs "github.com/lib4dev/cli/logger"
	"github.com/micro-plat/gitcli/md/db"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "create",
			Usage: "SQL语句",
			Subcommands: []cli.Command{
				{
					Name:   "sql",
					Usage:  "创建mysql文件",
					Action: createSQL,
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "cover,v",
							Usage: `-文件已存在时自动覆盖`,
						},
					},
				},
				{
					Name:   "go",
					Usage:  "创建go文件",
					Action: createGoFile,
				},
			},
		})
}

//createSQL 生成SQL语句
func createSQL(c *cli.Context) (err error) {
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
	files, err := db.GetSQL(tb.Tables, c.Args().Get(1), "")
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
	return nil
}

//CreatePath 创建文件，文件夹 存在时写入则覆盖
func createPath(path string, append bool) (file *os.File, err error) {

	dir := filepath.Dir(path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			err = fmt.Errorf("创建文件夹%s失败:%v", path, err)
			return nil, err
		}
	}

	_, err = os.Stat(path)
	var srcf *os.File
	if os.IsNotExist(err) {
		srcf, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
			return nil, err
		}
		return srcf, nil

	}
	if !append {
		return nil, fmt.Errorf("文件:%s已经存在", path)
	}
	srcf, err = os.OpenFile(path, os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("无法打开文件:%s(err:%v)", path, err)
		return nil, err
	}
	return srcf, nil

}
