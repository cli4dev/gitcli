package resets

import (
	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/gitcli/gitlabs"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:  "reset",
			Usage: "重置代码",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch,b",
					Usage: "分支",
				},
			},
			Action: reset,
		})
}

func reset(c *cli.Context) (err error) {
	reps, err := gitlabs.GetRepositories(c.Args().Get(0))
	if err != nil {
		return err
	}
	for _, rep := range reps {
		branch := types.GetString(c.String("branch"), "master")
		if err := rep.Reset(branch); err != nil {
			logs.Log.Error(err)
		}

	}
	return nil

}