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

package envx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/utils/envx"
)

// 不存在的环境变量
func TestGetNotExists(t *testing.T) {
	ret := envx.Get("NOT_EXISTS_ENV_KEY", "ENV_VAL")
	assert.Equal(t, "ENV_VAL", ret)
}

// 已存在的环境变量
func TestGetExists(t *testing.T) {
	ret := envx.Get("PATH", "")
	assert.NotEqual(t, "", ret)
}

// 不存在的环境变量
func TestMustGetNotExists(t *testing.T) {
	defer func() {
		assert.Equal(t, "required environment variable NOT_EXISTS_ENV_KEY unset", recover())
	}()

	_ = envx.MustGet("NOT_EXISTS_ENV_KEY")
}

// 已存在的环境变量
func TestMustGetExists(t *testing.T) {
	ret := envx.MustGet("PATH")
	assert.NotEqual(t, "", ret)
}

// 测试 GetBool 函数（True 场景）
func TestGetBoolTrueCase(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{"TRUE_1", "1"},
		{"TRUE_2", "t"},
		{"TRUE_3", "T"},
		{"TRUE_4", "true"},
		{"TRUE_5", "TRUE"},
		{"TRUE_6", "True"},
	}
	for _, tt := range tests {
		_ = os.Setenv(tt.key, tt.value)
		ret := envx.GetBool(tt.key)
		assert.True(t, ret)
	}
}

// 测试 GetBool 函数（False 场景）
func TestGetBoolFalseCase(t *testing.T) {
	tests := []struct {
		key   string
		value string
	}{
		{"FALSE_0", ""},
		{"FALSE_1", "0"},
		{"FALSE_2", "f"},
		{"FALSE_3", "F"},
		{"FALSE_4", "false"},
		{"FALSE_5", "FALSE"},
		{"FALSE_6", "False"},
	}
	for _, tt := range tests {
		_ = os.Setenv(tt.key, tt.value)
		ret := envx.GetBool(tt.key)
		assert.False(t, ret)
	}
}

// 测试 GetBool 函数（不存在的环境变量 -> false）
func TestGetBoolNotExistKey(t *testing.T) {
	ret := envx.GetBool("NOT_EXISTS_KEY")
	assert.False(t, ret)
}
