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

// Package memory 提供内存缓存服务（基于 freecache 封装，内存预分配 + LRU 算法）
// ref: https://github.com/coocood/freecache
package memory

import (
	"sync"

	"github.com/coocood/freecache"

	log "bk.tencent.com/blueapps-go/pkg/logging"
)

var (
	cache    *freecache.Cache
	initOnce sync.Once
)

// InitCache 根据指定容量初始化内存缓存（单位：MB）
func InitCache(capacity int) {
	initOnce.Do(func() {
		cache = freecache.NewCache(capacity * 1024 * 1024)
	})
}

// Cache 获取 cache 实例（提供 Get，Set，Del 等方法）
func Cache() *freecache.Cache {
	if cache == nil {
		log.Fatal("cache not init")
	}
	return cache
}
