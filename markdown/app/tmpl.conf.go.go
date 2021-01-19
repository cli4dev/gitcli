package app

const tmplConfGo = `package main

import (
	"github.com/micro-plat/hydra"
)

func install() {
	hydra.OnReady(func() error {		
		return nil
	})
}
`
