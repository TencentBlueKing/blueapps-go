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

package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/middleware"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

func TestI18nMiddlewareDefault(t *testing.T) {
	t.Parallel()

	// request with no language cookie
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := gin.Default()
	r.Use(middleware.I18n())
	r.GET("/ping", func(c *gin.Context) {
		lang := ginx.GetLang(c)
		assert.Equal(t, lang, i18n.LangDefault)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}

func TestI18nMiddlewareWithCookies(t *testing.T) {
	t.Parallel()

	// request with language cookie
	req, _ := http.NewRequest("GET", "/ping", nil)
	req.AddCookie(&http.Cookie{Name: common.UserLanguageKey, Value: string(i18n.LangEN)})

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/ping2", func(c *gin.Context) {
		lang := ginx.GetLang(c)
		assert.Equal(t, lang, i18n.LangEN)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}
