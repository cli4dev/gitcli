package ui

const TmplDetail = `
{{- $rows := .Rows -}}
<template>
  <div>
    <el-tabs v-model="tabName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="{{.Desc|shortName}}" name="{{.Name|varName}}">
        <div class="table-responsive">
          <table :date="info" class="table table-striped m-b-none">
            <tbody class="table-border">
            {{range $i,$c:=$rows|detail -}}
            {{if eq 0 (mod $i 2)}}
              <tr>
                <td>
            {{- end}}               
                  <el-col :span="6">
                    <div class="pull-right" style="margin-right:10px">{{$c.Desc|shortName}}:</div>
                  </el-col>
            {{- if or ($c.Con|SL) ($c.Con|RB) ($c.Con|CB)}}
                  <el-col :span="6">
                    <div class="grid-content">{{"{{info."}}{{$c.Name}} | fltrEnum("{{$c.Name|varName}}")}}</div>
                  </el-col>
            {{- else}}
                  <el-col :span="6">
                    <div class="grid-content">{{"{{info."}}{{$c.Name}}}}</div>
                  </el-col>
            {{- end}}
            {{- if and (eq (mod $i 2) 1) (ne ($rows|maxIndex) $i) }}
                </td>
              </tr>
            {{- end}}
            {{- if eq ($rows|maxIndex) $i }}
                </td>
              </tr>
            {{- end}}
            {{end -}}            
            </tbody>
          </table>
        </div>
	  </el-tab-pane>
	  </el-tabs>
	</div>
</template>

<script>
	export default {
    data(){
      return {
        tabName:"{{.Name|varName}}",
        info:{},
      }
    },
    created(){
      {{- range $i,$c:=$rows|detail -}}
      {{if or ($c.Con|SL) ($c.Con|CB) ($c.Con|RB) }}
        this.$enum.callback(function(){this.$http.xget("{{$c.Con|gbc|rpath}}/dictionary/get", {})},"{{$c.Name|varName}}")
      {{- end}} 
      {{- end}}
    },
    methods: {
      init(){
        this.queryData()
      },
      queryData:async function() {
        this.info=this.$http.xget(this.$route.query.getpath,this.$route.query)
      },
      handleClick(tab) {}
    },
	}
</script>`
