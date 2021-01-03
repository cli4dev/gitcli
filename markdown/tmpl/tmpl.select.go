package tmpl

const SelectSingle = `
{{$count:=.Rows|maxIndex -}}
{{$rcount:=.|pks|maxIndex -}}
//Select{{.Name|rmhd|pascal}} 查询单条数据{{.Desc}}
const Select{{.Name|rmhd|pascal}} = {###}
select 
{{- range $i,$c:=.Rows}}
t.{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=.|pks}}
t.{{$c}} = @{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} 
{###}`
