package tmpl

const FieldsTmpl = `
	
	//{{.Name|rmhd|varName}} {{.Desc}}------------------------------------ 

	{{range $j,$r:=.Rows}}
	//Field{{.Name|rmhd|varName}}{{$r.Name|varName}} 字段{{.Desc}}的数据库名称
	const Field{{.Name|rmhd|varName}}{{$r.Name|varName}} = "{{$r.Name}}"
	{{end}}		

`
