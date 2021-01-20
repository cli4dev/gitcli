package ui

const srcUitlityEnvJS = `
import $ from 'jquery';
//注入初始化
export default {
    install: function(Vue, path = "../public/env.conf.json"){
        Vue.prototype.$env = new Env(path)
    }
}

/*
* Env对象使用时须通过引用并进行初始化
* import evn from './evn'
* Vue.use(evn);
* 或 可配置加载文件地址(需json格式的文件)
* Vue.use(evn,"../static/env.conf.json")
*/
function Env(path) {
    Env.prototype.Conf = {}
    $.ajaxSettings.async = false; //同步
    $.getJSON (path, function (data){        
        if(!data){
            return
        }
        Object.assign(Env.prototype.Conf, data)
    });  
}

/*
*配置数据加载
*await this.$env.load(async function(){
*   var ress = await that.$http.xpost("/dds/dictionary/get", { dic_type: "operate_action" }, "", false, "") || {}
*   return ress[0]
*})
*/
Env.prototype.load = async function (f) {
    if (typeof f !== "function"){
        return
    }
    let conf = await f() || {}
    
    return Object.assign(Env.prototype.Conf, conf)   
}

`
