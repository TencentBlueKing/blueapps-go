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

// Package serializer ...
package serializer

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
)

// CacheRequest Cache API 请求结构
type CacheRequest struct {
	Message string `form:"message" binding:"required"`
	Backend string `form:"backend" binding:"required"`
	TTL     int    `form:"ttl" binding:"required"`
}

// CacheResponse Cache API 返回结构
type CacheResponse struct {
	Digest   string  `json:"digest"`
	HitCache bool    `json:"hitCache"`
	TimeCost float64 `json:"timeCost"`
}

// Validate ...
func (r *CacheRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	if !slices.Contains([]string{"memory", "redis"}, r.Backend) {
		return errors.New(i18n.T(ctx, "unsupported cache backend"))
	}
	if r.Backend == "redis" && config.G.Platform.Addons.Redis == nil {
		return errors.New(i18n.T(ctx, "redis cache backend is not enabled"))
	}
	return nil
}
