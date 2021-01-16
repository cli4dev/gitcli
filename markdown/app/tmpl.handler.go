package app

//TmplServiceHandler 服务处理函数
const TmplServiceHandler = `
package {{.PKG}}

import (
	"github.com/micro-plat/hydra"
	_ "github.com/go-sql-driver/mysql"
)

//{{.Name|rmhd|varName}}Handler {{.Desc}}处理服务
type {{.Name|rmhd|varName}}Handler struct {
}

//{{.Name|rmhd|varName}}Handler {{.Desc}}查询服务
func (o *{{.Name|rmhd|varName}}Handler) QueryHandle(ctx hydra.IContext) interface{} {
	items, err := hydra.C.DB().GetRegularDB().Query(sql{{.Name|rmhd|varName}}Query, ctx.Request().GetMap())
	if err != nil {
		return err
	}
	count, err := hydra.C.DB().GetRegularDB().Scalar(sql{{.Name|rmhd|varName}}Count, ctx.Request().GetMap())
	return map[string]interface{}{
		"items": items,
		"count": count,
	}
}

const sql{{.Name|rmhd|varName}}Query = {###}
select t.* from {{.Name}} t 1 = 1 
{{- range $i,$c:=.Rows|query}}
&t.{{$c.Name}}
{{- end}}{###}
const sql{{.Name|rmhd|varName}}Count = {###}
select count(1) from {{.Name}} t 1 = 1 
{{- range $i,$c:=.Rows|query}}
&t.{{$c.Name}}
{{- end}}{###}


`
