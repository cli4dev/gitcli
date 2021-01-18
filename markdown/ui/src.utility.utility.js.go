package ui

const srcUtilityUtilityJS = `
import enums from './enum'
import http from './http'
import conf from './env'

/*
* 初始化注入
* import utility from './utility'
* Vue.use(utility);
* 或传入加载配置文件路径
* Vue.use(utility, "../static/env.conf.json");
*/
export default {
    install: function(Vue, path){
        Vue.use(enums);
        Vue.use(http);
        Vue.use(conf, path)
    }
}
`
