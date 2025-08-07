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
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// I18n 国际化中间件
func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从 Cookie 中获取 User Language，如果不存在则使用默认值
		lang := i18n.LangDefault
		if cookie, err := c.Request.Cookie(common.UserLanguageKey); err == nil {
			lang = i18n.GetLangFromCookie(cookie)
		}

		// 在 context 中设置 User Language
		ctx := context.WithValue(c.Request.Context(), common.UserLangCtxKey, lang)
		c.Request = c.Request.WithContext(ctx)
		// 在 gin.Context 中设置 User Language
		ginx.SetLang(c, lang)

		c.Next()
	}
}
