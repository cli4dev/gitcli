package ui

const srcPagesSystemMenus = `

<template>
  <div id="app">
    <nav-menu
      :menus="menus"
      :copyright="copyright"
      :copyrightcode="copyrightcode"
      :themes="themes"
      :logo="logo"
      :systemName="systemName"
      :userinfo="userinfo"
      :items="items"
      :pwd="pwd"
      :signOut="signOutM"
      ref="NewTap"
    >
    </nav-menu>
  </div>
</template>

<script>
  import navMenu from 'nav-menu'; // 引入
  export default {
    name: 'app',
    data () {
      return {
        logo: "",
        copyright: "" + "Copyright©" + new Date().getFullYear() +"版权所有",
        copyrightcode: "" ,//"蜀ICP备20003360号",
        themes: "bg-danger|bg-danger|bg-dark dark-danger", //顶部左侧背景颜色,顶部右侧背景颜色,右边菜单背景颜色
        menus: [{}],  //菜单数据
        systemName: "订单交易系统",  //系统名称
        userinfo: {name:'admin',role:"管理员"},
        indexUrl: "/",
        items:[]
      }
    },
    components:{ //注册插件
      navMenu
    },
    created(){
    },
    mounted(){
      this.menus=[
    {
        "name":"日常管理",
        "icon":"-",
        "path":"-",
        "children":[
            {
               
                "name":"交易管理",
                "is_open":"1",
                "icon":"fa fa-line-chart text-danger",
                "path":"-",
                "children":[
                    {
                        "name":"交易订单",
                        "icon":"fa fa-user-circle text-primary",
                        "path":"/trade/order"
                    },
                    {
                        "name":"出库订单",
                        "icon":"fa fa-users text-danger",
                        "path":"/trade/delivery"
                    }
                ]
            },
            {
                "name":"商户管理",
                "is_open":"1",
                "icon":"fa fa-tasks text-info-lter",
                "path":"-",
                "children":[
                    {
                        "name":"商户信息",
                        "icon":"fa fa-folder text-success",
                        "path":"/merchant/info"
                    }
                ]
            }
        ]
    }
]
      this.setDocmentTitle();
    },
    methods:{
      pwd(){
        this.$sso.changePwd();
      },
      signOutM() {
        this.$sso.signOut();
      },
      setDocmentTitle() {
        document.title = this.systemName;
      }
    
    }
  }
</script>

<style scoped>
</style>

`
