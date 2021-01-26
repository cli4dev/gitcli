package ui

const srcPagesSystemMenus = `
<template>
  <div id="app">
    <nav-menu
      :menus="menus"
      :copyright="copyright"
      :copyrightcode="copyrightcode"
      :themes="system.themes"
      :logo="system.logo"
      :systemName="system.systemName"
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
  data() {
    return {
      system: {
        logo: "",
        systemName: "",  //系统名称
        themes: "", //顶部左侧背景颜色,顶部右侧背景颜色,右边菜单背景颜色
      },
      copyright: (this.$env.conf.copyright.company || "") + "Copyright©" + new Date().getFullYear() + "版权所有",
      copyrightcode: this.$env.conf.copyright.code,
      menus: [{}],  //菜单数据
      userinfo: {},
      items: []
    }
  },
  components: { //注册插件
    navMenu
  },
  created() {

  },
  mounted() {
    this.$auth.checkAuthCode(this)
    this.getMenu();
    this.getSystemInfo();
    this.userinfo = this.$auth.getUserInfo()
  },
  methods: {
    pwd() {
      this.$http.clearAuthorization();

      var keys = this.$cookies.keys();
      for (var i in keys) {
        this.$cookies.remove(keys[i]);
      }
      var url = this.$env.conf.sso.host + "/" + this.$env.conf.sso.ident + "/changepwd"
      window.location.href = url;
    },
    signOutM() {
      this.$auth.loginout();
    },
    getMenu() {
      this.$auth.getMenus(this).then(res => {
        this.menus = res;
        this.getUserOtherSys();
      });
    },
    //获取系统的相关数据
    getSystemInfo() {
      this.$auth.getSystemInfo().then(res => {
        this.system = res;
      })
    },
    //用户可用的其他系统
    getUserOtherSys() {
      this.$auth.getSystemList().then(res => {
        this.items = res;
      })
    },
  }
}
</script>
`
