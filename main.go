package main

import (
	"github.com/micro-plat/cli"
	_ "github.com/micro-plat/gitcli/clones"
	_ "github.com/micro-plat/gitcli/pulls"
	_ "github.com/micro-plat/gitcli/resets"
)

func main() {
	var app = cli.New(cli.WithVersion("0.1.0"))
	app.Start()
	// session := sh.InteractiveSession()

	// session.SetDir("/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/docs")

	// session.Command("git", "branch")

	// buff, err := session.Output()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(getBranch(string(buff)))
}
