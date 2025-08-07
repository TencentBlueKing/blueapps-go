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

package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/cache/memory"
	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/infras/redis"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

func initLogger(cfg *config.LogConfig) error {
	// 自动创建日志目录
	if err := os.MkdirAll(cfg.Dir, os.ModePerm); err != nil {
		// 只有当错误不是 “目录已存在” 时，需要抛出错误
		if !os.IsExist(err) {
			return errors.Wrapf(err, "creating log dir %s", cfg.Dir)
		}
	}

	// 输出位置
	writerName := "file"
	if cfg.ForceToStdout {
		writerName = "stdout"
	}

	// 初始化默认 Logger
	loggerName := "default"
	if err := log.InitLogger(loggerName, &log.Options{
		Level: cfg.Level,
		// 输出到标准输出时 Text 会更加友好
		HandlerName:  lo.Ternary(writerName == "stdout", "text", "json"),
		WriterName:   writerName,
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	// 初始化 Gorm Logger
	loggerName = "gorm"
	if err := log.InitLogger(loggerName, &log.Options{
		Level:        log.GormLogLevel,
		HandlerName:  "json",
		WriterName:   "file",
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	// 初始化 Gin Logger
	loggerName = "gin"
	if err := log.InitLogger(loggerName, &log.Options{
		Level:        log.GinLogLevel,
		HandlerName:  "json",
		WriterName:   "file",
		WriterConfig: map[string]string{"filename": filepath.Join(cfg.Dir, loggerName+".log")},
	}); err != nil {
		return errors.Wrapf(err, "creating logger %s", loggerName)
	}

	return nil
}

// 根据增强服务配置，初始化各类客户端
func initAddons(ctx context.Context, cfg *config.Config) error {
	// 初始化 DB Client
	database.InitDBClient(ctx, cfg.Platform.Addons.Mysql, log.GetLogger("gorm"))

	// 初始化 Redis Client
	if cfg.Platform.Addons.Redis != nil {
		redis.InitRedisClient(ctx, cfg.Platform.Addons.Redis)
	}

	// 初始化缓存
	memory.InitCache(cfg.Service.MemoryCacheSize)

	return nil
}
