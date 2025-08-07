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

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/middleware"
	testingutil "bk.tencent.com/blueapps-go/pkg/utils/testing"
)

func TestQueryTokenAuthRight(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c := testingutil.CreateTestContextWithDefaultRequest(w)

	q := c.Request.URL.Query()
	q.Add("token", "token_for_test")
	c.Request.URL.RawQuery = q.Encode()

	middleware.QueryTokenAuth("token_for_test")(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQueryTokenAuthBad(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c := testingutil.CreateTestContextWithDefaultRequest(w)

	middleware.QueryTokenAuth("token_for_test")(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
