package markdown

import (
	"fmt"

	"github.com/micro-plat/gitcli/markdown/app"
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
