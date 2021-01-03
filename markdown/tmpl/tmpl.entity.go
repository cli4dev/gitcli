package tmpl

const entityTmpl = `

//{{.name|rmheader|cname}} {{.desc}} 
type {{.name|rmheader|cname}} struct {		
	{{range $i,$c:=.columns -}}
	//{{$c.name|cname}} {{$c.desc}}
	{{$c.name|cname}} {{$c.type|cstype}} {###}json:"{{$c.name|lower}}" {{if not .isnull}} valid:"required"{{end}}{###}

	{{end -}}	
}
`
