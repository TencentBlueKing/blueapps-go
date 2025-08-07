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

package account

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bk.tencent.com/blueapps-go/pkg/common"
)

// StubType 测试用 AuthBackend 类型
type StubType string

const (
	// AllowAllStubType 允许所有用户
	AllowAllStubType StubType = "allow_all"
	// ForbidAllStubType 禁止所有用户
	ForbidAllStubType StubType = "forbid_all"
	// AllowSpecialTokenStubType 允许特定用户
	AllowSpecialTokenStubType StubType = "allow_special_token"
)

// SpecialTokenForTest 允许特定用户凭证（测试用）
const SpecialTokenForTest = "EverythingIsPermitted"

// StubAuthBackend 测试用 AuthBackend
type StubAuthBackend struct {
	stubType StubType
}

// Name 认证后端名称
func (b *StubAuthBackend) Name() string {
	return string(b.stubType)
}

// GetLoginUrl 获取登录地址
func (b *StubAuthBackend) GetLoginUrl(_ string) string {
	return "https://bklogin.example.com/plain/"
}

// GetUserToken 获取用户凭证
func (b *StubAuthBackend) GetUserToken(c *gin.Context) (string, error) {
	token, err := c.Request.Cookie(common.UserTokenKey)
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

// GetUserInfo 获取用户信息
func (b *StubAuthBackend) GetUserInfo(_ *gin.Context, token string) (*UserInfo, error) {
	// 禁止所有用户
	if b.stubType == ForbidAllStubType {
		return nil, errors.New("forbidden")
	}

	userInfo := &UserInfo{ID: "admin", Token: token, AuthSource: b.Name()}
	// 允许所有用户
	if b.stubType == AllowAllStubType {
		userInfo.ID = "anonymous"
		return userInfo, nil
	}
	// 允许特定用户 Token
	if b.stubType == AllowSpecialTokenStubType && token == SpecialTokenForTest {
		return userInfo, nil
	}
	return nil, errors.New("invalid token")
}

var _ AuthBackend = (*StubAuthBackend)(nil)

// NewStubAuthBackend ...
func NewStubAuthBackend(t StubType) AuthBackend {
	return &StubAuthBackend{stubType: t}
}
