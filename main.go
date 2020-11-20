package main

import (
	"github.com/micro-plat/cli"
	_ "github.com/micro-plat/gitcli/clones"
	_ "github.com/micro-plat/gitcli/pulls"
	_ "github.com/micro-plat/gitcli/resets"
	_ "github.com/micro-plat/gitcli/update"
)

func main() {
	var app = cli.New(cli.WithVersion("0.1.0"))
	app.Start()
}
