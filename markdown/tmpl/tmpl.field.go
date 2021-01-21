package tmpl

const FieldsTmpl = `
package {{.PKG}}
	//{{.Name}} {{.Desc}}的字段信息------------------------------------ 

	{{range $j,$r:=.Rows}}
	//Field{{$r.Name|varName}} 字段{{.Desc}}的数据库名称
	const Field{{$r.Name|varName}} = "{{$r.Name}}"
	{{end}}		

`
