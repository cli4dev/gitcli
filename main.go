package main

import (
	"github.com/lib4dev/cli"
	_ "github.com/micro-plat/gitcli/clones"
	_ "github.com/micro-plat/gitcli/email"
	_ "github.com/micro-plat/gitcli/markdown"
	"github.com/micro-plat/lib4go/logger"

	_ "github.com/micro-plat/gitcli/pulls"
	_ "github.com/micro-plat/gitcli/resets"
	_ "github.com/micro-plat/gitcli/update"
)

func main() {
	logger.Pause()
	var app = cli.New(cli.WithVersion("0.1.1"))
	app.Start()
}
