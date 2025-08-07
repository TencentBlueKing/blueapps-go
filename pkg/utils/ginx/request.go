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

package ginx

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

const (
	// NOTE: SaaS 开发者可根据需求自定义 Limit 上下限
	// Limit 下限
	minLimit = 5
	// Limit 上限
	maxLimit = 100
)

// GetPage 获取分页参数 Page
func GetPage(c *gin.Context) int {
	page := cast.ToInt(c.Query("page"))
	// page 必须为非负整数
	return max(1, page)
}

// GetLimit 获取分页参数 Limit
func GetLimit(c *gin.Context) int {
	limit := cast.ToInt(c.Query("limit"))
	limit = min(maxLimit, limit)
	limit = max(minLimit, limit)
	return limit
}

// GetOffset 获取分页参数 Offset
func GetOffset(c *gin.Context) int {
	page, limit := GetPage(c), GetLimit(c)
	return (page - 1) * limit
}
