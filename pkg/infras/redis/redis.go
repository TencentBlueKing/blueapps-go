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

// Package redis 提供了 Redis 相关的封装（基于 redis/go-redis/v9）
// SaaS 开发者查阅该文档以了解使用方法：https://redis.uptrace.dev/guide/go-redis.html
package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	"bk.tencent.com/blueapps-go/pkg/config"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

var (
	rds      *redis.Client
	initOnce sync.Once
)

const (
	// 尝试连接超时 单位：s
	dialTimeout = 2
	// 读超时 单位：s
	readTimeout = 1
	// 写超时 单位：s
	writeTimeout = 1
	// 闲置超时 单位: s
	idleTimeout = 3 * 60
	// 连接池大小 / 核
	poolSizeMultiple = 20
	// 最小空闲连接数 / 核
	minIdleConnectionMultiple = 10
)

// 生成 redis 连接选项
func buildOpts(cfg *config.RedisConfig) (*redis.Options, error) {
	opts, err := redis.ParseURL(cfg.DSN())
	if err != nil {
		log.Fatalf("redis parse url error: %s", err.Error())
	}

	// Redis 配置
	opts.DialTimeout = time.Duration(dialTimeout) * time.Second
	opts.ReadTimeout = time.Duration(readTimeout) * time.Second
	opts.WriteTimeout = time.Duration(writeTimeout) * time.Second
	opts.ConnMaxIdleTime = time.Duration(idleTimeout) * time.Second
	opts.PoolSize = poolSizeMultiple * runtime.NumCPU()
	opts.MinIdleConns = minIdleConnectionMultiple * runtime.NumCPU()

	// TLS 配置（ParseURL 会自动解析 scheme，若为 rediss 会提前初始化 opts.TLSConfig)
	if cfg.TLS.Enabled {
		// 再次确保 tls 已经初始化
		if opts.TLSConfig == nil {
			opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		// 服务器证书
		caCert, err := os.ReadFile(cfg.TLS.CertCaFile)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read ca cert: %s", cfg.TLS.CertCaFile)
		}
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM(caCert); !ok {
			return nil, errors.Errorf("failed to append ca cert: %s", cfg.TLS.CertCaFile)
		}
		opts.TLSConfig.RootCAs = pool
		// 证书验证配置
		opts.TLSConfig.InsecureSkipVerify = cfg.TLS.InsecureSkipVerify

		// 客户端证书
		if cfg.TLS.CertFile != "" && cfg.TLS.CertKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.CertKeyFile)
			if err != nil {
				return nil, errors.Wrapf(
					err, "failed to load x509 key pair, cert: %s, key: %s", cfg.TLS.CertFile, cfg.TLS.CertKeyFile,
				)
			}
			opts.TLSConfig.Certificates = []tls.Certificate{cert}
		}
	}
	return opts, nil
}

// InitRedisClient init redis client with config.RedisConfig
func InitRedisClient(ctx context.Context, cfg *config.RedisConfig) {
	if cfg == nil {
		log.Fatal("redis config is required when init redis client")
	}

	opts, err := buildOpts(cfg)
	if err != nil {
		log.Fatalf("unable to build redis options: %s", err.Error())
	}

	initOnce.Do(func() {
		rds = redis.NewClient(opts)
		if _, err = rds.Ping(ctx).Result(); err != nil {
			log.Fatalf("redis connect error: %s", err.Error())
		} else {
			log.Infof(ctx, "redis: %s:%d/%d connected", cfg.Host, cfg.Port, cfg.DB)
		}
		// OpenTelemetry Tracing
		if err = redisotel.InstrumentTracing(rds); err != nil {
			log.Fatalf("failed to enable redis tracing instrumentation: %s", err)
		}
	})
}

// Client 获取 redis 客户端
func Client() *redis.Client {
	if rds == nil {
		log.Fatal("redis client not init")
	}
	return rds
}
