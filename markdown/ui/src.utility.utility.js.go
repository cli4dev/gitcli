package ui

const srcUtilityUtilityJS = `
import enums from './enum'
import http from './http'
import conf from './env'
import { trim, isPhoneNumber, isEmailNumber, cardNumberFormat, phoneFormat} from './filter'

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
        Vue.use(conf, path);
        Vue.prototype.$utility = {
            trim :trim,
            isPhoneNumber:isPhoneNumber,
            isEmailNumber:isEmailNumber,
            phoneFormat:phoneFormat,
            cardNumberFormat:cardNumberFormat
            } 
    }
}
`
