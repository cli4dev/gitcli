package app

const tmplConfGo = `package main

import (
	"github.com/micro-plat/hydra"
)

func init() {
	hydra.OnReady(func() error {		
		return nil
	})
}

`
