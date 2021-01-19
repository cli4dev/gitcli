package ui

//TmplEditVue 添加编辑弹框页面
const TmplEditVue = `
{{- $rows := .Rows -}}
{{- $empty := "" -}}
<template>
	<el-dialog title="编辑{{.Desc}}"{{if gt ($rows|len) 5}} width="65%" {{- else}} width="25%" {{- end}} @closed="closed" :visible.sync="dialogFormVisible">
		<el-form :model="editData" {{if gt ($rows|update|len) 5 -}}:inline="true"{{- end}} :rules="rules" ref="editForm" label-width="110px">
    	{{- range $i,$c:=$rows|update}}
      {{if $c.Con|TA -}}
			<el-form-item label="{{$c.Desc|shortName}}" prop="{{$c.Name}}">
				<el-input type="textarea" :rows="2" placeholder="请输入{{$c.Desc|shortName}}" v-model="editData.{{$c.Name}}">
        </el-input>
			</el-form-item>
			{{- else if $c.Con|RB }}
			<el-form-item  label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-radio-group v-model="editData.{{$c.Name}}" style="margin-left:5px">
        	<el-radio v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :label="item.value">{{"{{item.name}}"}}</el-radio>
				</el-radio-group>
			</el-form-item>
			{{- else if $c.Con|SL }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}">
				<el-select  placeholder="---请选择---" clearable v-model="editData.{{$c.Name}}" style="width: 100%;">
					<el-option v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name" ></el-option>
				</el-select>
			</el-form-item>
			{{- else if $c.Con|CB }}
			<el-form-item label="{{$c.Desc|shortName}}:" prop="{{$c.Name}}"> 
				<el-checkbox-group v-model="editData.{{$c.Name}}">
					<el-checkbox v-for="(item, index) in {{$c.Name|lowerName}}" :key="index" :value="item.value" :label="item.name"></el-checkbox>
				</el-checkbox-group>
			</el-form-item>
			{{- else if $c.Con|DT  }}
			<el-form-item prop="{{$c.Name}}" label="{{$c.Desc|shortName}}:">
      	<el-date-picker class="input-cos" v-model="editData.{{$c.Name}}" popper-class="datetime-to-date" type="datetime" value-format="yyyy-MM-dd HH:mm:ss"  placeholder="选择日期"></el-date-picker>
      </el-form-item>
      {{- else -}}
      <el-form-item label="{{$c.Desc|shortName}}" prop="{{$c.Name}}">
				<el-input clearable v-model="editData.{{$c.Name}}" placeholder="请输入{{$c.Desc|shortName}}">
				</el-input>
      </el-form-item>
      {{- end}}
      {{end}}
    </el-form>
		<div slot="footer" class="dialog-footer">
			<el-button size="small" @click="dialogFormVisible = false">取 消</el-button>
			<el-button type="success" size="small" @click="edit">确 定</el-button>
		</div>
	</el-dialog>
</template>

<script>
export default {
	name: "{{.Name|lname}}.edit",
	data() {
		return {
			dialogFormVisible: false,    //编辑表单显示隐藏
			editData: {},                //编辑数据对象
      {{- range $i,$c:=$rows|update -}}
      {{if or ($c.Con|SL) ($c.Con|CB) ($c.Con|RB) }}
      {{$c.Name|lowerName}}:this.$enum.get("{{$c.Name|upperName}}"),
      {{- end}}
      {{- end}}
			rules: {                    //数据验证规则
				{{- range $i,$c:=$rows|update -}}
				{{if eq ($c|isNull) $empty}}
				{{$c.Name}}: [
					{ required: true, message: "请输入{{$c.Desc|shortName}}", trigger: "blur" }
				],
				{{- end}}
				{{- end}}
			},
		}
	},
	props: {
		refresh: {
			type: Function,
				default: () => {
				},
		}
	},
	created(){
	},
	methods: {
		closed() {
			this.refresh()
		},
		show() {
			this.dialogFormVisible = true;
		},
		edit() {
			{{- range $i,$c:=$rows|update -}}
			{{- if $c.Con|DT}}
			this.editData.{{$c.Name}} = this.DateConvert("yyyy-MM-dd hh:mm:ss", this.editData.{{$c.Name}})
			{{- end -}}
			{{- end}}
			this.$http.xput("{{.Name|rpath}}", this.editData,{})
			.then(res => {			
				this.dialogFormVisible = false;
				this.refresh()
			})
		},
	}
}
</script>

<style scoped>
</style>
`
