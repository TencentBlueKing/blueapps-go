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

package logging

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/trace"

	"bk.tencent.com/blueapps-go/pkg/common"
)

// Debugf 打印 debug 日志
func Debugf(ctx context.Context, format string, vars ...any) {
	Logf(ctx, slog.LevelDebug, format, vars...)
}

// Debug 打印 debug 日志
func Debug(ctx context.Context, msg string) {
	Logf(ctx, slog.LevelDebug, msg)
}

// Infof 打印 info 日志
func Infof(ctx context.Context, format string, vars ...any) {
	Logf(ctx, slog.LevelInfo, format, vars...)
}

// Info 打印 info 日志
func Info(ctx context.Context, msg string) {
	Logf(ctx, slog.LevelInfo, msg)
}

// Warnf 打印 warn 日志
func Warnf(ctx context.Context, format string, vars ...any) {
	Logf(ctx, slog.LevelWarn, format, vars...)
}

// Warn 打印 warn 日志
func Warn(ctx context.Context, msg string) {
	Logf(ctx, slog.LevelWarn, msg)
}

// Errorf 打印 error 日志
func Errorf(ctx context.Context, format string, vars ...any) {
	Logf(ctx, slog.LevelError, format, vars...)
}

// Error 打印 error 日志
func Error(ctx context.Context, msg string) {
	Logf(ctx, slog.LevelError, msg)
}

// Fatalf 打印 fatal 日志到标准输出并退出程序
// Q：为什么 Fatalf 是强制使用 stderr 而非 slog.Default() ？
// A：调用 Fatalf 意味着程序即将退出，此时往标准输出而不是文件打日志是更合理的（避免 Pod 崩溃导致日志无法采集）
func Fatalf(format string, vars ...any) {
	// 由于马上会退出，这里直接 New logger 而不是预先初始化也是可以的
	logger := log.New(os.Stderr, "", log.LstdFlags)
	logger.Fatalf(format, vars...)
}

// Fatal 打印 fatal 日志到标准输出并退出程序
func Fatal(msg string) {
	Fatalf(msg)
}

// Logf 打印日志
// ref: https://github.com/golang/go/blob/fc9f02c7aec81bcfcc95434d2529e0bb0bc03d66/src/log/slog/example_wrap_test.go#L19
// 注：该方法只能在 logging 包及其子包（如 logging/slogresty）中使用，不得在业务逻辑中直接使用
func Logf(ctx context.Context, level slog.Level, format string, vars ...any) {
	logger := slog.Default()
	if !logger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	r := slog.NewRecord(time.Now(), level, fmt.Sprintf(format, vars...), pcs[0])
	// 尝试获取 Request ID，若存在则需要记录到日志中
	if requestID, ok := ctx.Value(common.RequestIDCtxKey).(string); ok {
		r.AddAttrs(slog.String(common.RequestIDLogKey, requestID))
	}
	if traceID, ok := ExtractTraceID(ctx); ok {
		r.AddAttrs(slog.String(common.TraceIDLogKey, traceID.String()))
	}
	if spanID, ok := ExtractSpanID(ctx); ok {
		r.AddAttrs(slog.String(common.SpanIDLogKey, spanID.String()))
	}

	_ = logger.Handler().Handle(ctx, r)
}

// ExtractTraceID 从 context 中提取 traceID，若不存在则返回空
// ref: https://github.com/samber/slog-gin/blob/4d5fc6c3f623d0fc9b0a7860894c025b1a33fe8c/middleware.go#L289
func ExtractTraceID(ctx context.Context) (slog.Value, bool) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.HasTraceID() {
		traceID := spanCtx.TraceID().String()
		return slog.StringValue(traceID), true
	}

	return slog.Value{}, false
}

// ExtractSpanID 从 context 中提取 Span ID，若不存在则返回空
// ref: https://github.com/samber/slog-gin/blob/4d5fc6c3f623d0fc9b0a7860894c025b1a33fe8c/middleware.go#L289
func ExtractSpanID(ctx context.Context) (slog.Value, bool) {
	spanCtx := trace.SpanFromContext(ctx).SpanContext()
	if spanCtx.HasSpanID() {
		spanID := spanCtx.SpanID().String()
		return slog.StringValue(spanID), true
	}

	return slog.Value{}, false
}
