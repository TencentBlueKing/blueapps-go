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

// Package probe provides health probes for components
package probe

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/infras/objstorage"
	"bk.tencent.com/blueapps-go/pkg/infras/redis"
)

// GinProbe Gin 服务探针
type GinProbe struct{}

// NewGin ...
func NewGin() *GinProbe {
	return &GinProbe{}
}

// Perform ...
func (p GinProbe) Perform(_ context.Context) *Result {
	return &Result{Name: "Gin", Core: true, Healthy: true, Endpoint: "/", Issue: ""}
}

// MysqlProbe Mysql 服务探针
type MysqlProbe struct{}

// NewMysql ...
func NewMysql() *MysqlProbe {
	return &MysqlProbe{}
}

// Perform ...
func (p *MysqlProbe) Perform(ctx context.Context) *Result {
	cfg := config.G.Platform.Addons.Mysql
	if cfg == nil {
		return nil
	}

	healthy, issue := true, ""
	if err := database.Client(ctx).Exec("SELECT 1").Error; err != nil {
		healthy, issue = false, err.Error()
	}

	ep := fmt.Sprintf(
		"%s:***@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		cfg.User,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
	)
	return &Result{
		Name:     "Mysql",
		Core:     true,
		Healthy:  healthy,
		Endpoint: lo.Ternary(healthy, "", ep),
		Issue:    issue,
	}
}

var _ HealthProbe = &MysqlProbe{}

// RedisProbe redis 服务探针
type RedisProbe struct{}

// NewRedis ...
func NewRedis() *RedisProbe {
	return &RedisProbe{}
}

// Perform ...
func (p *RedisProbe) Perform(ctx context.Context) *Result {
	cfg := config.G.Platform.Addons.Redis
	if cfg == nil {
		return nil
	}

	healthy, issue := true, ""
	if err := redis.Client().Ping(ctx).Err(); err != nil {
		healthy, issue = false, err.Error()
	}

	ep := fmt.Sprintf("redis://%s:***@%s:%d/%d", cfg.Username, cfg.Host, cfg.Port, cfg.DB)
	return &Result{
		Name:     "Redis",
		Core:     false,
		Healthy:  healthy,
		Endpoint: lo.Ternary(healthy, "", ep),
		Issue:    issue,
	}
}

var _ HealthProbe = &RedisProbe{}

// BkRepoProbe BkRepo 服务探针
type BkRepoProbe struct{}

// NewBkRepo ...
func NewBkRepo() *BkRepoProbe {
	return &BkRepoProbe{}
}

// Perform ...
func (p *BkRepoProbe) Perform(ctx context.Context) *Result {
	cfg := config.G.Platform.Addons.BkRepo
	if cfg == nil {
		return nil
	}

	healthy, issue := true, ""
	if _, err := objstorage.NewClient(ctx).ListDir(ctx, "/", 1, 1); err != nil {
		healthy, issue = false, err.Error()
	}

	ep := fmt.Sprintf(
		"%s (project: %s, username: %s, bucket: %s)",
		cfg.EndpointUrl,
		cfg.Project,
		cfg.Username,
		cfg.Bucket,
	)
	return &Result{
		Name:     "BkRepo",
		Core:     false,
		Healthy:  healthy,
		Endpoint: lo.Ternary(healthy, "", ep),
		Issue:    issue,
	}
}

var _ HealthProbe = &BkRepoProbe{}
