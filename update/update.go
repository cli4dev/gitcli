package update

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/gitcli/gitlabs"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "update",
			Usage:  "更新服务",
			Action: update,
		})
}

//pull 根据传入的路径(分组/仓库)拉取所有项目
func update(c *cli.Context) (err error) {
	return gitlabs.Update()
}
