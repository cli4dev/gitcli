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
)

//{{.Name|rmhd|varName}}Handler {{.Desc}}处理服务
type {{.Name|rmhd|varName}}Handler struct {
}

{{- if gt (.Rows|create|len) 0}}
//PostHandle 添加{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PostHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------添加{{.Desc}}数据--------")
	ctx.Log().Info("1.参数校验")

	in := ctx.Request().GetMap()
	ck := map[string]interface{}{
		{{range $i,$c:=.Rows|create}}"{{$c.Name|lower}}":"required",
		{{end -}}	
	}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Insert{{.Name|upperName}},in)
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
	if err := ctx.Request().Check({{range $i,$c:=$pks}}"{{$c}}"{{end -}}); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err :=  hydra.C.DB().GetRegularDB().Query(sql.Get{{.Name|upperName}},ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	if items.Len() == 0 {
		return errs.NewError(http.StatusNoContent, "未查询到数据")
	}
	return items.Get(0)
}
{{- end}}

{{if gt ($rows|query|len) 0 -}}
//QueryHandle  获取{{.Desc}}数据列表
func (u *{{.Name|rmhd|varName}}Handler) QueryHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------获取{{.Desc}}数据列表--------")
	ctx.Log().Info("1.参数校验")

	in := ctx.Request().GetMap()
	ck := map[string]interface{}{
		{{range $i,$c:=.Rows|query}}"{{$c.Name|lower}}":"required",
		{{end -}}	
	}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	items, err := hydra.C.DB().GetRegularDB().Query(sql.Query{{.Name|upperName}}, ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据出错:%+v", err)
	}
	count, err := hydra.C.DB().GetRegularDB().Scalar(sql.Query{{.Name|upperName}}Count, ctx.Request().GetMap())
	if err != nil {
		return errs.NewErrorf(http.StatusNotExtended,"查询数据数量出错:%+v", err)
	}
	return map[string]interface{}{
		"items": items,
		"count": count,
	}
}
{{- end}}

{{- if gt ($rows|update|len) 0}}
//PutHandle 更新{{.Desc}}数据
func (u *{{.Name|rmhd|varName}}Handler) PutHandle(ctx hydra.IContext) (r interface{}) {

	ctx.Log().Info("--------更新{{.Desc}}数据--------")
	ctx.Log().Info("1.参数校验")

	in := ctx.Request().GetMap()
	ck := map[string]interface{}{
		{{range $i,$c:=.Rows|update}}"{{$c.Name|lower}}":"required",
		{{end -}}	
	}
	if ok, err := govalidator.ValidateMap(in, ck); !ok {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Update{{.Name|upperName}},in)
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

  if err := ctx.Request().Check({{range $i,$c:=$pks}}"{{$c}}"{{end -}}); err != nil {
		return errs.NewErrorf(http.StatusNotAcceptable, "参数校验错误:%+v", err)
	}

	ctx.Log().Info("2.执行操作")
	count,err := hydra.C.DB().GetRegularDB().Execute(sql.Delete{{.Name|upperName}}, ctx.Request().GetMap())
	if err != nil||count<1 {
		return errs.NewErrorf(http.StatusNotExtended,"删除数据出错:%+v", err)
	}

	ctx.Log().Info("3.返回结果")
	return "success"
}
{{- end}}
`
