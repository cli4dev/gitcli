package app

//TmplEnumsHandler 服务处理函数
const TmplEnumsHandler = `
package {{.PKG}}

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/types"
)

//EnumsHandler 枚举数据查询服务
type EnumsHandler struct {
}

//QueryHandle 枚举数据查询服务
func (o *EnumsHandler) QueryHandle(ctx hydra.IContext) interface{} {

	//根据传入的枚举类型获取数据
	tp := ctx.Request().GetString("type")
	if tp != "" {
		items, err := hydra.C.DB().GetRegularDB().Query(enumsMap[tp], ctx.Request().GetMap())
		if err != nil {
			return err
		}
		return items
	}

	//查询所有枚举数据
	list := types.XMaps{}
	for _, sql := range enumsMap {
		items, err := hydra.C.DB().GetRegularDB().Query(sql, ctx.Request().GetMap())
		if err != nil {
			return err
		}
		list = append(list, items...)
	}
	return list
}

var enumsMap = map[string]string{
{{- range $j,$t:=.Tbs -}}
{{if $t|isEnumTB -}}
"{{$t.Name|rmhd}}":{###}select '{{$t.Name|rmhd}}' type,{{range $i,$c:=.Rows -}}{{if $c.Con|isDI}}t.{{$c.Name}} 'value',{{end}}{{if $c.Con|isDN}}t.{{$c.Name}} 'name' {{end}}{{end}} from {{$t.Name}} t {###},
{{end -}}
{{- end -}}
}`
