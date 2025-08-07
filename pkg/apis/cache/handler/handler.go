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

// Package handler ...
package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/apis/cache/serializer"
	"bk.tencent.com/blueapps-go/pkg/cache/memory"
	"bk.tencent.com/blueapps-go/pkg/cache/redis"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// CacheQuery ...
//
//	@Summary	缓存示例
//	@Tags		cache
//	@Param		message query		string	true	"待哈希数据"
//	@Param		backend	query		string	true	"缓存类型，可选值：redis、memory"
//	@Param		expire	query		int		false	"缓存过期时间（单位：秒）"
//	@Success	200		{object}	ginx.Response{data=serializer.CacheResponse}
//	@Router		/api/cache [get]
func CacheQuery(c *gin.Context) {
	var req serializer.CacheRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	startAt := time.Now()
	var digest string
	var hitCache bool

	ctx := c.Request.Context()
	if req.Backend == "redis" {
		digest, hitCache = redisCacheQuery(ctx, req.Message, req.TTL)
	} else {
		digest, hitCache = memoryCacheQuery(ctx, req.Message, req.TTL)
	}

	respData := serializer.CacheResponse{
		Digest:   digest,
		HitCache: hitCache,
		TimeCost: time.Since(startAt).Seconds(),
	}
	ginx.SetResp(c, http.StatusOK, respData)
}

// 慢哈希
func slowHash(data string) string {
	// 模拟长耗时
	time.Sleep(5 * time.Second)

	hash := md5.New()
	hash.Write([]byte(data))
	sum := hash.Sum(nil)

	return hex.EncodeToString(sum)
}

// 尝试从内存缓存中获取数据，未命中则更新
func memoryCacheQuery(ctx context.Context, message string, ttl int) (digest string, hitCache bool) {
	cache := memory.Cache()

	// 命中缓存
	if value, err := cache.Get([]byte(message)); err == nil {
		return string(value), true
	}

	// 未命中缓存
	digest = slowHash(message)
	if err := cache.Set([]byte(message), []byte(digest), ttl); err != nil {
		log.Errorf(ctx, "set cache `%s: %s`: %s", message, digest, err.Error())
	}
	return digest, false
}

// 尝试从 redis 中获取数据，未命中则更新
func redisCacheQuery(ctx context.Context, message string, ttl int) (digest string, hitCache bool) {
	cache := redis.New("SlowHash", time.Duration(ttl)*time.Second)

	// 命中缓存
	if err := cache.Get(ctx, message, &digest); err != nil {
		log.Errorf(ctx, "get cache `%s`: %s", message, err.Error())
	} else {
		return digest, true
	}

	// 未命中缓存
	digest = slowHash(message)
	if err := cache.Set(ctx, message, digest, 0); err != nil {
		log.Errorf(ctx, "set cache `%s: %s`: %s", message, digest, err.Error())
	}
	return digest, false
}
