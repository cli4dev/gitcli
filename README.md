# gitcli

对gitlab仓库的分组目录提供clone,pull,reset等操作指令

### 一、克隆项目

  gitcli clone gitlab.100bm.cn/micro-plat/fas/apiserver

```sh
> gitcli clone gitlab.100bm.cn/micro-plat/fas/apiserver
get clone https://gitlab.100bm.cn/micro-plat/fas/apiserver 
    /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver
```


### 二、克隆分组下所有项目

 gitcli clone gitlab.100bm.cn/micro-plat/oms

```sh
> gitcli clone gitlab.100bm.cn/micro-plat/oms
get clone https://gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-web /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-web
get clone https://gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-api /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/oms/release/lxl/mgrsys/oms-api

```


### 三、拉取分组下所有项目的指定分支

 gitcli clone gitlab.100bm.cn/micro-plat/fas -branch dev

 ```sh
> gitcli pull gitlab.100bm.cn/micro-plat/fas -branch dev
get clone https://gitlab.100bm.cn/micro-plat/fas/docs /home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/docs
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/docs > git pull origin dev:dev
....
 ```

 ### 四、撤销分组下所有项目的修改

 gitcli reset gitlab.100bm.cn/micro-plat/fas -branch dev

 ```sh
> gitcli reset gitlab.100bm.cn/micro-plat/fas -branch dev
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git reset --hard
HEAD 现在位于 1575b23 RPC SDK
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git checkout dev
切换到分支 'dev'
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/fds > git reset --hard
HEAD 现在位于 1575b23 RPC SDK
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git reset --hard
HEAD 现在位于 420caa3 
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git checkout dev
切换到分支 'dev'
/home/yanglei/work/src/gitlab.100bm.cn/micro-plat/fas/apiserver > git reset --hard
HEAD 现在位于 420caa3
....
 ```