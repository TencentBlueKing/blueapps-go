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

package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCSRF(t *testing.T) {
	t.Parallel()

	r := gin.Default()
	appID := "demo"

	r.Use(CSRF("app-test", "fake-secret"))
	r.Use(CSRFToken(appID, ""))

	r.GET("/test_csrf_token", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/test_csrf", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Step 1: 发起一个 GET 请求来获取 CSRF token
	req, _ := http.NewRequest(http.MethodGet, "/test_csrf_token", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var csrfToken string
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == appID+"-csrf-token" {
			// Token 被 url 编码了，需要解码
			csrfToken, _ = url.QueryUnescape(cookie.Value)
		}
	}
	assert.NotEmpty(t, csrfToken)

	// Step 2: 利用获取到的 CSRF token 发起一个 POST 请求
	postReq, _ := http.NewRequest(http.MethodPost, "/test_csrf", strings.NewReader(""))
	postReq.Header.Set("X-CSRF-Token", csrfToken)
	postReq.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	postW := httptest.NewRecorder()
	r.ServeHTTP(postW, postReq)
	assert.Equal(t, http.StatusOK, postW.Code)

	// Step 3: 使用一个错误的 CSRF token 发起 POST 请求
	invalidPostReq, _ := http.NewRequest(http.MethodPost, "/test_csrf", strings.NewReader(""))
	invalidPostReq.Header.Set("X-CSRF-Token", "invalid-token")
	postReq.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	invalidPostW := httptest.NewRecorder()
	r.ServeHTTP(invalidPostW, invalidPostReq)
	assert.Equal(t, http.StatusForbidden, invalidPostW.Code)
}
