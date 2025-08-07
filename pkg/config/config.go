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

// Package config 管理蓝鲸 SaaS 配置项，支持从配置文件 / 环境变量中读取配置
package config

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/samber/lo"

	log "bk.tencent.com/blueapps-go/pkg/logging"
)

var G *Config

// Load 加载配置
func Load(ctx context.Context, cfgFile string) (*Config, error) {
	var cfg *Config
	var err error

	if cfgFile != "" {
		// 若已经指定配置文件，则从配置文件中加载
		log.Infof(ctx, "load config from file: %s", cfgFile)
		cfg, err = loadConfigFromFile(cfgFile)
	} else {
		// 若没有指定配置文件，则环境变量构建配置
		log.Info(ctx, "config file not specified, load config from env vars")
		cfg, err = loadConfigFromEnv()
	}

	if err != nil {
		cfgFrom := lo.Ternary(cfgFile != "", "file: "+cfgFile, "env vars")
		return nil, errors.Wrapf(err, "load config from "+cfgFrom)
	}

	// 后置校验
	// 1. AppSecret 为必填项
	if cfg.Platform.AppSecret == "" {
		return nil, errors.New("config item platform.appSecret is required")
	}
	// 2. 对 EncryptSecret 进行 base64 解码
	if cfg.Service.EncryptSecret != "" {
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(cfg.Service.EncryptSecret)
		if err != nil {
			log.Fatalf("EncryptSecret not valid base64: %s", err)
		}
		cfg.Service.EncryptSecret = string(decoded)
	}

	// 设置全局环境变量
	G = cfg
	return cfg, nil
}
