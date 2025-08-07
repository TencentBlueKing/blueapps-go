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

package common

const (
	// RequestIDHeaderKey Request ID 在 HTTP Header 中的 key
	RequestIDHeaderKey = "X-Request-Id"
)

const (
	// RequestIDCtxKey Request ID 在 context 中的 key
	RequestIDCtxKey = "requestID"

	// UserIDCtxKey user id 在 context 中的 key
	UserIDCtxKey = "userID"

	// UserLangCtxKey user language 在 context 中的 key
	UserLangCtxKey = "userLang"

	// ErrorCtxKey error 在 context 中的 key
	ErrorCtxKey = "error"

	// TracerCtxKey tracer 在 gin.context 中的 key
	// copy: go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin@v0.57.0/gintrace.go:23
	TracerCtxKey = "otel-go-contrib-tracer"
)

const (
	// RequestIDLogKey Request ID 在日志中的 key
	RequestIDLogKey = RequestIDCtxKey

	// TraceIDLogKey trace id 在日志中的 key
	TraceIDLogKey = "otelTraceID"

	// SpanIDLogKey span id 在日志中的 key
	SpanIDLogKey = "otelSpanID"
)

const (
	// UserIDKey user id 在 context / session 中的 key
	UserIDKey = "user_id"

	// UserTokenKey user token 在 cookies / session 中的 key
	UserTokenKey = "user_token"

	// UserAuthSourceKey user 认证来源（AuthBackend Name）在 session 中的 key
	UserAuthSourceKey = "user_auth_source"

	// UserLanguageKey user language 在 cookies / session 中的 key
	UserLanguageKey = "blueking_language"
)
