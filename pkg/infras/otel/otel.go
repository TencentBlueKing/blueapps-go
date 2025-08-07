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

// Package otel 用于初始化 OpenTelemetry
package otel

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"bk.tencent.com/blueapps-go/pkg/config"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

// 采样比例
const ratio = 0.5

// 采样策略映射表
var samplerMap = map[string]sdktrace.Sampler{
	"always_on":                sdktrace.AlwaysSample(),
	"always_off":               sdktrace.NeverSample(),
	"parentbased_always_on":    sdktrace.ParentBased(sdktrace.AlwaysSample()),
	"parentbased_always_off":   sdktrace.ParentBased(sdktrace.NeverSample()),
	"traceidratio":             sdktrace.TraceIDRatioBased(ratio),
	"parentbased_traceidratio": sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio)),
}

// 创建 OpenTelemetry exporter
func newGRPCTracerExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

// 创建 OpenTelemetry 资源
func newResource(cfg *config.BkOtelConfig, serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		attribute.Key("bk.data.token").String(cfg.BkDataToken),
	)
}

// 创建 OpenTelemetry 追踪器
func newTracerProvider(
	res *resource.Resource,
	exporter *otlptrace.Exporter,
	sampler sdktrace.Sampler,
) *sdktrace.TracerProvider {
	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
}

// GenServiceName 给服务名称添加前缀
func GenServiceName(proc string) string {
	cfg := config.G.Platform
	return fmt.Sprintf("%s-%s-%s-%s", cfg.AppID, cfg.ModuleName, cfg.RunEnv, proc)
}

// 创建采样策略
func newSampler(sampler string) sdktrace.Sampler {
	if s, ok := samplerMap[sampler]; ok {
		return s
	}
	return sdktrace.AlwaysSample()
}

// 创建 OpenTelemetry 传播器
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// ShutdownFunc 用于关闭 OpenTelemetry 追踪器
type ShutdownFunc func(context.Context) error

// InitTracer 初始化全局 OpenTelemetry 追踪器
func InitTracer(
	ctx context.Context, cfg *config.BkOtelConfig, serviceName string,
) (ShutdownFunc, error) {
	// 只保留 IP:Port / Domain:Port（即 Host），无需 Scheme 等其他信息
	u, err := url.Parse(cfg.GrpcUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing url %s", cfg.GrpcUrl)
	}

	client, err := grpc.NewClient(u.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "creating grpc client: %s", cfg.GrpcUrl)
	}

	exporter, err := newGRPCTracerExporter(ctx, client)
	if err != nil {
		return nil, errors.Wrapf(err, "creating OTLP trace exporter: %s", cfg.GrpcUrl)
	}

	tracerProvider := newTracerProvider(newResource(cfg, serviceName), exporter, newSampler(cfg.Sampler))
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(newPropagator())
	log.Infof(ctx, "otel tracer provider: %s successfully initialized", serviceName)
	return tracerProvider.Shutdown, nil
}
