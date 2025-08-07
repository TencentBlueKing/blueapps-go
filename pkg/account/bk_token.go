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
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"bk.tencent.com/blueapps-go/pkg/config"
	slogresty "bk.tencent.com/blueapps-go/pkg/logging/slog-resty"
)

// BkTokenAuthBackend 用于社区开源版本的用户登录 & 信息获取
type BkTokenAuthBackend struct{}

// Name 认证后端名称
func (b *BkTokenAuthBackend) Name() string {
	return bkTokenAuthBackendName
}

// GetLoginUrl 获取登录地址
func (b *BkTokenAuthBackend) GetLoginUrl(callbackUrl string) string {
	loginUrl := fmt.Sprintf("%s/plain/", config.G.Platform.BkPlatUrl.BkLogin)
	if callbackUrl != "" {
		loginUrl += fmt.Sprintf("?c_url=%s", url.QueryEscape(callbackUrl))
	}
	return loginUrl
}

// GetUserToken 获取用户凭证
func (b *BkTokenAuthBackend) GetUserToken(c *gin.Context) (string, error) {
	token, err := c.Request.Cookie("bk_token")
	if err != nil {
		return "", err
	}
	return token.Value, nil
}

// GetUserInfo 获取用户信息
func (b *BkTokenAuthBackend) GetUserInfo(c *gin.Context, token string) (*UserInfo, error) {
	getUserInfoUrl := fmt.Sprintf("%s/accounts/get_user/", config.G.Platform.BkPlatUrl.BkLogin)

	client := resty.New().SetLogger(slogresty.New(c.Request.Context())).SetTimeout(10 * time.Second)

	respData := map[string]any{}
	_, err := client.R().
		SetQueryParams(map[string]string{"bk_token": token}).
		ForceContentType("application/json").
		SetResult(&respData).
		Get(getUserInfoUrl)
	if err != nil {
		return nil, err
	}

	if retCode, cErr := cast.ToIntE(respData["code"]); cErr != nil {
		return nil, errors.Errorf("get user info api %s return code isn't integer", getUserInfoUrl)
	} else if retCode != 0 {
		return nil, errors.Errorf("failed to get user info from %s, message: %s", getUserInfoUrl, respData["message"])
	}

	data, ok := respData["data"].(map[string]any)
	if !ok {
		return nil, errors.Errorf("failed to get user info from %s, response data not json format", getUserInfoUrl)
	}
	return &UserInfo{ID: data["username"].(string), Token: token, AuthSource: b.Name()}, nil
}

var _ AuthBackend = (*BkTokenAuthBackend)(nil)

// NewBkTokenAuthBackend ...
func NewBkTokenAuthBackend() AuthBackend {
	return &BkTokenAuthBackend{}
}
