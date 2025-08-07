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
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace/noop"

	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

func TestGetRequestID(t *testing.T) {
	c := &gin.Context{}

	requestID := ginx.GetRequestID(c)
	assert.Equal(t, "", requestID)
}

func TestSetRequestID(t *testing.T) {
	c := &gin.Context{}

	ginx.SetRequestID(c, "test")
	assert.Equal(t, "test", ginx.GetRequestID(c))
}

func TestGetError(t *testing.T) {
	c := &gin.Context{}

	err, ok := ginx.GetError(c)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)
}

func TestSetError(t *testing.T) {
	c := &gin.Context{}
	err := errors.New("test")

	ginx.SetError(c, err)
	gErr, ok := ginx.GetError(c)

	assert.Equal(t, true, ok)
	assert.Equal(t, err, gErr)
}

func TestGetUserID(t *testing.T) {
	c := &gin.Context{}

	userID := ginx.GetUserID(c)
	assert.Equal(t, "", userID)
}

func TestSetUserID(t *testing.T) {
	c := &gin.Context{}

	ginx.SetUserID(c, "test")
	assert.Equal(t, "test", ginx.GetUserID(c))
}

func TestGetLang(t *testing.T) {
	c := &gin.Context{}

	lang := ginx.GetLang(c)
	assert.Equal(t, i18n.LangDefault, lang)
}

func TestSetLang(t *testing.T) {
	c := &gin.Context{}

	ginx.SetLang(c, i18n.LangEN)
	assert.Equal(t, i18n.LangEN, ginx.GetLang(c))
}

func TestSetAndGetTracer(t *testing.T) {
	c := &gin.Context{}

	ginx.SetTracer(c, noop.Tracer{})
	assert.NotNil(t, ginx.GetTracer(c))
}
