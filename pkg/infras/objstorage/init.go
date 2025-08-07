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

// Package objstorage 提供对象存储相关封装，目前接入的是蓝盾制品库（bkrepo）
// 如果 SaaS 开发者需要使用其他云对象存储（如 COS，S3, Ceph 等），可参考相关实现
package objstorage

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"bk.tencent.com/blueapps-go/pkg/config"
	otelresty "bk.tencent.com/blueapps-go/pkg/infras/otel/otel-resty"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	slogresty "bk.tencent.com/blueapps-go/pkg/logging/slog-resty"
)

// NewClient 获取对象存储客户端
func NewClient(ctx context.Context) *BkGenericRepoClient {
	if !IsBkRepoAvailable() {
		log.Error(ctx, "bkrepo is not available")
		return nil
	}

	cli := newBkGenericRepoClient(config.G.Platform.Addons.BkRepo)
	// 根据 ctx 设置 Logger，以支持记录 Request / Span / Trace ID 等信息
	cli.client = cli.client.SetLogger(slogresty.New(ctx))
	return cli
}

// IsBkRepoAvailable 判断蓝盾制品仓库是否可用
func IsBkRepoAvailable() bool {
	return config.G.Platform.Addons.BkRepo != nil
}

// 初始化客户端 ...
func newBkGenericRepoClient(cfg *config.BkRepoConfig) *BkGenericRepoClient {
	// 使用连接池
	transport := &http.Transport{
		MaxIdleConns:        5,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     30 * time.Second,
	}

	client := resty.New().
		SetTransport(transport).
		SetBaseURL(strings.TrimSuffix(cfg.EndpointUrl, "/")).
		SetTimeout(60*time.Second).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetRetryCount(2).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(
			func(response *resty.Response, err error) bool {
				// Retry on 5xx status codes
				return response.StatusCode() >= http.StatusInternalServerError
			},
		).
		// OpenTelemetry 相关中间件
		OnBeforeRequest(otelresty.RequestMiddleware).
		OnAfterResponse(otelresty.ResponseMiddleware)

	return &BkGenericRepoClient{cfg: cfg, client: client}
}
