# 开发使用指南

## 快速开始

### 创建应用

通过 `蓝鲸开发者中心 -> 创建应用`，填写应用 ID & 名称，应用模板选择 `Go 开发框架 -> Golang 开发框架（Gin）` 并配置代码仓库（如：工蜂）。

### 下载模板 & 初始化

在创建应用成功页面，点击 **初始化模板链接** 下载开发框架模板到本地并解压，参考页面指引完成 Git 仓库 **初始化 & 代码提交**。

### 应用部署

点击指引页面 **应用部署** 前往部署管理页面，点击部署按钮，耐心等待预发布环境部署完成，点击 **访问** 查看示例功能。

## 本地开发

Golang 项目推荐使用 `Makefile` 来管理常用命令，开发者可以查阅该文件以获取更多有用的命令。

### 目录说明

```shell
.
├── ChangeLog.md
├── Dockerfile
├── Makefile
├── README.md
├── app_desc.yaml             # 蓝鲸应用（SaaS）配置文件
├── cmd
│   ├── extract_i18n_msgs.go  # extract-i18n-msgs 命令，用于从源码 / 模板中提取国际化数据
│   ├── init.go
│   ├── init_data.go          # init-data 命令，用于首次部署时初始化数据
│   ├── make_migration.go     # make-migration 命令，用于生成数据库版本文件（需手动实现具体变更内容）
│   ├── migrate.go            # migrate 命令，用于执行数据库表结构变更
│   ├── root.go
│   ├── scheduler.go          # scheduler 命令，用于启动定时任务服务器（必须单实例）
│   ├── version.go            # version 命令，用于查阅目前服务的版本信息
│   ├── view_config.go        # view-config 命令，用于查阅目前服务加载的配置信息
│   └── webserver.go          # webserver 命令，用于启用提供 API & 前端页面的 Web 服务
├── configs
│   └── config.yaml           # 配置参考模板
├── go.mod
├── go.sum
├── main.go
├── pkg
│   ├── account             # 登录 & 认证相关
│   │   └── ...
│   ├── apis                # Web API
│   │   ├── asynctask         # 异步任务示例
│   │   │   ├── handler         # API 处理逻辑
│   │   │   │   └── ...
│   │   │   ├── router.go       # 子路由注册
│   │   │   └── serializer      # API 出/入数据结构体定义
│   │   │       └── ...
│   │   ├── basic             # 基础 API（如 ping，healthz，version 等）
│   │   │   └── ...
│   │   ├── cache             # 缓存使用示例
│   │   │   └── ...
│   │   ├── cloudapi          # 云 API 调用示例
│   │   │   └── ...
│   │   ├── crud              # 数据库操作（CRUD）示例
│   │   │   └── ...
│   │   └── objstorage        # 对象存储服务示例
│   │       └── ...
│   ├── async               # 简单的异步任务服务封装
│   │   └── ...
│   ├── cache              # 基础类设施（依赖的外部服务）
│   │   ├── memory           # 内存缓存（基于 freecache)
│   │   │   └── ...
│   │   └── redis             # Redis 缓存
│   │       └── ...
│   ├── common              # 项目通用的设置
│   │   ├── probe             # 健康探针
│   │   │   └── ...
│   │   └── ...
│   ├── config              # 配置建模 & Loader
│   │   └── ...
│   ├── infras              # 基础类设施（依赖的外部服务）
│   │   ├── cloudapi          # 云 API 相关封装
│   │   │   └── cmsi            # 通用消息发送服务 API 封装
│   │   │       └── ...
│   │   ├── database          # 数据库接入（Mysql + GORM）
│   │   │   └── ...
│   │   ├── objstorage        # 对象存储（BkRepo)
│   │   │   └── ...
│   │   ├── redis             # Redis 服务
│   │   │   └── ...
│   │   └── otel              # OpenTelemetry
│   │       └── ...
│   ├── logging             # 日志相关
│   │   └── ...
│   ├── middleware          # 自定义 Gin 中间件
│   │   └── ...
│   ├── migration           # 数据库版本控制
│   │   └── ...
│   ├── model               # 数据库模型（GORM）
│   │   ├── ...
│   │   └── types.go          # 自定义字段
│   ├── router              # web 服务路由主入口
│   │   └── ...
│   ├── utils               # 项目工具集
│   │   ├── crypto            # 加解密工具
│   │   │   └── ...
│   │   ├── envx              # 环境变量工具
│   │   │   └── ...
│   │   ├── ginx              # Gin 框架工具
│   │   │   └── ...
│   │   ├── testing           # 单元测试工具
│   │   │   └── ...
│   │   └── uuidx             # uuid 工具
│   │       └── ...
│   ├── version             # 项目版本信息
│   │   └── ...
│   └── webfe               # 前端页面路由
│       └── ...
├── static
│   └── image               # 静态文件 - 图片
│       └── ...
├── tailwind.config.js
└── templates
    └── webfe               # 前端页面 HTML 模板
        └── ...
```

