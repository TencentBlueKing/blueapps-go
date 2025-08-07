# ChangeLog

## 1.6.1 (2025-08-05)

- 优化：`app_desc.yaml` 文件中优先使用 `procCommand` 字段
- 优化：API 文档目录变更为 `pkg/docs`（原位置：`docs`）

## 1.6.0 (2025-07-21)

- 新增：支持多种认证方式（AuthBackend）依次认证
- 新增：支持通过太湖（Taihu）进行用户身份认证
- 新增：添加 `Dockerfile-Debug` 用于构建调试用镜像
- 优化：基于 AST 提取 Go 源码中的国际化消息

## 1.5.0 (2025-06-05)

- 新增：增强服务（MySQL，Redis，RabbitMQ）TLS 配置
- 新增：支持国际化（Golang 源代码 & 模板）
- 优化：MySQL dsn 添加 loc & time_zone 参数

## 1.4.3 (2025-04-11)

- 修复：API 文档按 swagger 规范，在路径中使用 `{id}` 而非 `:id`

## 1.4.2 (2025-03-20)

- 修复：不再对 MySQL、RabbitMQ 密码默认执行 `url.QueryEscape()`

## 1.4.1 (2025-01-17)

- 优化：GORM 日志配置 & 文档
- 优化：log 包新增 `Debug`, `Info`, `Warn`, `Error`, `Fatal` 方法
- 调整：对象存储（bkrepo）客户端不再提前初始化（使用时需新建）
- 调整：增强服务 蓝鲸 APM 配置项名称 `BkAPM -> BkOtel`

## 1.4.0 (2024-11-28)

- 新增：支持接入蓝鲸 APM（Otel / OpenTelemetry）
- 新增：`view-config` 命令以查看加载的配置信息
- 优化：细化 Healthz API 返回组件指标信息

## 1.3.0 (2024-11-20)

- 新增：日志中默认记录 Request ID（日志平台可用）
- 新增：支持数据库字段加密（AES-GCM 对称加密算法）
- 优化：Migration 执行 / 回滚时打印的日志
- 优化：不再指定 web 进程副本数量，推荐页面扩缩容

## 1.2.0 (2024-11-12)

- 新增：Redis & 内存缓存接入 & 使用示例
- 优化：Metric API 上报的指标（支持蓝鲸监控默认看板）
- 优化：CSRF Token 防护相关中间件配置
- 修复：scheduler 在重启前不会注销已删除的定时任务

## 1.1.0 (2024-10-25)

- 新增：默认启用 CSRF 中间件
- 新增：数据库版本控制（基于 gormigrate）

## 1.0.0 (2024-09-03)

- 新增：基础开发框架，应用配置，日志配置等
- 新增：接入统一登录，支持认证以获取用户身份
- 新增：数据库操作，云 API，异步任务，对象存储示例页面
