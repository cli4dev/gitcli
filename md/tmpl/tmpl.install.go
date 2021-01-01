package tmpl

const mysqlInstallTmpl = `
package {{.pkg}}

import (
	"github.com/micro-plat/hydra"
)
		
func init() {
	//注册服务包
	hydra.DBCli.OnStarting(func(c hydra.ICli) error {
		hydra.Installer.DB.AddSQL(
		{{range $i,$c:=.tbs -}}
		{{$c.Name}},
		{{end -}}
		)
		return nil
	})

}
`