### 配置说明

目前框架配置模型所在的位置是 `pkg/config/types.go`，其中的字段都有相应的注释，开发者可以查阅以获得更多有用的信息。

#### Config 中的 Platform，Service，Biz 字段分别是什么定位？

- Platform 用于承接平台提供给 SaaS 的内置环境变量
    - 开发者可以在 `应用详情页 -> 模块配置 -> 环境变量 -> 内置环境变量` 中查看字段对应的配置值
    - 比较特殊的是 `Platform.Addons` 部分，这需要在 `模块配置 -> 增强服务` 中查看字段对应的配置值
- Service 用于存放一些 web / scheduler 服务本身运行所需要的配置，如端口、日志、密钥等
- Biz 用于存放 SaaS 开发者业务需要的一些配置，如 `MaxLimit，ExpireTime` 等

#### 配置项 AppSecret 是什么？如何获取？

AppSecret 是开发者中心分配给每个应用的密钥，与 AppID 搭配使用，开发者可以在 `应用详情页 -> 应用配置 -> 基本信息 -> 密钥信息` 查看并管理应用密钥。

#### 配置项 HealthzToken 和 MetricToken 是什么？如何获取？

由于健康探针访问 `/healthz` API 或 监控通过 `/metrics` 拉取指标数据时是没有用户登录态的，因此这些接口走的是 Token 鉴权的方式（具体实现：`pkg/middleware/token.go`）

我们推荐开发者使用足够长的随机字符串来作为 Token（建议至少 32 位）

### 开发环境搭建

#### 配置 Golang 开发环境

我们推荐使用 JetBrain 的 Goland IDE 来开发 Golang 项目，通过 Goland 打开项目后，通过 `Goland -> Settings -> Go -> GOROOT` 设置使用最高的 Go 1.22 版本（如 Go 1.22.5）。

在 IDE 中打开 Terminal，输入 `go version` 确定 Go 版本是否正确，而后执行 `make tidy` 下载项目所需要的 Golang 依赖包（此步需要耐心等待）。

#### 完善项目配置

本框架支持从环境变量 / 配置文件读取项目配置（`pkg/config/loader.go`）；在本地开发时，我们推荐使用配置文件来管理配置，开发者可以基于 `configs/config.yaml` 进行修改以获取可用的配置文件。

#### 配置 Hosts & BkDomain

本框架已经接入蓝鲸统一登录（`pkg/account`），开发者需要配置本地开发域名以共享身份凭证信息（在 Cookies 中）。

举个例子：如果你所使用的蓝鲸开发者中心域名为 `bkpaas.example.com`，则推荐你使用 `appdev.example.com` 来作为你本地开发的域名。

如何你使用的是 MacOS / Linux 系统，可以通过修改 `/etc/hosts` 来配置 Hosts，参考命令如下：

```shell
$ sudo vim /etc/hosts
# 新添加一行: 127.0.0.1 appdev.example.com
```

注：如果你使用的是 Windows 系统，则需要修改 `C:\Windows\System32\drivers\etc\hosts` 文件。

同时，你还需要修改配置项 `platform.bkDomain` 或 `BK_DOMAIN` 环境变量的值（如 `example.com`)，该配置会影响到跨域 & CSRF 防护相关中间件。

#### 数据库初始化

在启动服务之前，我们需要对本地的数据库进行初始化，开发者可以通过各类数据库可视化软件（如 MySQL Workbench，Navicat 等）来创建数据库，需要注意有：

- 数据库配置 & 名称应该与 `config.yaml` 文件的 `platform.addons.mysql` 中的配置一一对应
- 数据库字符集需要使用 `utf8mb4`（或 `utf8`） 以避免在添加中文数据时候出错

除了使用可视化软件外，你也可以通过在 mysql shell 中执行 SQL 语句来为应用创建一个数据库：

```shell
CREATE DATABASE `gin-demo` DEFAULT CHARACTER SET = `utf8mb4` DEFAULT COLLATE = `utf8mb4_general_ci`;
```

#### 启动 web & scheduler 进程

```shell
# web 页面 & API 服务
$ go run main.go webserver --conf=configs/config.yaml
# 定时任务调度器进程
$ go run main.go scheduler --conf=configs/config.yaml
```

你可以通过以下方式，判断服务是否正常启动：

- 观察控制台输出的标准输出日志
- 浏览器访问 <http://appdev.example.com:8080>

## 进阶功能

### Pre-Commit hook

