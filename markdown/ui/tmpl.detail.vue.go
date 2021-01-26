package ui

const TmplDetail = `
{{- $string := "string" -}}
{{- $int := "int" -}}
{{- $int64 := "int64" -}}
{{- $decimal := "types.Decimal" -}}
{{- $time := "time.Time" -}}
{{- $len := 32 -}}
{{- $rows := .Rows|detail -}}
<template>
  <div>
    <el-tabs v-model="tabName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="{{.Desc|shortName}}" name="{{.Name|rmhd|varName}}Detail">
        <div class="table-responsive">
          <table :date="info" class="table table-striped m-b-none">
            <tbody class="table-border">
            {{- range $i,$c:=$rows -}}
            {{- if eq 0 (mod $i 2)}}
              <tr>
                <td>
            {{- end}}                 
                  <el-col :span="6">
                    <div style="margin-right: 10px">{{$c.Desc|shortName}}:</div>
                  </el-col>
            {{- if or ($c.Con|SL) ($c.Con|RB) ($c.Con|CB)}}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrEnum("{{(or ($c.Con|moduleCon|firstStr|rmhd) $c.Name)|lower}}")}}</div>
                  </el-col>
            {{- else if and (eq ($c.Type|codeType) $string) (gt $c.Len $len )}}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrEnum("{{(or ($c.Con|moduleCon|firstStr|rmhd) $c.Name)|lower}}")}}</div>
                  </el-col>
          	{{- else if or (eq ($c.Type|codeType) $int64) (eq ($c.Type|codeType) $int) }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrNumberFormat(0)}}</div>
                  </el-col>
            {{- else if eq ($c.Type|codeType) $decimal }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrNumberFormat(2)}}</div>
                  </el-col>
            {{- else if eq ($c.Type|codeType) $time }}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}} | fltrDate}}</div>
                  </el-col>
            {{- else}}
                  <el-col :span="6">
                    <div>{{"{{ info."}}{{$c.Name}}}}</div>
                  </el-col>
            {{- end}}
            {{- if and (eq (mod $i 2) 1) (ne ($rows|maxIndex) $i) }}
                </td>
              </tr>
            {{- end}}
            {{- if eq ($rows|maxIndex) $i }}
                </td>
              </tr>
            {{- end -}}
            {{- end}}            
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
        tabName: "{{.Name|rmhd|varName}}Detail",
        info: {},
      }
    },
    mounted() {
      this.init();
    },
    created(){
      {{- range $i,$c:=$rows|detail -}}
      {{if or ($c.Con|SL) ($c.Con|CB) ($c.Con|RB) }}
        this.$enum.get("{{(or ($c.Con|moduleCon|firstStr|rmhd) $c.Name)|lower}}")
      {{- end}}
      {{- end}}
    },
    methods: {
      init(){
        this.queryData()
      },
      queryData:async function() {
        this.info = await this.$http.xget("{{.Name|rmhd|rpath}}",this.$route.query)
      },
      handleClick(tab) {}
    },
	}
</script>`
