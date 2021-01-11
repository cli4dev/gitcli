package markdown

import (
	"fmt"

	"github.com/micro-plat/gitcli/markdown/ui"
	"github.com/urfave/cli"
)

//createUI 创建web界面
func createUI(c *cli.Context) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定项目名称")
	}
	if c.Bool("clear") {
		return ui.Clear(c.Args().First())
	}

	return ui.CreateWeb(c.Args().First())

}

//createUI 创建web界面
func clear(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		return ui.Clear("")
	}
	return ui.Clear(c.Args().First())

}
