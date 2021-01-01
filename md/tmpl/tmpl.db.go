package tmpl

const SQLTmpl = `
{{- if .pkg}}package {{.pkg}}
{{end -}}

{{- if .pkg}} 
//{{.name}} {{.desc}}
const {{.name}}={###}{{end -}}
	{{$luks:=len .uks|sub1 -}}
	{{$lkeys:=len .keys|sub1 -}}
	CREATE TABLE  {{.name}} (
		{{range $i,$c:=.columns -}}
		{{$c.name}} {{$c.type}} {{$c.def}} {{$c.null}} {{$c.seq}} comment '{{$c.desc}}{{$c.desc_ext}}' {{if or $c.not_end $.pk $.uks $.keys}},{{end}}
		{{end -}}
		{{- if .pk}}PRIMARY KEY ({{.pk}}){{- if or .uks .keys}},{{end -}}
		{{end -}}
		{{- if .uks}}
		{{range $i,$c:=.uks}}UNIQUE KEY {{$c.ukname}} ({{$c.ukfield}}){{if or (lt $i $luks) $.keys}},{{end}}
		{{end}}{{end -}}
		{{- if .keys}}
		{{range $i,$c:=.keys -}}KEY {{$c.kname}} ({{$c.kfield}}){{if lt $i $lkeys}},{{end}}
		{{end}}{{end -}}
  ) ENGINE=InnoDB {{.auto_increment}} DEFAULT CHARSET=utf8 COMMENT='{{.desc}}';
  {{- if .pkg}}{###}{{end -}} 
`
