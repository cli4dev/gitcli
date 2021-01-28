package tmpl

const FieldsTmpl = `package field

{{range $j,$r:=. -}}
{{if not ($r.Name|varName|fieldExist) }}
//Field{{$r.Name|varName}} 字段{{.Desc}}的数据库名称
const Field{{$r.Name|varName}} = "{{$r.Name}}"{{$r.Name|varName|fieldSet}}
{{end}}
{{- end -}}
`
