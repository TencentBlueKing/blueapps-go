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
	"bk.tencent.com/blueapps-go/pkg/middleware"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

func TestRequestIDMiddlewareWithoutRequestID(t *testing.T) {
	t.Parallel()

	// request with no request_id
	req, _ := http.NewRequest("GET", "/ping", nil)

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/ping", func(c *gin.Context) {
		requestID := ginx.GetRequestID(c)
		assert.NotNil(t, requestID)
		assert.Len(t, requestID, 32)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}

func TestRequestIDMiddlewareWithRequestID(t *testing.T) {
	t.Parallel()

	// request with X-Request-Id
	req, _ := http.NewRequest("GET", "/ping", nil)
	originRID := "ca7ff4ce433447a99e8175f28af31460"
	req.Header.Set(common.RequestIDHeaderKey, originRID)

	r := gin.Default()
	r.Use(middleware.RequestID())
	r.GET("/ping2", func(c *gin.Context) {
		requestID := ginx.GetRequestID(c)
		assert.NotNil(t, requestID)
		assert.Equal(t, originRID, requestID)
		c.String(http.StatusOK, "pong")
	})

	r.ServeHTTP(httptest.NewRecorder(), req)
}
