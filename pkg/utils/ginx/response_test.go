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

package ginx_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

func TestSetSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/success", func(c *gin.Context) {
		ginx.SetResp(c, http.StatusOK, "test data")
	})
	req, _ := http.NewRequest(http.MethodGet, "/success", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedResponse := ginx.Response{Message: "", Data: "test data"}

	var actualResponse ginx.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestSetErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/error", func(c *gin.Context) {
		ginx.SetErrResp(c, http.StatusInternalServerError, "error occurred")
	})
	req, _ := http.NewRequest(http.MethodGet, "/error", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)

	expectedResponse := ginx.Response{Message: "error occurred", Data: nil}

	var actualResponse ginx.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestNewPaginatedRespData(t *testing.T) {
	data := ginx.NewPaginatedRespData(100, []string{"alpha", "beta", "gamma"})
	assert.Equal(t, ginx.PaginatedResp{Count: int64(100), Results: []string{"alpha", "beta", "gamma"}}, data)
}
