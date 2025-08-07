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

// Package handler ...
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/apis/basic/serializer"
	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/common/probe"
	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
	"bk.tencent.com/blueapps-go/pkg/version"
)

// Ping ...
//
//	@Summary	服务探活
//	@Tags		basic
//	@Produce	text/plain
//	@Success	200	{string}	string	pong
//	@Router		/ping [get]
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Healthz ...
//
//	@Summary	提供服务健康状态
//	@Tags		basic
//	@Param		token	query		string	true	"healthz api token"
//	@Success	200		{object}	serializer.HealthzResponse
//	@Router		/healthz [get]
func Healthz(c *gin.Context) {
	ctx := c.Request.Context()

	healthy, fatal := true, false
	var results []probe.Result

	probes := []probe.HealthProbe{
		probe.NewGin(),
		probe.NewMysql(),
		probe.NewRedis(),
		probe.NewBkRepo(),
	}
	for _, p := range probes {
		ret := p.Perform(ctx)
		if ret == nil {
			continue
		}

		// 任意探针失败，则为不健康
		healthy = healthy && ret.Healthy
		// 任意核心组件探针失败，则为致命异常
		fatal = fatal || (ret.Core && !ret.Healthy)
		results = append(results, *ret)
	}
	respData := serializer.HealthzResponse{
		Time:    time.Now().Format(time.RFC3339),
		Healthy: healthy,
		Fatal:   fatal,
		Results: results,
	}
	// 如果核心服务不可用，应该返回 503 而非 200
	c.JSON(lo.Ternary(fatal, http.StatusServiceUnavailable, http.StatusOK), respData)
}

// Version ...
//
//	@Summary	服务版本信息
//	@Tags		basic
//	@Success	200	{object}	serializer.VersionResponse
//	@Router		/version [get]
func Version(c *gin.Context) {
	respData := serializer.VersionResponse{
		Version:     version.AppVersion,
		GitCommit:   version.GitCommit,
		BuildTime:   version.BuildTime,
		TmplVersion: version.TmplVersion,
		GoVersion:   version.GoVersion,
	}
	c.JSON(http.StatusOK, respData)
}

// Language ...
//
//	@Summary	修改语言
//	@Tags		basic
//	@Success	204	"No Content"
//	@Router		/language [put]
func Language(c *gin.Context) {
	var req serializer.LanguageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	val := i18n.GetLangCookieValue(req.Lang)
	c.SetCookie(common.UserLanguageKey, val, 0, "/", config.G.Platform.BkDomain, false, false)
	c.String(http.StatusNoContent, "")
}

// Metrics ...
//
//	@Summary	Prometheus 指标
//	@Tags		basic
//	@Produce	text/plain
//	@Param		token	query		string	true	"metrics api token"
//	@Success	200		{string}	string	metrics
//	@Router		/metrics [get]
func Metrics() {} // nolint
