package tmpl

const SQLTmpl = `
{{- if .PKG}}package {{.PKG}}
{{end -}}

{{$count:=.Rows|maxIndex -}}

{{- if .PKG}} 
//{{.Name}} {{.Desc}}
const {{.Name}}={###}{{end -}}
	CREATE TABLE  {{.Name}} (
		{{range $i,$c:=.Rows -}}
		{{$c.Name}} {{$c.Type|sql}} {{$c.Def|def}} {{$c.IsNull|isnull}} {{$c|seq}} comment '{{$c.Desc}}' {{if lt $i $count}},{{end}}
		{{end -}}
	{{.|index}}{{.|pk}}) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='{{.Desc}}'
  {{- if .PKG}}{###}{{end -}} `
