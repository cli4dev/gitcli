package tmpl

const UpdateSingle = `
{{$count:=.Rows|maxIndex -}}
{{$rcount:=.|pks|maxIndex -}}
//Update{{.Name|rmhd|varName}} 查询单条数据{{.Desc}}
const Update{{.Name|rmhd|varName}} = {###}
Update {{.Name}} t
{{- range $i,$c:=.Rows}}
t.{{$c.Name}} = @{{$c.Name}}{{if lt $i $count}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- range $i,$c:=.|pks}}
t.{{$c}} = @{{$c}}{{if lt $i $rcount}} and {{end}}
{{- end}} 
{###}`
