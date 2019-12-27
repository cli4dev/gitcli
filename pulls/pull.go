package pulls

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
			Name:  "pull",
			Usage: "拉取最新",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "branch,b",
					Usage: "分支",
				},
			},
			Action: pull,
		})
}

func pull(c *cli.Context) (err error) {
	branch := types.GetString(c.String("branch"), "master")
	reps, err := gitlabs.GetRepositories(c.Args().Get(0))
	if err != nil {
		return err
	}
	for _, rep := range reps {
		if !rep.Exists() {
			logs.Log.Infof("get clone %s %s", rep.FullPath, rep.GetLocalPath())
			if err := rep.Clone(); err != nil {
				logs.Log.Error(err)
			}
		}
		if err := rep.Pull(branch); err != nil {
			logs.Log.Error(err)
		}
	}
	return nil

}
