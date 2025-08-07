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

import "github.com/gin-gonic/gin"

const (
	// 认证后端名称
	bkTicketAuthBackendName = "BkTicket"
	bkTokenAuthBackendName  = "BkToken"
	taihuAuthBackendName    = "Taihu"
)

// UserInfo 用户信息
type UserInfo struct {
	ID         string
	Token      string
	AuthSource string
}

// AuthBackend 认证后端
type AuthBackend interface {
	// Name 获取认证后端名称
	Name() string
	// GetLoginUrl 获取登录地址，callbackUrl 为当前请求的地址，用于登录成功后跳转
	GetLoginUrl(callbackUrl string) string
	// GetUserToken 获取用户凭证
	GetUserToken(c *gin.Context) (string, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(c *gin.Context, token string) (*UserInfo, error)
}
