/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - Go 开发框架 (BlueKing - Go Framework) available.
 * Copyright (C) 2017 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 *	https://opensource.org/licenses/MIT
 *
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package router 是项目 API 服务的主路由入口
package router

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/samber/lo"
	slogGin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"bk.tencent.com/blueapps-go/pkg/account"
	"bk.tencent.com/blueapps-go/pkg/apis/asynctask"
	"bk.tencent.com/blueapps-go/pkg/apis/basic"
	"bk.tencent.com/blueapps-go/pkg/apis/cache"
	"bk.tencent.com/blueapps-go/pkg/apis/cloudapi"
	"bk.tencent.com/blueapps-go/pkg/apis/crud"
	"bk.tencent.com/blueapps-go/pkg/apis/objstorage"
	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/otel"
	"bk.tencent.com/blueapps-go/pkg/middleware"
	"bk.tencent.com/blueapps-go/pkg/web"
)

// New create server router
func New(slogger *slog.Logger) *gin.Engine {
	gin.SetMode(config.G.Service.Server.GinRunMode)
	gin.DisableConsoleColor()

	router := gin.New()
	_ = router.SetTrustedProxies(nil)
	store := cookie.NewStore([]byte(config.G.Platform.AppSecret))
	store.Options(sessions.Options{MaxAge: int(30 * time.Minute)})
	router.Use(sessions.Sessions(fmt.Sprintf("%s-session", config.G.Platform.AppID), store))

	// 服务指标
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	// 探针相关 API 不应被记录
	m.SetExcludePaths([]string{"/ping", "/healthz"})
	// 默认超过 1s 算是慢请求
	m.SetSlowTime(1)
	// 请求时间分段记录
	m.SetDuration([]float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5})
	m.UseWithoutExposingEndpoint(router)

	// Gin 中间件
	setMiddlewares(router, slogger)

	// 设置静态文件
	router.Static("/static", config.G.Service.StaticFileBaseDir)
	// 设置模板方法
	router.SetFuncMap(funcMap())
	// 加载 HTML 模板文件
	router.LoadHTMLGlob(config.G.Service.TmplFileBaseDir + "/web/*")
	// 404 访问路径不存在
	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

	// 基础 API
	basic.Register(router)

	// 用户认证后端
	authBackends := account.GetAuthBackends(config.G.Service.AuthTypes)
	// 前端页面
	{
		webfeRG := router.Group("")
		webfeRG.Use(middleware.UserAuth(authBackends))
		webfeRG.Use(middleware.AccessControl(config.G.Service.AllowedUsers))
		web.Register(webfeRG)
	}

	// 后端 API
	{
		apiRG := router.Group("/api")
		apiRG.Use(middleware.UserAuth(authBackends))
		apiRG.Use(middleware.AccessControl(config.G.Service.AllowedUsers))

		// 数据库 CRUD 示例
		crud.Register(apiRG)
		// 内存 / Redis 缓存示例
		cache.Register(apiRG)
		// 云 API 调用示例
		cloudapi.Register(apiRG)
		// 异步任务调用示例
		asynctask.Register(apiRG)
		// 对象存储调用示例
		objstorage.Register(apiRG)
	}

	return router
}

// 为 gin.Engine 设置中间件
// otelgin：OpenTelemetry - Gin Tracing 上报
// RequestID：在 Context，HTTP Header 中设置 Request ID
// I18n：从 cookies 中读取，并在 Context 中设置国际化语言信息
// slogGin：记录 Gin 框架结构化日志
// CORS / CSRF / CSRFToken：跨域设置 / CSRF 防护
// Recovery：请求 Panic 恢复
func setMiddlewares(router *gin.Engine, slogger *slog.Logger) {
	router.Use(otelgin.Middleware(
		otel.GenServiceName("web"),
		otelgin.WithGinFilter(
			func(c *gin.Context) bool {
				// 忽略部分路径避免过于骚扰
				excludedPaths := []string{"/metrics", "/ping"}
				return !lo.Contains(excludedPaths, c.Request.URL.Path)
			},
		),
	))
	router.Use(middleware.RequestID())
	router.Use(middleware.I18n())
	// 替换 slogGin 配置以保持一致
	slogGin.RequestIDKey = common.RequestIDLogKey
	slogGin.SpanIDKey = common.SpanIDLogKey
	slogGin.TraceIDKey = common.TraceIDLogKey
	cfg := slogGin.Config{WithTraceID: true, WithSpanID: true, WithRequestID: true}
	router.Use(slogGin.NewWithConfig(slogger, cfg))
	router.Use(middleware.CORS(config.G.Service.AllowedOrigins))
	router.Use(middleware.CSRF(config.G.Platform.AppID, config.G.Platform.AppSecret))
	router.Use(middleware.CSRFToken(config.G.Platform.AppID, config.G.Service.CSRFCookieDomain))
	router.Use(gin.Recovery())
}

// 自定义的模板方法
func funcMap() template.FuncMap {
	fm := sprig.FuncMap()
	// 添加国际化方法
	fm["i18n"] = i18n.TranslateWithLang
	return fm
}
