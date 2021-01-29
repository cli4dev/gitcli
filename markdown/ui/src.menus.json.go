package ui

const SrcMenusJson = `
{{- $rows:=. -}}
[
{{- range $i,$v:=$rows}}
  {
    "name": "{{$v.Desc}}",
    "icon": "fa fa-user-circle text-primary",
    "path": "/{{$v.Name|rmhd|rpath}}"
  }{{if lt $i ($rows|maxIndex)}},{{end}}
{{- end}}
]
`
