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

// Package ginx 提供一些 Gin 框架相关的工具
package ginx

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/i18n"
)

// GetRequestID ...
func GetRequestID(c *gin.Context) string {
	return c.GetString(common.RequestIDCtxKey)
}

// SetRequestID ...
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(common.RequestIDCtxKey, requestID)
}

// GetError ...
func GetError(c *gin.Context) (err any, ok bool) {
	return c.Get(common.ErrorCtxKey)
}

// SetError ...
func SetError(c *gin.Context, err error) {
	c.Set(common.ErrorCtxKey, err)
}

// GetUserID ...
func GetUserID(c *gin.Context) string {
	return c.GetString(common.UserIDCtxKey)
}

// SetUserID ...
func SetUserID(c *gin.Context, userID string) {
	c.Set(common.UserIDCtxKey, userID)
}

// GetLang ...
func GetLang(c *gin.Context) i18n.Lang {
	lang, ok := c.Get(common.UserLanguageKey)
	if !ok {
		return i18n.LangDefault
	}
	return lang.(i18n.Lang)
}

// SetLang ...
func SetLang(c *gin.Context, lang i18n.Lang) {
	c.Set(common.UserLanguageKey, lang)
}

// GetTracer 获取 tracer（the creator of Spans）
func GetTracer(c *gin.Context) trace.Tracer {
	tracer, ok := c.Get(common.TracerCtxKey)
	if !ok {
		return nil
	}
	return tracer.(trace.Tracer)
}

// SetTracer ...
func SetTracer(c *gin.Context, tracer trace.Tracer) {
	c.Set(common.TracerCtxKey, tracer)
}
