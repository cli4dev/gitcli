package app

//TmplServiceHandler 服务处理函数
const TmplServiceHandler = `
package {{.PKG}}

import (
	"net/http"
	"github.com/micro-plat/hydra"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro-plat/lib4go/errs"
)

//{{.Name|rmhd|varName}}Handler {{.Desc}}处理服务
type {{.Name|rmhd|varName}}Handler struct {
}

//QueryHandle {{.Desc}}查询服务
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
//SingleHandle 查询单条数据
func (o *{{.Name|rmhd|varName}}Handler) SingleHandle(ctx hydra.IContext) interface{} {
	items, err := hydra.C.DB().GetRegularDB().Query(sqlSignle{{.Name|rmhd|varName}}, ctx.Request().GetMap())
	if err != nil {
		return err
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}
	return items.Get(0)
}

//sql{{.Name|rmhd|varName}}Query 查询数据({{.Desc}})
const sql{{.Name|rmhd|varName}}Query = {###}
select 
{{- $count:=.Rows|maxIndex -}}
{{- range $i,$c:=.Rows}}
t.{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t where 1 = 1
{{- range $i,$c:=.Rows|query}}
&t.{{$c.Name}}
{{- end}}{###}

//sql{{.Name|rmhd|varName}}Count 查询条数({{.Desc}})
const sql{{.Name|rmhd|varName}}Count = {###}
select count(1) from {{.Name}} t where 1 = 1 
{{- range $i,$c:=.Rows|query}}
&t.{{$c.Name}}
{{- end}}{###}


const sqlSignle{{.Name|rmhd|varName}} = {###}
select 
{{- $count:=.Rows|maxIndex -}}
{{- $rcount:=.|pks|maxIndex -}}
{{- range $i,$c:=.Rows}}
t.{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t where 
{{- range $i,$c:=.|pks}}
t.{{$c}} = @{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} {###}
`
