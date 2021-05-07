# luban
=====

luban 鲁班是也；是一个用于生成新项目代码的工具。当要开发一个全新的项目的时，使用鲁班可以
快速的生产符合平台要求的golang服务。可以让开发人员直接投入到需求开发当中。

**名称由来**： 平台的golang基础框架的名字是takumi, 在日语里是工匠的意思。鲁班是工匠的祖师
爷，由他来创建一个个的takumi工程，似乎比较贴切。

# Get started
-----------

## 安装

在shell中执行
```
go get -u gopkg.mihoyo.com/plat-go/luban
```

## 使用

### step 1
找到有权限在gitlab中创建项目的人，让他在gitlab正确的分组中建立好空的项目，如：

> https://platgit.mihoyo.com/plat/go/community/mission

### step 2

在你本地的`$GOPATH/src`目录下执行
```
luban create plat/go/组名/工程名 项目业务组名
```

> 例如: https://platgit.mihoyo.com/plat/go/community/mission  这个项目。
> 它的工程名是：mission, 该工程的path: plat/go/community
> 所以，当初使用luban初始化这个项目时使用的命令应该是：
>     luban create plat/go/community/mission

### step 3

luban全在当前目录下新建一个以你输入的 `工程名` 命名的目录。项目中必要的代码和文件都在里面了。
```
YourProj
├── Gopkg.toml
├── Makefile
├── README.md
├── api
│   └── YourProj
│       └── YourProj.proto
├── app.yml
├── commands
│   └── server.go
├── deploy.yml
├── handlers
│   └── giraffe.go
├── main.go
├── models
│   ├── mysql.go
│   └── redis.go
└── queue
    └── kafka.go
```

### step 4
按照 gitlab中 空项目的指引将你本地生成的项目代码push到仓库中


### step 5
.......coding hardly...........

**这里你应该遵循平台的git work flow, 如果有不清楚的地方请咨询：parker || 郭忠祥 ||张伟 || 王志刚 **

## 配置文件

* 注：推荐使用配置配置中心管理配置，使用配置中心testing、pre、prod不生产配置文件。[配置中心使用文档](https://www.tapd.cn/22435861/markdown_wikis/view/#1122435861001010574)

takumi应用启动的时候默认读取app.yml文件来加载配置。平台的开发和发布过程涉及到：

- dev
- test
- pre
- prod

4套环境。4套环境的配置文件以不同的文件保存，我们对配置文件的命名作如下约定:


| 环境 |  配置文件名    |  代码分支  |
| ---- | ---- | ------------------ |
| dev | app.yml | 除`develop`,`master`和`release/*`之外的任何分支，一般为`feature/xx`或`hotfix/xx` |
| testing | app.testing.yml | develop |
| pre | app.pre.yml | release/xx |
| prod | app.prod.yml | release/xx 并且打了tag |

各套环境对应的配置文件的格式应该保持一致。

在发布的时候，ci/cd脚本会根据发布的目标环境选取对应的配置文件将其改名为app.yml，并放置到服务进程所在的目录以使配置生效。

## ci/cd

游戏平台服务的发布流程基于gitlab的ci/cd。ci/cd job的执行需要一些预定义的环境变量。

通常每个project只需要在`settings` -> `CI / CD` 界面中展开`Environment variables`设置项增加以下环境变量：

| 变量名 | 值     |  说明  |
| ---- | ---- | ---- |
| BK_MODULE_ID_PROD_ENV     |  int    |  生产环境的模块ID, 请咨询运维同学   |
| BK_GW_MODULE_ID_PROD_ENV     |  int    |  生产环境的gateway模块ID, 请咨询运维同学   |

当新建一个group时，你需要在group的settings中预定义一些必要的环境变量。参见：[init_group.md](init_group.md)

gitlab中的环境变量按 `yml文件` > `project settings` >  `sub group` > `group` 的顺序覆盖生效。以上预定义的环境变量也是按这个逻辑，将一些共用的，提取到root group和各级sub group中。
如果项目中有不同的要求，可以在项目的`settings` -> `CI / CD` 界面中展开`Environment variables`中增加同名的变量设置来覆盖预定义的环境变量。

### 指定打包内容
默认情况下，go服务会打包编译出来的可执行文件，*.yml （即框架的配置文件和服务自定义的配置文件）。

当一个服务有需要额外打包发布的文件时，可以在仓库的根目录下，创建一个名为：`manifest`文件。文件中的每一个行纪录需要打包的文件的路径（**必须是基于仓库根目录的相对路径**）

如果服务需要将proto文件发布到gateway。需要在仓库的根目录下，创建一个名为：`api.manifest`文件。每一行指定一个proto文件的路径。

### 指定发布目标module id
当一个服务完成开发准备上线的时候，需要和运维同学申请生产环境的运行资源。得到一个module_id。这个ID需要在project的`setting` - `ci/cd ` 中填写，展开`Environment variables`设置选项。

以下三个变量，分别用于指定测试，预发布和生产环境的module:

- BK_MODULE_ID_TEST_ENV
- BK_MODULE_ID_PRE_ENV
- BK_MODULE_ID_PROD_ENV

现阶段所有的服务共用一套 测试环境和预发布环境。所以 `BK_MODULE_ID_TEST_ENV`和`BK_MODULE_ID_PRE_ENV`已经在group中预设置好。一般只需要在project中设置`BK_MODULE_ID_PROD_ENV`。

*tips: project中的变量值会覆盖group中的值。如果项目有独立的测试环境，只要在项目的设置中填写测试环境的值就好，不可以去修改group中的值。预发布的情况同理*

## 注意事项 to Dev
- 使用dep管理包依赖
- proto产生的 pb.go文件应该纳入git repo
- vendor应该纳入git repo
- 在项目的根目录提供Makefile，Makefile需包含以下目标
	- app 编译产生可执行文件
	- lint 对代码进行静态检查
	- unittest 对代码进行单元测试
- 确保代码被git clone之后执行`make`即可完成编译


关于`ci/cd`流程的介绍，请看这里：[go微服务cicd流程说明](https://www.tapd.cn/22435861/markdown_wikis/view/#1122435861001007147)
