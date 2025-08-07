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

// Package envx 提供环境变量相关工具
package envx

import (
	"fmt"
	"os"
)

// Get 读取环境变量，支持默认值
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// MustGet 读取环境变量，若不存在则 panic
func MustGet(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("required environment variable %s unset", key))
}

// GetBool 读取环境变量，并自动转换为 bool 类型
// true: 1, t, T, TRUE, true, True
// false: 0, f, F, FALSE, false, False
// 其他情况返回 false
func GetBool(key string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false
	}

	switch value {
	case "1", "t", "T", "TRUE", "true", "True":
		return true
	default:
		return false
	}
}
