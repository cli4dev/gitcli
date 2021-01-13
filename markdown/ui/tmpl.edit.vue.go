package ui

//TmplEditVue 添加编辑弹框页面
const TmplEditVue = `
<template>
       <bootstrap-modal
        ref="editPage"
        :need-header="true"
        :need-footer="true"
        :closed="resetSys"
        no-close-on-backdrop="true"
      >
        <div slot="title" v-if="isAdd == 0">添加{{.Desc|shortName}}</div>
        <div slot="title" v-if="isAdd != 0">编辑{{.Desc|shortName}}</div>
        <div slot="body">
          <div class="panel panel-default">
            <div class="panel-body">
              <form role="form" class="ng-pristine ng-valid ng-submitted height-min">
              {{- range $i,$c:=.Rows|addAndEdit}}               
              {{-- if $c|input}}
                <div class="form-group">
                  <label>{$c.Desc|shortName}</label>
                  <input
                    name="{{$c.Name|varName}}"
                    type="text"
                    class="form-control"
                  {{if $c|require}}v-validate="'required'" {{end}}
                    v-model="{{.Name|varName}}.{{$c.Name|varName}}"
                    placeholder="请输入{{.Desc|shortName}}"
                   {{if $c|require}}required {{end}}                    
                    maxlength="{{$c|len}}"
                  />
                  {{- end}}
                  {{if $c|require}}   
                  <div class="form-heigit">
                    <span v-show="errors.first('{{$c.Name|varName}}')" class="text-danger">{{.Desc|shortName}}不能为空！</span>
                  </div> 
                  {{end}}
                </div> 
              {{-- else if $c|select}}}




              {{-- end}}
                <div class="form-group">
                  <div>
                    <label>邮箱</label>
                  </div>                 
                  <select class="email-select" v-model="userInfo.email_suffix">
                    <option selected="selected" value="@100bm.cn">@100bm.cn</option>
                    <option value="@hztx18.com">@hztx18.com</option>
                  </select>                 
                </div>
                <div class="form-group">
                  <label>扩展参数/label>
                  <textarea
                    maxlength="1000"
                    name="ext_params"
                    style="resize:none"
                    rows="5"
                    type="text"
                    class="form-control"
                    v-model="userInfo.ext_params"
                    placeholder="扩展参数"
                  ></textarea>
                </div>              
                <div class="form-group" v-if="isAdd == 1">
                  <label class="checkbox-inline">
                    <input id="statuscheck" type="checkbox" />是否启用
                  </label>
                </div>
              </form>
            </div>
          </div>
        </div>
        <div slot="footer">
          <el-button size="small" @click="onClose">取消</el-button>
          <el-button type="success" size="small" @click="submitUser">提交</el-button>
        </div>
      </bootstrap-modal>
</template>

<script>
export default {
    data() {
         return {
             {{.Name|shortName}}:{},
         }
    }
}
</script>`
