package tmpl

const UpdateSingle = `
{{$count:=.Rows|rMaxIndex -}}
{{$rcount:=.|pks|sMaxIndex -}}
//Update{{.Name|rmhd|pascal}} 查询单条数据{{.Desc}}
const Update{{.Name|rmhd|pascal}} = {###}
Update 
{{- range $i,$c:=.Rows}}
t.{{$c.Name}} = @{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=.|pks}}
t.{{$c}}=@{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} 
{###}`
