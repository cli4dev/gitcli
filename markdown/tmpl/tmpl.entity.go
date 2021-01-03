package tmpl

const EntityTmpl = `
//{{.Name|rmhd|pascal}} {{.Desc}} 
type {{.Name|rmhd|pascal}} struct {
			
	{{range $i,$c:=.Rows -}}
	//{{$c.Name|pascal}} {{$c.Desc}}
	{{$c.Name|pascal}} {{$c.Type|codeType}} {###}json:"{{$c.Name|lower}}"{###}

	{{end -}}	
}
`
