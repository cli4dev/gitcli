package ui

const srcRouterIndexJS = `
import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);
export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'menus',
      component: () => import('../pages/system/menus.vue'),
      children:[
        // {
        // path: 'index',
        // name: 'index',
        // component: () => import('../pages/system/index.vue'),
        // meta: { title: "首页" }
        // },
      ]
    }
  ]
})

`
