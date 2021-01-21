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
		key := tp
		if _, ok := enumsMap[key]; !ok {
			key = "*"
		}

		items, err := hydra.C.DB().GetRegularDB().Query(enumsMap[key], ctx.Request().GetMap())
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
{{ range $j,$t:=.Tbs -}}
{{if $t|fIsEnumTB -}}
{{$count:= 0 -}}
"{{$t.Name|rmhd|upperName}}":{###}select '{{$t.Name|rmhd|upperName}}' type {{$count = 1}}
{{- range $i,$c:=.Rows -}}
{{if $c.Con|fIsDI -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} value {{end -}}
{{if $c.Con|fIsDN -}}{{if gt $count 0}},{{end}}{{$count = 1}} t.{{$c.Name}} name {{end -}}
{{end}} from {{$t.Name}} t {###},
{{end -}}
{{- end -}}
"Province":{###}select 'Province' type , canton_code value, chinese_name name from dds_area_info where parent_code = '*' or parent_code='QG' order by canton_code{###},
"City": {###}select 'City' type , canton_code value, chinese_name name from dds_area_info where grade='2' order by canton_code {###},
"Region":{###}select 'Region' type ,canton_code value, chinese_name name from dds_area_info order by canton_code {###},
"*":  {###}select 	type, name,	value from dds_dictionary_info where type = @type and status = 0 order by sort_no,id {###},
}`
