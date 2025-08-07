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
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"bk.tencent.com/blueapps-go/pkg/async"
	"bk.tencent.com/blueapps-go/pkg/i18n"
)

// PeriodicTaskListResponse List PeriodicTask API 返回结构
type PeriodicTaskListResponse struct {
	ID      int64  `json:"id"`
	Cron    string `json:"cron"`
	Name    string `json:"name"`
	Args    string `json:"args"`
	Enabled bool   `json:"enabled"`
	Creator string `json:"creator"`
}

// PeriodicTaskCreateRequest Create PeriodicTask API 请求结构
type PeriodicTaskCreateRequest struct {
	Name string `json:"name"`
	Cron string `json:"cron"`
	Args []any  `json:"args"`
}

// Validate ...
func (r *PeriodicTaskCreateRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	// 检查 name 是否合法
	if r.Name == "" {
		return errors.New(i18n.T(ctx, "Task name required"))
	}
	if _, ok := async.RegisteredTasks[r.Name]; !ok {
		return errors.Errorf(i18n.T(ctx, "Task name %s invalid"), r.Name)
	}
	// 检查 cron 表达式是否合法
	if r.Cron == "" {
		return errors.New(i18n.T(ctx, "cron required"))
	}
	if _, err := cron.ParseStandard(r.Cron); err != nil {
		return errors.Wrap(err, i18n.T(ctx, "cron invalid"))
	}
	return nil
}

// TogglePeriodicTaskEnabledResponse TogglePeriodicTaskEnabled API 返回结构
type TogglePeriodicTaskEnabledResponse struct {
	Enabled bool `json:"enabled"`
}
