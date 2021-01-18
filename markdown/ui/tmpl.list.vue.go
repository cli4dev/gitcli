package ui

const TmplList = `
<template>
  <div ref="main">
    <div class="panel panel-default">
    <div class="panel-body">
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
{{- else}}{{var "tp" $c.Desc|shortName}}{{end -}}{{- end}} 
<div class="form-group input">
<input
type="text"
class="form-control"
v-model="query.input"
onkeypress="if(event.keyCode == 13) return false;"
placeholder="{{- range $i,$c:=vars "tp"}}{{$c}}/{{- end}}"
/>
</div>
<a class="btn btn-success" @click="search">查询</a>
</form>

<el-scrollbar style="height:100%">
<el-table :data="dataList.items" stripe style="width: 100%">

{{- range $i,$c:=.Rows|list}}
{{ if $c.Type|isTime}}
  <el-table-column align="center" width="100" prop="{{$c.Name}}" label="{{$c.Desc|shortName}}">
  <template slot-scope="scope">
    <i class="el-icon-time"></i>
    <span style="margin-left: 10px" v-text="scope.row.{{$c.Name}}"></span>
  </template>
  </el-table-column>
{{else}}
<el-table-column align="center" width="100" prop="{{$c.Name}}" label="{{$c.Desc|shortName}}"></el-table-column>
{{end}}

{{- end}} 
  <el-table-column align="center" label="操作">
	<template slot-scope="scope">
	  <el-button plain type="primary" size="mini" @click="showModal(2,scope.row)">编辑</el-button>
	  <el-button plain type="success" size="mini" @click="userChange(0,scope.row.user_id, scope.row.user_name)" v-if="scope.row.status == 2">启用</el-button>
	</template>
  </el-table-column>
</el-table>
</el-scrollbar>

<div style="padding: 10px 15px;text-align: right;">
<el-pagination
  @size-change="pageSizeChnage"
  @current-change="pageIndexChange"
  :current-page="paging.pi"
  :page-size="paging.ps"
  :page-sizes="paging.sizes"
  layout="total, sizes, prev, pager, next, jumper"
  :total="dataList.count"
></el-pagination>
</div>

</div>
</div>
</div>
</template>


<script>
{{$count:= vars "fields" |maxIndex -}}
export default {
  components: {
   // "bootstrap-modal": require("vue2-bootstrap-modal"),
   // pager: require("vue-simple-pager"),
  },
  data() {
    return {
      paging: {ps: 10, pi: 1,total:0,sizes:[5, 10, 20, 50]},
      query:{ {{- range $i,$c:=vars "fields"}}{{$c}}:""{{if lt $i $count}},{{end}} {{- end}} },
      dataList: {count: 0,items: []},
    };
  },
  created() {
    this.queryData()
  },
 mounted() {},
 methods: {
   pageIndexChange:function(val) {
     this.paging.pi = val;
     this.queryData();
   },
   pageSizeChnage:function(val) {
     this.paging.ps = val;
     this.queryData();
   },
   search: function() {
     this.paging.pi = 1;
     this.queryData();
   },
   queryData:async function(){
     Object.assign(this.query, this.paging);  
     let res = await this.$http.xpost("/{{- range $i,$c:=.Name|rmhd|lower|names}}{{$c}}/{{- end}}query",this.query)
     this.dataList.items = res.items
     this.dataList.count = res.count
   },

  },
}
</script>
<style>
</style>
`
