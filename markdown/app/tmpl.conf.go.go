package app

const tmplConfGo = `package main

func install() {
	hydra.Conf.Web("8089").Header(header.WithCrossDomain())
	hydra.Conf.Vars().DB().MySQL("db", "root", "rTo0CesHi2018Qx", "192.168.0.36:3306", "xxx", db.WithConnect(20, 10, 600))
}
`
