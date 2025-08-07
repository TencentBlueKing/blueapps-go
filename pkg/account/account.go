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

// Package account 提供不同版本的用户认证后端
package account

// GetAuthBackends 获取需要的 AuthBackend
func GetAuthBackends(authTypes []string) []AuthBackend {
	if len(authTypes) == 0 {
		panic("at least one auth type is required")
	}

	funcMap := map[string]func() AuthBackend{
		bkTicketAuthBackendName: NewBkTicketAuthBackend,
		bkTokenAuthBackendName:  NewBkTokenAuthBackend,
		taihuAuthBackendName:    NewTaihuAuthBackend,
	}

	var backends []AuthBackend
	for _, authType := range authTypes {
		fn, ok := funcMap[authType]
		if !ok {
			panic("unsupported auth type: " + authType)
		}
		backends = append(backends, fn())
	}
	return backends
}
