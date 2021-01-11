package ui

const TmplQueryVue = `
<form class="form-inline" role="form">
{{var "tp"}}
{{- range $i,$c:=.Rows}}
<div class="form-group query">
{{- if $c.Con|SL}} 
	<select
	v-model="query.{{$c.Name}}"
	name="{{$c.Name}}"
	class="form-control visible-md visible-lg">
	<option value selected="selected">---请选择{{$c.Desc|shortName}}---</option>
	</select>
{{- else}}
{{var "tp" $c.Desc|shortName}}
{{end -}}

</div>
{{- end}} 
<div class="form-group">
<input
type="text"
class="form-control"
v-model="query.input"
onkeypress="if(event.keyCode == 13) return false;"
placeholder="请输入{{var "tp"}}"
/>
</div>
<a class="btn btn-success" @click="searchClick">查询</a>
</form>

`
