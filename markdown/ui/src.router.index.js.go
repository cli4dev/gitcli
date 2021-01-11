package ui

const srcRouterIndexJS = `
import Vue from 'vue';
import Router from 'vue-router';
import menus from '@/pages/system/menus';
import index from '@/pages/system/index';

Vue.use(Router);
export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'menus',
      component: menus,
      children:[{
        path: 'index',
        name: 'index',
        component: index,
        titile:"首页"
      }]
    }
  ]
})

`
