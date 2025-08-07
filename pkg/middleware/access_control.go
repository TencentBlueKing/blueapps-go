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

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// AccessControl 用户访问控制（重要：应该在 UserAuth 中间件后使用）
func AccessControl(allowedUsers []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 没有配置 -> 允许所有用户访问
		if len(allowedUsers) == 0 {
			c.Next()
			return
		}

		// 检查用户是否可访问
		userID := ginx.GetUserID(c)
		if lo.Contains(allowedUsers, userID) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
