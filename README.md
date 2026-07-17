# Go+Layui 快速开发框架

一款基于 **Gin + Gorm + Redis + Zap + Layui** 搭建的轻量级后台快速开发框架，结构清晰、开箱即用，适合快速开发**后台管理系统、电商后台、订单管理、CRM、SAAS**等中小型业务系统。

- **GitHub**：https://github.com/tkdaojia/ego
- **Gitee**：https://gitee.com/cndouya/egoweb

## 🔗 应用演示
> **基于Ego框架研发的项目管理系统**
**演示地址：**[http://wuji.wenxiaoxi.cc/](http://wuji.wenxiaoxi.cc/)
- **测试账号：** ego
- **登录密码：** helloego

## ✨ 技术栈
- **Web框架**：Gin
- **ORM框架**：Gorm
- **缓存中间件**：Redis
- **日志组件**：Zap
- **前端界面**：Layui

## 📌 环境要求
- Go **1.18+**
- MySQL **8.0**
- Redis 任意稳定版本

## 📁 项目目录结构
适配本 Go 后端 + Layui 模板项目，标准企业级架构：
```Plain Text
ego/
├── data/               # 附件与日志
├── src/                # 框架后端代码
│   ├── api/            # 核心公共接口
│   ├── boot/           # 系统引导（路由、全局变量、中间件）
│   ├── service/        # 业务应用层（按应用模块划分）
│   └── utils/          # 工具包
├── static/             # 前端静态资源包（Layui, CSS, JS）
├── tool/               # 自动化工具包
├── web/                # 前端模版页面（HTML）
├── config.example.yaml # 配置示例文件 (需改名为config.yaml)
├── main.go             # 项目入口文件
├── go.mod              # Go依赖管理
├── go.sum
└── README.md
```

### 📖 核心文件说明
| 目录/文件 | 功能说明 |
| --- | --- |
| `src/boot/router/sys_api.go` | 路由管理，统一注册所有接口与自动化应用路由 |
| `config.yaml` | 系统全局配置，包含数据库、Redis、日志、服务端口等 |
| `main.go` | 项目启动入口，初始化路由、配置、数据库、Redis |

---

## 🚀 安装部署教程
### 1. 配置文件初始化
将根目录示例配置文件重命名：
```bash
config.example.yaml  -->  config.yaml
```
打开 `config.yaml`，修改 **MySQL、Redis、端口** 等连接信息为本地环境。

### 2. 安装项目依赖
```bash
go mod tidy
```

### 3. 自动初始化数据库
执行内置安装脚本，自动创建数据表 + 初始数据：
```bash
cd tool
go run install.go
```

### 4. 启动项目
返回项目根目录执行：
```bash
go run main.go
```
启动成功后，通过配置的端口访问后台系统。

---

## ✨ 创建新应用标准指南
开发一个新应用（例如 `erp` 或 `crm`）遵循以下标准化流程规则：

### 1. 目录职责划分
> **目录职责划分：**
> * **`module/` (Module - 业务组件层)**：存放该应用专属的高级业务逻辑封装、第三方接口对接、或抽取出的复杂公共算法，作为 `do` 层的强力后盾。
> * **`act/` (Action - 动作入口层)**：负责一级、二级路由的分发。拦截并解析 URL 中的 `?act=xxx` 和 `?do=xxx` 参数，精准导向具体的执行函数。
> * **`do/` (Data Operation - 数据操作层)**：核心业务逻辑与执行层。负责解析具体的 `do` 动作、进行数据绑定与校验（`ShouldBindJSON`）、执行数据库增删改查（GORM）、记录审计日志并最终返回页面或 JSON 响应。

### 2. 挂载自动化路由
打开 `src/boot/router/sys_api.go`，在 `InitApiRouter` 中注册你的应用路由：
```go
AppRouter := Router.Group("/myroute").Use(middleware.VerifyAuthApp())
{
AppRouter.GET("/", app.ApiAppGet)
AppRouter.POST("/", app.ApiAppPost)
}
```
*系统会自动将其挂载到 `/api/app/app.go`（可以自己换文件），定义自己的module层*

### 3. 后台配置（权限与菜单）
1. 登录后台 -> **后台设置** -> **系统组** -> **新建应用名称**（如：ERP管理系统）。
2. 进入 **权限管理/菜单管理**，创建对应的系统模块、菜单页面。
3. 在权限管理中，将对应的菜单和模块权限分配给指定的用户组。

### 4. 编写业务代码并重启
在生成的 `act` 和 `do` 中实现你的业务逻辑。开发完成后，重启 `main.go` 即可直接在 Layui 前端通过 Ajax 异步调用该接口。

---

## 💬 技术交流 & 框架答疑
> 项目使用、二次开发、定制需求、Go/Layui技术问题均可联系作者交流；框架持续迭代更新，入群获取最新版本、开发文档、示例源码。
### 联系方式
扫码或手动添加作者微信（socitys 备注：ego框架），交流进技术交流群：
<div style="text-align: left; margin: 10px 0;">
  <img src="http://www.wenxiaoxi.cc/images/wxx.png" style="width: 200px; height: auto;" alt="作者微信二维码">
</div>
