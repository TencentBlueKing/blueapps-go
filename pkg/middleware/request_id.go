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
	"context"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
	"bk.tencent.com/blueapps-go/pkg/utils/uuidx"
)

// RequestID 中间件用于注入 RequestID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(common.RequestIDHeaderKey)
		// 若 RequestID 不存在或不是 32 位随机字符串，则需要重新生成并写入到 Request Header
		if len(requestID) != 32 {
			requestID = uuidx.New()
			// Request 的 Header 中需要注入 RequestID，方便 slog 中获取
			c.Request.Header.Set(common.RequestIDHeaderKey, requestID)
		}

		// 在 context 中设置 RequestID
		ctx := context.WithValue(c.Request.Context(), common.RequestIDCtxKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		// 在 gin.Context 中设置 RequestID
		ginx.SetRequestID(c, requestID)
		// Writer 的 Header 中需要注入 RequestID，用于提供给请求方
		c.Writer.Header().Set(common.RequestIDHeaderKey, requestID)

		c.Next()
	}
}
