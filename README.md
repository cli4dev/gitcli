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
 ### 五、根据markdown创建项目
 ##### 创建前端vue项目
 1.创建vue项目文件
 ```
 gitcli ui create web
 ```
 2.创建vue项目页面文件
 ```
 gitcli ui page md文档路径 [输出文件路径] -table [指定表名] -f -cover
 ```
 ##### 创建后端hydra项目
 1.创建web服务文件
 ```
 gitcli app create webserver
 ```
 2.创建服务层和sql文件
 ```
 gitcli app service md文档路径 [输出文件路径] -table [指定表名] --exclude [排除表名] -f -cover
 ```
 ##### 数据字典约束配置
 ```
PK:主键
SEQ:SEQ，mysql自增，oracle序列
C: 创建数据时的字段
R: 单条数据读取时的字段 
U: 修改数据时需要的字段
D: 删除，默认为更新字段状态值为1，D[(更新状态值)]
Q: 查询条件的字段
L：(前端页面)列表里列出的字段
OB：查询时的order by字段；默认降序； OB[(顺序)]，越小越先排序
DI: 字典编号，数据表作为字典数据时的id字段
DN: 字典名称，数据表作为字典数据时的name字段
SL: "select"      //表单下拉框,默认使用dds字典表枚举,指定表名的SL[(字典表名)]
CB: "checkbox"    //表单复选框,默认使用dds字典表枚举,指定表名的CB[(字典表名)]
RD: "radio"       //表单单选框,默认使用dds字典表枚举,指定表名的RB[(字典表名)]
TA: "textarea"    //表单文本域
DATE: "date-picker" //表单日期选择器
DTIME: "datetime-picker" //表单日期时间选择器

//C,R,U,Q,L子约束
f:前端过滤器，L(f:过滤器参数)
 ```

 
 
 
 