当使用 Git 来管理项目代码的时候，我们推荐为项目配置 [Pre-Commit hook](https://pre-commit.com/)

参考配置命令如下：

```shell
# 在项目根目录下（与 .pre-commit-config.yaml 文件同层级）
pip install pre-commit
pre-commit install
```

pre-commit 可以在代码提交时，自动帮助开发者完成文档生成（make doc），代码格式化（make fmt），潜在错误检查（make vet）等工作。

### 关于 API 开发

本开发框架一般作为 Web 服务后端使用，即用于为前端（如 Vue）项目提供 API。

在 API 编写上，我们推荐以下实践：

- 尽可能使用 [RESTful 风格](https://www.ruanyifeng.com/blog/2014/05/restful_api.html) 来设计 API URL
- 使用 Request & Response 结构体来限制 API 的输入输出
- 使用 binding & Validate Func 或 [Validator](https://gin-gonic.com/zh-cn/docs/examples/custom-validators/) 来检查参数合法性
- 使用 `pkg/utils/ginx` 中提供的各类工具来获取 UserID、分页参数、设置返回数据等

如果更多参考，可以查阅 `pkg/apis` 包中的框架功能示例

### 用户认证 / 豁免登录

目前开发框架已支持蓝鲸统一登录、太湖（TAI）等多种用户认证方式，提供了获取用户身份 & 登录态的功能，相关代码实现可查阅 `pkg/account`。

#### 如何从 gin context 中获取用户信息？

框架中提供 `ginx.GetUserID` 方法来获取当前登录的用户 ID。注意：调用该方法前，确保相关 URL 访问会经过 `UserAuth` 中间件。

#### 如何实现指定的 URL 豁免登录？

如上所述，目前用户认证 & 登录依赖 `UserAuth` 中间件，只需要对应的 URL（router）不经过该中间件即可。

#### 如何使用太湖（TAI）认证？

具体步骤如下：
1. 开发者需要在太湖上创建应用，并申请自定义域名，同时将该域名作为 PC / 移动站点接入到太湖。
2. 在 `config.yaml` 中的 `service.authTypes` 配置中添加 `Taihu`，示例如：`["Taihu", "BkTicket"]`。
3. 在 `config.yaml` 中配置 `service.taihuAppToken`，其值可在太湖服务页面上的 `应用管理 - 应用概览` 处查询得到。
4. 如果在太湖上的应用配置的不是安全模式（即兼容 / 明文模式），则还需要配置 `service.taihuInsecure` 值为 `true`。
5. 在开发者中心上配置自定义域名并部署服务后，当使用已经接入太湖的域名访问服务时，则会自动通过太湖进行身份认证。

注意：如果使用的是环境变量，则以上配置项对应的环境变量名分别为 `AUTH_TYPES`, `TAIHU_APP_TOKEN`, `TAIHU_INSECURE`。

配置格式形如：`AUTH_TYPES=Taihu,BkTicket, TAIHU_APP_TOKEN=xxxx, TAIHU_INSECURE=false`。

#### 如何同时配置多种认证方式？

开发框架支持同时配置多种认证方式，开发者仅需配置 `service.authTypes` 即可（示例：`["Taihu", "BkTicket"]`），认证逻辑如下：
- 开发框架会按顺序逐个使用 AuthBackend 来尝试获取用户信息，一旦成功获取到用户信息，则停止后续的认证方式。
- 若所有认证方式均未找到用户信息，则会跳转到登录页面，若 **最后一种** 认证方式有提供登录链接，则登录页面上会有 **立即登录** 按钮。

更多参考：

- 框架中间件使用位置：`pkg/router/router.go`
- AuthBackend 具体实现：`pkg/account/`
- Gin 中间件设计：[Custom Middleware](https://gin-gonic.com/zh-cn/docs/examples/custom-middleware/)
- Gin 中间件使用：[Using Middleware](https://gin-gonic.com/zh-cn/docs/examples/using-middleware/)

### CORS 配置

CORS（跨域资源共享）是一种安全机制，允许服务器指定哪些源（域名）可以访问其资源，从而避免暴露资源给非预期的域名服务。

目前开发框架已内置 CORS 中间件（`pkg/middleware/cors.go`），默认情况下，该中间件允许任意域名跨域访问（即 `allowedOrigins: ["*"]`）

**注意：[凭据模式](https://htmlspecs.com/fetch/#cors-protocol-and-credentials) 为 include 时无法命中 `"*"` 的匹配规则，必须指定 `allowedOrigins` 才能放行对应域名的请求**

我们推荐开发者修改配置文件中的 `service.allowedOrigins` 项来限制允许访问的域名（尤其是使用多个域名的前后端分离应用），参考格式如下：

```yaml
service:
  allowedOrigins:
    - http://frontend.example.com
    - http://backend.example.com
```

**注意：填写的值必须包含 scheme（http / https），如果使用的不是 80 / 443 端口则需要显式指定。**

如果你使用的是环境变量来配置，则参考格式如下：

```shell
export ALLOWED_ORIGINS="http://frontend.example.com:6060,http://backend.example.com:8080"
```

除此之外，如果还有修改 AllowHeaders 等配置的需求（例如跨域需要使用自定义头），可以去 `pkg/middleware/cors.go` 中修改相关代码。

### CSRF 配置

开发框架中已内置 CSRF 防护中间件（`pkg/middleware/csrf.go`），默认情况下，该中间件会拦截所有不安全的请求（例如 POST，PUT 等会修改数据的方法，更多参考：[RFC7231](https://datatracker.ietf.org/doc/html/rfc7231#section-4.2.1)），要求请求头中必须包含 CSRF Token。

同时，我们提供 CSRFToken 中间件来把 Token 存入 Cookie 中，开发者可以修改配置文件中的 `service.csrfCookieDomain` 或设置 `CSRF_COOKIE_DOMAIN` 环境变量来设置 CSRF Token 在 Cookie 中生效的域名，例如：

```yaml
service:
  # example.com 表示对所有子域名 (如 demo.example.com/login.example.com) 均生效
  # 也可以设置成具体的子域名，例如：demo.example.com，设置为空字符串 "" 表示为当前域名
  csrfCookieDomain: "example.com"
```

注意：发送 POST / PUT 等不安全请求，需要把从 cookie 中获取的 Token 放入名称为 `X-CSRF-TOKEN` 的 Header 中，同时需要带上 `blueapps-go-csrf` 这个 cookie 用于验证 Token。

若开发者有跨域的需求（例如前后端分离开发）需要在 `pkg/middleware/cors.go` 中修改 CORS 相关配置：在 `AllowHeaders` 中添加 `X-CSRF-Token`，并且配置合理的 `AllowOrigins`。

#### API 请求出现 `Forbidden - CSRF token invalid` 问题排查

1. 检查浏览器 Cookies 中是否存在 **同名但不同域** 的 `blueapps-go-csrf` / `blueapps-go-csrf-token`，若有需清理后重试。
2. 如果是前后端分离开发，则需要检查是否有按照本节说明，合理配置 CORS 并且在请求的 Header 中设置 CSRF Token。
3. 如果不涉及前后端分离开发（使用与示例相同的 GoTemplate 编写 html 页面），则可以参考文件 `templates/web/crud.html` 中的实现，为 axios 设置 CSRFToken。

### 用户访问控制（白名单）

开发框架中内置简单的用户访问权限控制中间件（`pkg/middleware/access_control`），默认情况下 allowedUsers 为空，表示允许任意用户访问。

如果开发者有启用用户访问白名单的需求，可以通过修改 `service.allowedUsers` 配置或 `ALLOWED_USERS` 环境变量，添加 UserID 来启用该中间件，具体示例如下：

```yaml
# 使用配置文件的方式
service:
  allowedUsers:
    - admin
    - userAlpha
```

```shell
# 使用环境变量的方式（注意不能有多余的空格）
export ALLOWED_USERS="admin,userAlpha"
```

### 日志

开发框架默认提供基础的日志功能，开发者可以使用 `pkg/logging/shim.go` 中提供的方法来打印日志，参考示例如下：

```go
import log "bk.tencent.com/blueapps-go/pkg/logging"

func main() {
	// 注：ctx 应该从程序入口一路传递下来（如 cmd 函数中的 context.Background()）
	// 在 gin.HandlerFunc 中可以直接取 c.Request.Context()，中间件会默认添加 Request ID
	ctx := context.Background()
	// 调试日志（需要指定 service.log.level 为 debug 才会输出）
	// 注：如果你不需要格式化字符串，可以使用名称不带 f 后缀的方法
	log.Debug(ctx, "program runs to this point")
	// 普通日志
	log.Infof(ctx, "user %s logged in", "username")
	// 警告日志
	log.Warnf(ctx, "value %d is out of range: [%d, %d]", 10, 0, 9)
	// 错误日志
	log.Errorf(ctx, "connection error: %s", err)
}
```

注:

- 应使用 `pkg/logging/shim.go` 中提供的方法来打印日志，而非标准库 `log` 或 `log/slog`。
- 若开发者对日志性能有更高的要求，可以使用 [uber-go/zap](https://github.com/uber-go/zap) 作为 slog 的 [Handler](https://github.com/samber/slog-zap)。
- 若本地开发时希望输出日志到 console 而非文件，修改配置 `service.log.forceToStdout` 为 `true` 即可。
- gorm & gin 的日志级别不是由配置 `service.log.level` 控制的，而是固定的常量（`GinLogLevel` & `GormLogLevel`），其默认值均为 `warn`，开发者可按需修改。

### ORM

开发框架默认使用 [GORM](https://gorm.io/docs/) 作为与数据库交互的 ORM，这是一个比较成熟的 Golang ORM，简单易上手，也有较好的社区 & 文档支持。

目前 ORM 接入实现的位置是 `pkg/infras/database`，相关的命令（cmd）有 `migrate` 以及 `init_data`；注意：`webserver` & `scheduler` 命令也依赖 ORM 支持。

除 GORM 外，还有许多优秀的 Golang ORM，如 [SQLBoiler](https://github.com/volatiletech/sqlboiler) / [Ent](https://github.com/ent/ent) 等，它们通过代码生成而非反射从而提供了静态类型检查 & 更好的运行性能。

如果你对性能有更进一步的需求，还可以考虑下 [sqlx](https://github.com/jmoiron/sqlx)，这是一个高性能的标准 sql 库增强 & 扩展包，缺点是需要写比较多的 SQL，并且需要自行处理诸如 SQL 注入等的安全问题。

如果你不希望使用 GORM，可以参考上面的建议，自行替换掉 `pkg/infras/database` 中的实现。

### 数据库版本控制

由于我们的开发框架默认采用 GORM，因此我们选择简单可靠的 [gormigrate](https://github.com/go-gormigrate/gormigrate) 来控制数据库的版本。

开发者可以通过执行 `make-migration` 命令来生成新的版本文件（存放目录：`pkg/migration`）。

**⚠ 注意：开发者需要手动在新生成的版本文件中实现 `Migrate` 和 `Rollback` 方法。**

```shell
# 生成新的数据库版本文件
go run main.go make-migration
```

完成新版本文件的编写后，开发者可以通过执行 `migrate` 命令来更新数据库。我们还提供了 `--migration` 参数，允许迁移到指定的数据库版本。

```shell
# 更新数据库版本
go run main.go migrate --conf=configs/config.yaml

# 更新 或 回滚 到指定的数据库版本
go run main.go migrate --conf=configs/config.yaml --migration=20241022_105518
```

如果想了解更多关于数据库版本控制的内容与建议，请参阅 [数据库迁移（Migration）指南](../pkg/migration/README.md)。

### 缓存

开发框架目前支持接入内存和 Redis 两种缓存（`pkg/cache/memory + pkg/cache/redis`），可以查看 `apis/cache` 目录下的代码以获取参考使用方法。

内存缓存基于 [freecache](https://github.com/coocood/freecache) 进行封装，采用预分配内存 + LRU 算法，支持简单的 Set，Get，Del 等操作，具体可参考 [接口文档](https://pkg.go.dev/github.com/coocood/freecache?utm_source=godoc)。

开发者可以通过修改配置 `service.memoryCacheSize` 或环境变量 `MEMORY_CACHE_SIZE` 来闲置缓存使用的内存（默认为 100 MB）。

redis 缓存基于 [go-redis](https://github.com/go-redis/redis) + [go-redis/cache](https://github.com/go-redis/cache) 进行封装，若部署到开发者中心，推荐结合 Redis 增强服务使用，配置位置：`platform.addons.redis`。

开发框架中保留 Redis 的原生实现，开发者可以直接使用 `pkg/infras/redis` 中提供的 Redis 客户端实例 `redis.Client()` 来实现其他基于 Redis 的需求（如：分布式锁）。

### 云 API

云 API 是 **蓝鲸开发者中心** 与 **蓝鲸 API 网关** 联合提供的扩展能力，开发者可在开发者中心中查阅 & 申请相应的 API 权限，并通过 SDK 进行调用。

在 Golang 开发框架中，我们推荐使用蓝鲸 API 网关提供的 [SDK](https://github.com/TencentBlueKing/bk-apigateway-sdks) 来接入并调用 API，目前开发框架提供了对接组件 API `(cmsi.send_mail)` 的示例 `(pkg/infras/cloudapi/cmsi)` 以供开发者参考。

云 API 文档 & 权限申请入口：`蓝鲸开发者中心 -> 应用详情页面 -> 云 API 权限（左侧导航栏）`

### 异步/定时任务

如果查阅 `pkg/async` 包的实现，你就会发现目前开发框架的异步任务示例，并没有使用类似与 Python 中 celery 这样的大型框架，也没有强依赖外部的消息队列（RabbitMQ / Redis）。

这样做的原因是：我们考虑到 Go 原生支持异步（简单的 `go` 关键字即可启动协程跑异步任务），如果直接引入大型异步任务框架会显得过重，也不一定是开发者需要的功能（增加学习成本）。

在定时任务方面，我们目前示例（scheduler）使用的是 `robfig/cron + singleton` 来实现，通过单实例运行来确保不会出现重复执行定时任务的问题。

注意：定时任务会在预定的时间通过 `ApplyTask` 方法下发执行，该方法要求任务执行函数第一个参数 **必须** 是 `context.Context`，最后一个返回值 **推荐** 为 `error`。

你可以查看 `pkg/async/tasks.go` 中的包注释以获得更多的信息 & 建议。

#### 异步任务框架

在开发框架设计阶段，我们调研了使用量比较高的的 Golang 异步任务框架，最后锁定其中两个：

- [machinery](https://github.com/RichardKnop/machinery)（🌟 7.4k）
    - 优点：功能齐全，而且是极少数支持 rabbitmq 作为 broker 的异步框架
    - 缺点：框架偏重，单测覆盖不高，看 issue / commit 维护情况不是很乐观
- [asynq](https://github.com/hibiken/asynq)（🌟 9.2k）
    - 优点：轻量，功能齐全，自带 web ui、支持优先级队列，命令行工具等
    - 缺点：必须强依赖 redis 作为消息队列（不支持其他类型的 broker）

在实现 machinery 接入后，我们发现一些问题：

- 框架太重，对任务函数定义，参数限制都有比较高的要求，开发者需要系统学习该框架，有一定的成本
- 测试发现 redis lock 存在 Bug，这会导致在多副本运行时候，无法通过锁来控制定时任务不会重复执行

因此我们在讨论后，决定框架中的异步任务应该尽可能轻量，仅提供基于 goroutine 的异步任务 & 基于 cron 的定时任务示例。

如果开发者有使用异步任务的需求，可以参考框架文档自行接入，配置项可使用 `config.Platform.Addons`。

### 蓝鲸监控看板

**注：该功能需要应用部署环境（集群）支持使用蓝鲸监控，具体可咨询应用部署环境的维护者 / 助手服务**

Go 开发框架现已接入蓝鲸监控仪表盘，开发者可在成功部署 & 稍等片刻后通过 `访问蓝鲸监控 -> 仪表盘 -> 左侧选择当前应用空间` 查看默认监控看板。

默认监控看板提供以下纬度的指标：

- Gin：HTTP 请求数，耗时，独立访客数，慢请求等
- Go 进程：内存，FD 使用 / 分配情况，Goroutine 数量等
- 系统资源：CPU / 内存使用率，磁盘读写，网络 IO 等

开发者可以通过监控看板更好地了解当前服务的总体状态；针对指标异常的情况，可以配置告警以便及时知晓并处理。

#### 问题排查

##### 没有默认看板

开发者应确认应用服务已经 **成功部署**，若部署失败是不会默认创建仪表盘的。

若应用已成功部署但还是没有看板，则需要联系部署环境的维护者 / 助手服务，确认以下前置准备是否就绪：

- 应用部署的集群已支持蓝鲸监控
- 开发者中心已添加默认看板配置

##### 看板无数据

1. 确认当前应用服务已经成功部署，且各个进程正常运行（首次部署数据上报可能有延迟，需耐心等待）。
2. 若 Gin / Go 进程类指标无数据，则可以先访问已部署的服务，触发数据上报后再进行检查。
3. 检查环境变量 `METRIC_TOKEN` 或配置项 `service.metricToken` 是否与 `app_desc.yaml` 中可观测性配置的 `token` 一致。
4. 联系该部署环境的维护者 / 助手服务协助排查无数据问题。

### 蓝鲸 APM（Otel / OpenTelemetry）

OpenTelemetry 是开源的分布式追踪解决方案，它可以帮助我们更好地了解应用的性能问题，更好地对程序进行优化；开发者可以阅读 [官方文档](https://opentelemetry.io/docs/what-is-opentelemetry/) 以了解相关的背景知识。

蓝鲸 APM（应用性能监控）通过服务之间的调用来分析问题，其完全兼容 OpenTelemetry 协议，支持 Metrics、Logs、Traces、Profiling 等观测数据，实现了主机、容器、事件、告警等数据关联，提供基于应用 / 服务视角的数据观测及故障分析能力。

#### 如何启用蓝鲸 APM

开发框架目前已经支持接入蓝鲸 APM，但由于并非所有环境的开发者中心都支持蓝鲸监控，因此没有默认启用。

如果需要启用蓝鲸 APM，需要前往 `应用详情页 -> 模块配置 -> 增强服务 -> 蓝鲸 APM` 处 **手动启用并重新部署**。

#### 接入 OpenTelemetry 需要怎么做

目前开发框架内置接入 OpenTelemetry 的第三方包有 [Gin](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gin-gonic/gin)，[Resty](../pkg/infras/otel/otel-resty)，[GORM](https://github.com/go-gorm/opentelemetry)，[Redis](https://github.com/redis/go-redis/tree/master/extra/redisotel)，在应用部署且调用后即可在蓝鲸监控中看到对应数据。

如果开发者有为使用的其他第三方包接入 OpenTelemetry 的需求，可以在 [Registry](https://opentelemetry.io/ecosystem/registry/?s=&component=&language=&flag=all) 上搜索对应的 otel 包并参考指引自行接入。

#### 本地开发如何调试

当开发者在本地开发时，可以使用 jaeger 作为数据上报后端（替代蓝鲸 APM），在上线前测试并检查数据上报 & 内容是否符合预期。

```shell
# Docker 一键拉起 jaeger 服务
docker run -d -name jaeger-dev -p 4317:4317 -p 16686:16686 jaegertracing/all-in-one
```

注：本地开发时需修改 `platform.addons.bkOtel.grpcUrl` 或环境变量 `OTEL_GRPC_URL` 为 `http://localhost:4317`（无需配置 token）；在数据上报后，通过浏览器访问 `http://localhost:16686` 查看上报的 Tracing 数据。

#### 其他注意事项

##### 可观测性 x Resty

开发框架使用 resty 作为 HTTP Client，我们可以简单地通过 `resty.New()` 初始化 resty 客户端，但是此时这个客户端不会携带 Request / Span / Trace ID 信息。

如果开发者需要记录 & 关注 resty 的调用链信息，可以参考 `pkg/infras/objstorage/init.go` 中的实现，或者参考以下示例代码：

```golang
package main

import (
  "context"
  "fmt"

  "github.com/go-resty/resty/v2"

  otelresty "bk.tencent.com/blueapps-go/pkg/infras/otel/otel-resty"
  slogresty "bk.tencent.com/blueapps-go/pkg/logging/slog-resty"
)

func main() {
	// 需替换为顶层（如 Gin.Handler）传入的 Context
	ctx := context.Background()

	client := resty.New().
		// 设置定制化 Logger 以在日志中打印 Request / Span / Trace ID
		SetLogger(slogresty.New(ctx)).
		// OpenTelemetry 相关中间件，用于上报 Tracing 数据
		OnBeforeRequest(otelresty.RequestMiddleware).
		OnAfterResponse(otelresty.ResponseMiddleware)
		// 其他 Resty 设置
		//...

	url := "http://bk.example.com/echo"
	params := map[string]string{"foo": "bar"}

	var respData any
	resp, err := client.R().SetContext(ctx).SetResult(&respData).SetQueryParams(params).Get(url)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(resp.String())
}
```

### 数据库字段加密

敏感数据存储时加密是 SaaS 开发中的常见需求，可以避免在数据库中明文存储敏感信息。

开发框架在 `pkg/model/types.go` 中提供支持 AES 对称加密算法的 GORM 自定义数据类型 `AESEncryptedString`，从而允许数据在 DB 中以加密形式存储，在进程中则可以使用其字面值。

使用示例：

```go
type User struct {
	// 非加密字段
	username string `gorm:"type:varchar(128);not null;unique"`
	// 【推荐】加密字段，特殊指定数据类型为 varchar(256)
	email AESEncryptedString `gorm:"type:varchar(256)`
	// 【不推荐】加密字段，不指定数据类型，默认使用 varchar(128)
	phone AESEncryptedString
	...
}
```

重要：目前 DB 加密使用的密钥为 `service.encrtpySecret`（对应环境变量 `ENCRYPT_SECRET`），开发者可以执行以下命令来生成可用的加密密钥：

```shell
python -c "import base64, os; print(base64.b64encode(os.urandom(32)).decode('utf-8'))"
```

若开发者有使用其他加密算法的需求（如 SM4），可以参考 `pkg/utils/crypto/crypto.go` 以及 `pkg/model/types.go` 中的代码，基于 [crypto-golang-sdk](https://github.com/TencentBlueKing/crypto-golang-sdk) 自行实现相应的 GORM 自定义数据类型（如：`SM4EncryptedString`）。

注：当应用部署在蓝鲸开发者中心后，可以通过读取环境变量 `BKPAAS_BK_CRYPTO_TYPE` 获取推荐的加密方式；具体对应关系：`SHANGMI -> SM4CTR`，`CLASSIC -> Fernet`。

### Swagger

Swagger 是一种 API 协议描述的规范，被广泛用于描述 API 接口的定义；在前后端联调时，swagger 文档可以帮助前端同事更好地了解 API 接口的定义，减轻后端同事编写文档 & 沟通的成本。

除此之外，目前蓝鲸 API 网关还支持通过 swagger 文档来 [管理](https://bk.tencent.com/docs/markdown/ZH/APIGateway/1.10/UserGuide/apigateway/reference/swagger.md) 你的应用网关，支持网关资源的导入、导出等操作。

目前开发框架使用 [swag](https://github.com/swaggo/swag) 来支持从代码注释自动生成 Swagger 文档（`docs/swagger.json`），参考示例如下：

```go
// CreateEntry ...
//
// @Summary    创建条目
// @Tags       crud
// @Param      body    body        serializer.EntryCreateRequest  true  "创建条目请求体"
// @Success    201     {object}    ginx.Response{data=nil}
// @Router     /api/entries [post]
func CreateEntry(c *gin.Context) {...}
```

开发框架在 Makefile 中提供 `make doc` 命令来支持一键生成 `swagger.json`，开发者可以根据需要在终端中执行。

目前 `.pre-commit-config.yaml` 中已默认添加 `make doc` 步骤以在代码提交时自动更新 swagger 文档，若开发者没有该需求则可自行移除。

若开发者配置 `service.enableSwagger` 或环境变量 `ENABLE_SWAGGER` 值为 `true`，则服务会提供 swagger-ui web 服务。

开发者可通过浏览器访问 `http://{host}:{port}/swagger-ui/index.html` 以查看 swagger 文档。

更多参考：

- [Swaggo 使用详解](https://blog.csdn.net/qq_41630102/article/details/128411210)

### 单元测试

我们推荐开发者尽可能多地为代码写单元测试，良好的单元测试能帮助开发者快速发现并定位代码中的错误，保证线上服务的质量；在重构代码时，充足的单元测试是重要的基础，能够显著提升代码质量和开发效率。

开发框架目前提供的单元测试是基于标准库 `testing` 和 [stretchr/testify/assert](https://github.com/stretchr/testify) 实现的，属于比较轻量级的单元测试。开发者可以查阅相关文档，参考现有示例完成单元测试的编写。

如果你想使用 BDD 风格的框架来编写单元测试，可以了解一下 [Ginkgo](https://github.com/onsi/ginkgo) & [Gomega](https://github.com/onsi/gomega) 这对搭档。相比于简单的 assert 测试，它们提供了更强的测试编排、断言能力和可维护性，缺点则是上手需要一些学习成本。

更多参考：

- [有关单元测试的 5 个建议](https://www.piglei.com/articles/5-tips-on-unit-testing/)
- [Golang 单元测试详尽指引](https://cloud.tencent.com/developer/article/1729564)

### 基于 Dockerfile 构建

目前蓝鲸应用默认使用 buildpacks 构建应用镜像，buildpacks 通过预定义的一批脚本，对应用的依赖项进行探测（如 go.mod）并安装依赖，而后执行编译，缓存，导出等行为。

如果你希望在镜像构建环节获得更高的自由度，可以考虑在创建应用 / 模块时，选择构建方式为 `Dockerfile`，并参考项目根目录下的 `Dockerfile` 文件进行修改（推荐本地构建镜像时使用 `make docker-build` 命令）。

### 国际化（i18n）

当蓝鲸应用需动态支持多语言时，Golang 开发框架内置的国际化能力可快速实现需求。

#### 语言配置规范

按蓝鲸统一规范，用户语言配置需存储于名为 `blueking_language` 的 Cookie 中，因此框架仅支持从 Cookies 获取用户语言配置，如需从 Header / QueryParams 获取则需开发者自行实现。

Golang 开发框架通过 `I18n` 中间件将语言信息注入到两种 Context 中：

- `gin.Context`：
    - 在 Gin 处理函数中调用 `ginx.GetLang(c)` 获取
- `context.Context`（需源于 `c.Request.Context()`，其中 `c` 即为 `gin.Context`）：
    - 翻译文本：`i18n.T(ctx, "文本键")`
    - 获取语言配置：`i18n.GetLangFromContext(ctx)`

#### Go 源码中的国际化

```go
import (
	"github.com/pkg/errors"

	"bk.tencent.com/blueapps-go/pkg/i18n"
)

func main() {
	ctx := c.Request.Context()
	fmt.Println(i18n.T(ctx, "这是一个提示信息"))
	_ = errors.New(i18n.T(ctx, "Name required"))
	_ = errors.Errorf(i18n.T(ctx, "名称 %s 不合法"), Name)
}
```

#### Go 模板中的国际化

在开发框架示例中，使用到 Go 模板来渲染出前端的 html 文件，因此也同样支持在 Go 模板中使用国际化。

开发框架通过以下改造来支持国际化：

1. 在 `pkg/router/router.go` 中的 `funcMap` 方法中为 Go 模板添加自定义的国际化函数 `i18n`。
2. 在 `pkg/web/handler/handler.go` 中的 `renderHTML` 方法中默认注入语言信息 `lang` 作为上下文。

在 Go 模板中，你可以通过下面这种写法来实现国际化：

```html
<title>{{ i18n "这是一个标题" .lang }}</title>
```

更多示例可以参考 `templates/web` 目录下的模板。

#### 翻译文件更新

开发者可通过执行 `make i18n` 命令以从源代码 & 模板中收集国际化文本，并默认存储到 `i18n/messages.yaml` 文件中。

接下来开发者需要打开该文件，搜索 `<TODO>` 占位符，并手动翻译成对应的语言版本，新翻译的文本将会在服务重启后生效。

## 更多帮助

新版 Golang Gin 开发框架相比于老版本有很大的变化，如果 SaaS 开发者在使用过程中遇到问题，或者有什么优化的建议，请及时联系蓝鲸助手或平台管理员进行反馈～
