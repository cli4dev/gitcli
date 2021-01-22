package app

//TmplServiceHandler 服务处理函数
const TmplServiceHandler = `
{{- $empty := "" -}}
{{- $rows := .Rows -}}
{{- $pks := .|pks -}}
package {{.PKG}}

import (
	"net/http"
	"github.com/micro-plat/hydra"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/types"
	"{{.BasePath}}/modules/const/sql"
	"{{.BasePath}}/modules/const/field"
)

//{{.Name|rmhd|varName}}Handler {{.Desc}}处理服务
type {{.Name|rmhd|varName}}Handler struct {
	{{- if gt (.Rows|create|len) 0}}
	postCheckFields   map[string]interface{} {{end}}
	{{- if gt (.Rows|detail|len) 0}}
	getCheckFields    map[string]interface{} {{end}}
	{{- if gt (.Rows|list|len) 0}}
	queryCheckFields  map[string]interface{} {{end}}
	{{- if gt (.Rows|update|len) 0}}
	updateCheckFields map[string]interface{} {{end}}
	{{- if gt (.Rows|delete|len) 0}}
	deleteCheckFields map[string]interface{} {{end}}
}

func New{{.Name|rmhd|varName}}Handler() *{{.Name|rmhd|varName}}Handler {
	return &{{.Name|rmhd|varName}}Handler{
		{{- if gt (.Rows|create|len) 0}}
		postCheckFields: map[string]interface{}{
			{{range $i,$c:=.Rows|create}}field.Field{{$c.Name|varName}}:"required",
			{{end -}}
		},
		{{- end}}
		{{- if gt (.Rows|detail|len) 0}}
		getCheckFields: map[string]interface{}{
			{{range $i,$c:=$pks}}field.Field{{$c|varName}}:"required",{{end}}
		},
		{{- end}}
		{{- if gt (.Rows|list|len) 0}}
		queryCheckFields: map[string]interface{}{
			{{range $i,$c:=.Rows|query}}field.Field{{$c.Name|varName}}:"required",
			{{end -}}
		},
		{{- end}}
		{{- if gt (.Rows|update|len) 0}}
		updateCheckFields: map[string]interface{}{
			{{range $i,$c:=.Rows|update}}field.Field{{$c.Name|varName}}:"required",
			{{end -}}
		},
		{{- end}}
		{{- if gt (.Rows|delete|len) 0}}
		deleteCheckFields: map[string]interface{}{
			{{range $i,$c:=$pks}}field.Field{{$c|varName}}:"required",{{end}}
		},
		{{- end}}
 }
}

{{- if gt (.Rows|create|len) 0}}
//PostHandle 添加{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PostHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------添加{{.Desc}}数据--------")
	
	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(u.postCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Insert{{.Name|rmhd|upperName}},ctx.Request().GetMap())
	if err != nil||count<1 {
		return errs.NewErrorf(http.StatusNotExtended,"添加数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}


{{if gt ($rows|detail|len) 0 -}}
//GetHandle 获取{{.Desc}}单条数据
func (u *{{.Name|rmhd|varName}}Handler) GetHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}单条数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(u.getCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err :=  hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}},ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}

	ctx.Log().Info("3.返回结果")
	return items.Get(0)
}
{{- end}}

{{if gt ($rows|list|len) 0 -}}
//QueryHandle  获取{{.Desc}}数据列表
func (u *{{.Name|rmhd|varName}}Handler) QueryHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}数据列表--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(u.queryCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	m := ctx.Request().GetMap()
	m["offset"] = (ctx.Request().GetInt("pi") - 1) * ctx.Request().GetInt("ps")

	count, err := hydra.C.DB().GetRegularDB().Scalar(sql.GetOrderInfoListCount, m)
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended, "查询数据数量出错:%+v", err)
	}
	
	var items types.XMaps
	if types.GetInt(count) > 0 {
		items, err = hydra.C.DB().GetRegularDB().Query(sql.GetOrderInfoList, m)
		if err != nil {
			return errs.NewErrorf(http.StatusNotExtended, "查询数据出错:%+v", err)
		}
	}

	ctx.Log().Info("3.返回结果")
	return map[string]interface{}{
		"items": items,
		"count": types.GetInt(count),
	}
}
{{- end}}

{{- if gt ($rows|update|len) 0}}
//PutHandle 更新{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PutHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------更新{{.Desc}}数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(u.updateCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Update{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}},ctx.Request().GetMap())
	if err != nil||count<1 {
		return errs.NewErrorf(http.StatusNotExtended,"更新数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}

{{- if gt ($rows|delete|len) 0}}
//DeleteHandle 删除{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) DeleteHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------删除{{.Desc}}数据--------")

	ctx.Log().Info("1.参数校验")
	if err := ctx.Request().CheckMap(u.deleteCheckFields); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Delete{{.Name|rmhd|upperName}}By{{$pks|firstStr|upperName}}, ctx.Request().GetMap())
	if err != nil||count<1 {
		return errs.NewErrorf(http.StatusNotExtended,"删除数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}
`
