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

// Package i18n provide i18n (Internationalization) support
package i18n

import (
	"context"
	"net/http"
	"strings"

	"bk.tencent.com/blueapps-go/pkg/common"
)

// GetLangFromCookie 从 Cookie 中获取语言版本
func GetLangFromCookie(ck *http.Cookie) Lang {
	if ck == nil {
		return LangDefault
	}

	// 忽略大小写匹配，检查是否为受支持的语言
	if lang, ok := langMap[strings.ToLower(ck.Value)]; ok {
		return lang
	}
	return LangDefault
}

// GetLangCookieValue 获取语言版本 Cookie 值
func GetLangCookieValue(lang Lang) string {
	switch lang {
	case LangZH:
		return CookieValueZH
	case LangEN:
		return CookieValueEN
	}
	return CookieValueEN
}

// T 获取国际化翻译文本
func T(ctx context.Context, msgID string) string {
	return TranslateWithLang(msgID, GetLangFromContext(ctx))
}

// GetLangFromContext 从 Context 中获取语言版本
func GetLangFromContext(ctx context.Context) Lang {
	if lang, ok := ctx.Value(common.UserLangCtxKey).(Lang); ok {
		return lang
	}
	return LangDefault
}

// TranslateWithLang 获取国际化文本
func TranslateWithLang(msgID string, lang Lang) string {
	// 默认情况下直接返回 ID
	if lang == LangDefault {
		return msgID
	}
	if m, ok := i18nMsgMap[msgID]; ok {
		// 特殊忽略占位符的情况
		if ms, ok := m[lang]; ok && ms != Placeholder {
			return ms
		}
	}
	// 找不到时返回传入的 ID 值
	return msgID
}
