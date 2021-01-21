package tmpl

const MarkdownCurdSql = `
{{- $string := "string" -}}
{{- $int := "int" -}}
{{- $int64 := "int64" -}}
{{- $decimal := "types.Decimal" -}}
{{- $time := "time.Time" -}}
{{- $length := 32 -}}
{{- $createrows := .Rows|create -}}
{{- $updaterows := .Rows|update -}}
{{- $detailrows := .Rows|detail -}}
{{- $deleterows := .Rows|delete -}}
{{- $listrows := .Rows|list -}}
{{- $queryrows := .Rows|query -}}
{{- $ismysql := .DBType|ismysql -}}
{{- $isoracle := .DBType|isoracle -}}
{{- $pks := .|pks -}}
{{- $order:=.|order -}}
package sql

{{- if and $ismysql (gt ($createrows|len) 0)}}
//Insert{{.Name|upperName}} 添加{{.Desc}}
const Insert{{.Name|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{- range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values
(
	{{- range $i,$c:=$createrows}}
	@{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
){###}
{{end -}}

{{- if and $isoracle (gt ($createrows|len) 0)}}
//Insert{{.Name|upperName}} 添加{{.Desc}}
const Insert{{.Name|upperName}} = {###}
insert into {{.Name}}{{.DBLink}}
(
	{{- range $i,$c:=.seqs}}
	{{$c.Name}},{{end}}{{range $i,$c:=$createrows}}
	{{$c.Name}}{{if lt $i ($createrows|maxIndex)}},{{end}}
	{{- end}}
)
values(
	{{- range $i,$c:=.seqs}}
	{{$c.seqname}}.nextval,
	{{- end}}
	{{- range $i,$c:=$createrows}}
	{{if eq ($c.Type|codeType $time)}}to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss'){{else -}}
	@{{$c.Name}}{{end}}{{if $c.comma}},{{end}}
	{{- end}}
){###}
{{end -}}

{{- if $ismysql}}
{{if gt ($detailrows|len) 0 -}}
//Get{{.Name|upperName}} 查询单条数据{{.Desc}}
{{$tbname := .Name -}}
const Get{{.Name|upperName}} = {###}
select 
{{- range $i,$c:=$detailrows}}
t.{{$c.Name}}{{if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}}
from {{.Name}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
&{{$c}} 
{{- end}}{{end}}{###}
{{- end}}

//Query{{.Name|upperName}}Count 获取{{.Desc}}列表条数
const Query{{.Name|upperName}}Count = {###}
select count(1)
from {{.Name}} t
where 
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$queryrows -}}
{{if eq ($c.Type|codeType) $time }}
and t.{{$c.Name}}>=@{{$c.Name}} and t.{{$c.Name}}<date_add(@{{$c.Name}}, interval 1 day)
{{- else if and (gt $c.Len $length) (eq ($c.Type|codeType) $string)}}
and if(isnull(@{{$c.Name}})||@{{$c.Name}}='',1=1,t.{{$c.Name}} like CONCAT('%',@{{$c.Name}},'%'))
{{- else}}
and if(isnull(@{{$c.Name}})||@{{$c.Name}}='',1=1,t.{{$c.Name}}=@{{$c.Name}}){{end}}
{{- end}}{{end}}{###}

//Query{{.Name|upperName}} 查询{{.Desc}}列表数据
const Query{{.Name|upperName}} = {###}
select 
{{- range $i,$c:=$listrows}}
t.{{$c.Name}}{{if lt $i ($listrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}} t
where
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$queryrows -}}
{{if eq ($c.Type|codeType) $time}}
and t.{{$c.Name}}>=@{{$c.Name}} and t.{{$c.Name}}<date_add(@{{$c.Name}}, interval 1 day)
{{- else if and (gt $c.Len $length) (eq ($c.Type|codeType) $string)}}
and if(isnull(@{{$c.Name}})||@{{$c.Name}}='',1=1,t.{{$c.Name}} like CONCAT('%',@{{$c.Name}},'%'))
{{- else}}
and if(isnull(@{{$c.Name}})||@{{$c.Name}}='',1=1,t.{{$c.Name}}=@{{$c.Name}}){{end}}
{{- end}} 
{{- if gt ($order|len) 0}}
order by {{range $i,$c:=$order}}t.{{$c.Name}}{{if $c.comma}},{{else}} desc{{end}}{{end}}
{{- else}}
order by {{range $i,$c:=$pks}}t.{{$c}} desc{{end}}
{{- end}}
limit #pageSize offset #currentPage
{{end -}}{###}{{end}}

{{- if and $isoracle }}
{{if  (gt ($detailrows|len) 0) -}}
//Get{{.Name|upperName}} 查询单条数据{{.Desc}}
const Get{{.Name|upperName}} = {###}
select 
{{- range $i,$c:=$detailrows}}
t.{{$c.Name}}{{if lt $i ($detailrows|maxIndex)}},{{end}}
{{- end}} 
from {{.Name}}{{.DBLink}} t
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
&{{$c}}
{{- end}}{{- end}}{###}
{{- end}}

//Query{{.Name|upperName}}Count 获取{{.Desc}}列表条数
const Query{{.Name|upperName}}Count = {###}
select count(1)
from {{.Name}}{{.DBLink}} t
where
{{- if eq ($queryrows|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$queryrows -}}
{{if eq ($c.Type|codeType) $time}}
and t.{{$c.Name}}>=to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')
and t.{{$c.Name}}<to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')+1
{{- else if  and (gt $c.Len $length) (eq ($c.Type|codeType) $string)}}
and t.{{$c.Name}} like '%' || @{{$c.Name}} || '%'
{{- else}}
&t.{{$c.Name}}{{end}}
{{- end}}{{end}}
{###}

//Query{{.Name|upperName}} 查询{{.Desc}}列表数据
const Query{{.Name|upperName}} = {###}
select 
	TAB1.*
from (select L.*  
	from (select rownum as rn,R.* 
		from (
			select 
			{{- range $i,$c:=$listrows}}
			{{if eq ($c.Type|codeType) $time }}to_char(t.{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')	{{$c.Name}}{{else -}}
				t.{{$c.Name}}{{end}}{{if lt $i ($listrows|maxIndex)}},{{end}}
			{{- end}} 
			from {{.Name}}{{.DBLink}} t
			where
			{{- if eq ($listrows|len) 0}}
				1=1
			{{- else -}}
			{{- range $i,$c:=$queryrows -}} 
			{{if eq ($c.Type|codeType) $time}}
				and t.{{$c.Name}}>=to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')
				and t.{{$c.Name}}<to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss')+1
			{{- else if and (gt $c.Len $length) (eq ($c.Type|codeType) $string)}}
				and t.{{$c.Name}} like '%' || @{{$c.Name}} || '%'
			{{- else}}
				&t.{{$c.Name}}{{end}}
			{{- end}}{{end}}
			{{- if gt ($order|len) 0}}
			order by {{range $i,$c:=$order}}t.{{$c.Name}}{{if $c.comma}},{{else}} desc{{end}}{{end}}
			{{- else}}
			order by {{range $i,$c:=pks}}t.{{$c}} desc{{end}}
			{{- end}}
			) R 
	where rownum <= @pi * @ps) L 
where L.rn > (@pi - 1) * @ps) TAB1{###}
{{end}}

{{- if  and $ismysql (gt ($updaterows|len) 0)}}
//Update{{.Name|upperName}} 更新{{.Desc}}
const Update{{.Name|upperName}} = {###}
update 
{{.Name}}{{.DBLink}} 
set
{{- range $i,$c:=$updaterows}}
	{{$c.Name}}=@{{$c.Name}}{{if lt $i ($updaterows|maxIndex)}},{{end}}
{{- end}}
where
{{- if eq ($pks|len) 0}}
	1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}{{end}}{###}
{{end -}}

{{- if and $isoracle (gt ($updaterows|len) 0)}}
//Update{{.Name|upperName}} 更新{{.Desc}}
const Update{{.Name|upperName}} = {###}
update 
{{.Name}}{{.DBLink}} 
set
{{- range $i,$c:=$updaterows}}
	{{if eq ($c.Type|codeType) $time}}{{$c.Name}}=to_date(@{{$c.Name}},'yyyy-mm-dd hh24:mi:ss'){{else -}}
	{{$c.Name}}=@{{$c.Name}}{{end}}{{if lt $i ($updaterows|maxIndex)}},{{end}}
{{- end}}
where
{{- if eq ($pks|len) 0}}
	1=1
{{- else -}}
{{- range $i,$c:=$pks}}
	&{{$c}}
{{- end}}{{end}}{###}
{{end -}}

{{- if gt ($deleterows|len) 0}}
//Delete{{.Name|upperName}} 删除{{.Desc}}
const Delete{{.Name|upperName}} = {###}
delete from {{.Name}}{{.DBLink}} 
where
{{- if eq ($pks|len) 0}}
1=1
{{- else -}}
{{- range $i,$c:=$pks}}
&{{$c}}
{{- end}}{{end }}{###}
{{end}}
`
