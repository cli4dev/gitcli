package ui

const TmplList = `
<template>
  <div ref="main">
    <div class="panel panel-default">

    <form class="form-inline" role="form">
{{var "tp"}}
{{var "fields"}}
{{- range $i,$c:=.Rows|query}}
{{- if $c.Con|SL}} 
<div class="form-group select">
	<select
	v-model="query.{{$c.Name}}"
	name="{{$c.Name}}"
	class="form-control visible-md visible-lg">
	<option value selected="selected">---请选择{{$c.Desc|shortName}}---</option>
  </select>
  </div>
  {{var "fields" $c.Name}}
{{- else}}
{{var "tp" $c.Desc|shortName}}
{{end -}}


{{- end}} 
<div class="form-group input">
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

<el-scrollbar style="height:100%">
<el-table :data="datalist.items" stripe style="width: 100%">

{{- range $i,$c:=.Rows|list}}
<el-table-column align="center" width="100" prop="{{$c.Name}}" label="{{$c.Desc|shortName}}"></el-table-column>
{{- end}} 
  <el-table-column align="center" label="操作">
	<template slot-scope="scope">
	  <el-button plain type="primary" size="mini" @click="showModal(2,scope.row)">编辑</el-button>
	  <el-button plain type="success" size="mini" @click="userChange(0,scope.row.user_id, scope.row.user_name)" v-if="scope.row.status == 2">启用</el-button>
	</template>
  </el-table-column>
</el-table>
</el-scrollbar>

<div class="page-pagination">
<el-pagination
  @size-change="handleSizeChange"
  @current-change="pageChange"
  :current-page="paging.pi"
  :page-size="paging.ps"
  :page-sizes="pageSizes"
  layout="total, sizes, prev, pager, next, jumper"
  :total="totalpage"
></el-pagination>
</div>

</div>
</div>
</template>


<script>
import pager from "vue-simple-pager";

export default {
  components: {
    "bootstrap-modal": require("vue2-bootstrap-modal"),
    pager: pager,
  },
  data() {
    return {
      paging: {
        ps: 10,
        pi: 1
      },
      pageSizes: [5, 10, 20, 50], //可选显示数据条数
      query:{ 
        {{- range $i,$c:=vars "fields"}}
        {{$c}}="",
        {{- end}} 

      },
      datalist: {
        count: 0,
        items: []
      },
      
    };
  },
}
</script>
<style>
.page-pagination {
  padding: 10px 15px;
  text-align: right;
}
</style>
`
