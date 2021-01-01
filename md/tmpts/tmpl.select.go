package tmpts

const selectSingle = `
//Get{{.name|rmheader|cname}} 查询单条数据{{.desc}}
const Get{{.name|rmheader|cname}} = {###}
select 
{{$lkeys:=len .keys|sub1 -}}
{{- range $i,$c:=.columns}}
t.{{$c.name}}{{if lt $i $lkeys}},{{end}}
{{- end}} 
from {{.name}} t
where
{{.pk}}=@{{.pk}}
{###}`
