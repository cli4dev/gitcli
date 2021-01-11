package ui

const srcMainJS = `

import "jquery"
import "bootstrap"
 
import Vue from 'vue'
import App from './App'
import router from './router'
import store from './store'

import {
    post,
    fetch,
    put
} from './utility/http'
import VueCookies from 'vue-cookies'
Vue.use(VueCookies);

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(ElementUI);


//定义全局变量
Vue.prototype.$http = {get:fetch,post:post,put:put}

Vue.config.productionTip = false;


router.beforeEach((to, from, next) => {
    /* 路由发生变化修改页面title */
    if (to.title) {
      document.title = to.title
    }
    next()
  })


  /* eslint-disable no-new */
new Vue({
    el: '#app',
    store,
    router,
    components: {
        App
    },
    template: '<App/>'
});

`
