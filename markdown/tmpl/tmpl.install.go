package tmpl

const InstallTmpl = `
package {{.PKG}}

import (
	"github.com/micro-plat/hydra"
)
		
func init() {
	//注册服务包
	hydra.DBCli.OnStarting(func(c hydra.ICli) error {
		hydra.Installer.DB.AddSQL(
		{{range $i,$c:=.Tbs -}}
		{{$c.Name}},
		{{end -}}
		)
		return nil
	}) 
}
`
