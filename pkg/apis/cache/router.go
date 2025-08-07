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

// Package cache 实现缓存实例 (redis/内存缓存)
package cache

import (
	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/apis/cache/handler"
)

// Register ...
func Register(rg *gin.RouterGroup) {
	cacheRouter := rg.Group("/cache")
	cacheRouter.GET("", handler.CacheQuery)
}
