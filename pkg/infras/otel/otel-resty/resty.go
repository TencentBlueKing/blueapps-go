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

// Package otelresty provides OpenTelemetry middleware for resty
package otelresty

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("resty-client")

// RequestMiddleware 用于在 resty 发起请求前，记录请求相关 tracing 信息
func RequestMiddleware(_ *resty.Client, req *resty.Request) error {
	ctx, span := tracer.Start(req.Context(), fmt.Sprintf("HTTP %s", req.Method),
		trace.WithAttributes(
			attribute.String("http.url", req.URL),
			attribute.String("http.method", req.Method),
		),
	)
	ctx = context.WithValue(ctx, "otel-span", span)
	req.SetContext(ctx)

	return nil
}

// ResponseMiddleware 用于在 resty 发起请求后，记录响应相关 tracing 信息
func ResponseMiddleware(_ *resty.Client, resp *resty.Response) error {
	span := resp.Request.Context().Value("otel-span").(trace.Span)
	defer span.End()

	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode()),
	)
	if resp.IsError() {
		span.RecordError(errors.Errorf("HTTP error: %s", resp.String()))
	}
	return nil
}
