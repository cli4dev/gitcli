package ui

const TmplDetail = `
<template>
  <div>
    <el-tabs v-model="tabName" type="border-card" @tab-click="handleClick">
      <el-tab-pane label="{{.Desc|shortName}}" name="{{.Name|varName}}">
        <div class="table-responsive">
          <table :date="info" class="table table-striped m-b-none">
            <tbody class="table-border">
            {{range $i,$c:=.Rows -}}
            {{if eq 0 mod $i 2}}  <tr>  <td>{{end}}               
                  <el-col :span="6">
                    <div class="pull-right" style="margin-right:10px">{{$c.Desc|shortName}}:</div>
                  </el-col>
                  <el-col :span="6">
                    <div class="grid-content" v-text="{{$c.Name}}" ></div>
                  </el-col>
               {{if eq 0 mod $i 2}}  </td> </tr> {{end}}
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
		data() {
		  return {
		  }
		}
	}
	</script>`
