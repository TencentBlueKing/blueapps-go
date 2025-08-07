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

package config

import (
	"fmt"
	"net/url"

	"github.com/samber/lo"
)

// TLSConfig tls 证书相关配置
type TLSConfig struct {
	Enabled            bool   // 是否启用 TLS 证书连接
	CertCaFile         string // CA 证书文件路径
	CertFile           string // Cert 证书文件路径
	CertKeyFile        string // Cert Key 文件路径
	InsecureSkipVerify bool   // 是否跳过 TLS 校验（不推荐在生产环境使用）
}

// MysqlConfig Mysql 增强服务配置
type MysqlConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	Charset  string
	TLS      TLSConfig
}

// DSN ...
func (cfg *MysqlConfig) DSN() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s&time_zone=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		// 指定从 MySQL 读取的时间转到 go 的 time.Time 所用时区
		url.QueryEscape("UTC"),
		// 指定连接 MySQL 的时区
		// Q：为什么是使用指定时区，而不使用命名时区，比如 UTC、Asia/Shanghai 等
		// A：仅当 mysql 中存在时区信息表时，才能使用命名时区；所以使用指定时区可适应所有场景
		url.QueryEscape("'+00:00'"),
	)
	if cfg.TLS.Enabled {
		dsn = fmt.Sprintf("%s&tls=%s", dsn, cfg.TLSCfgName())
	}
	return dsn
}

// TLSCfgName mysql tls 配置名称
func (cfg *MysqlConfig) TLSCfgName() string {
	return "custom"
}

// RabbitMQConfig RabbitMQ 增强服务配置
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Vhost    string
	Password string
	TLS      TLSConfig
}

// DSN ...
func (cfg *RabbitMQConfig) DSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		lo.Ternary(cfg.TLS.Enabled, "amqps", "amqp"),
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Vhost,
	)
}

// RedisConfig Redis 增强服务配置
type RedisConfig struct {
	Username string
	Host     string
	Port     int
	Password string
	DB       int
	TLS      TLSConfig
}

// DSN ...
func (cfg *RedisConfig) DSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%d",
		lo.Ternary(cfg.TLS.Enabled, "rediss", "redis"),
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}

// BkRepoConfig BkRepo 增强服务配置
type BkRepoConfig struct {
	EndpointUrl   string
	Project       string
	Username      string
	Password      string
	Bucket        string
	PublicBucket  string
	PrivateBucket string
}

// BkOtelConfig BkOtel 增强服务配置
type BkOtelConfig struct {
	Trace       bool
	GrpcUrl     string
	BkDataToken string
	Sampler     string
}

// AddonsConfig 增强服务配置
type AddonsConfig struct {
	Mysql    *MysqlConfig
	RabbitMQ *RabbitMQConfig
	Redis    *RedisConfig
	BkRepo   *BkRepoConfig
	BkOtel   *BkOtelConfig
}

// BkPlatUrlConfig 蓝鲸各平台服务地址
type BkPlatUrlConfig struct {
	// 蓝鲸开发者中心地址
	BkPaaS string
	// 统一登录地址
	BkLogin string
	// 组件 API 地址
	BkCompApi string
	// NOTE: SaaS 开发者可按需添加诸如 BkIAM，BkLog 等服务配置
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	// 蓝鲸应用 ID
	AppID string
	// 蓝鲸应用密钥
	AppSecret string
	// 模块名称
	ModuleName string
	// 运行环境：stag 预发布环境，prod 生产环境
	RunEnv string

	// 应用引擎版本
	Region string
	// 推荐的 DB 加密算法有：SHANGMI（对应 SM4CTR 算法）和 CLASSIC（对应 Fernet 算法）
	CryptoType string

	// 蓝鲸根域名，用于获取登录票据，国际化语言等 cookie 信息
	BkDomain string
	// 网关 API 访问地址模板
	ApiUrlTmpl string

	// 蓝鲸平台服务地址配置
	BkPlatUrl BkPlatUrlConfig
	// 增强服务配置
	Addons AddonsConfig
}

// LogConfig 日志配置
type LogConfig struct {
	// 日志级别，可选值为：debug、info、warn、error
	Level string
	// 日志目录，部署于 PaaS 平台上时，该值必须为 /app/v3logs，否则无法采集日志
	Dir string
	// 是否强制标准输出，不输出到文件（一般用于本地开发，标准输出日志查看比较方便）
	ForceToStdout bool
}

// ServerConfig Gin Web Server 配置
type ServerConfig struct {
	// 服务端口
	Port int
	// 优雅退出等待时间
	GraceTimeout int
	// Gin 运行模式
	GinRunMode string
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	// Web Server 配置
	Server ServerConfig
	// 日志配置
	Log LogConfig

	// CORS 允许来源列表
	AllowedOrigins []string
	// 允许访问的用户列表（UserID），为空时表示不限制
	AllowedUsers []string
	// DB 加密密钥，若未使用加密功能可不配置
	// 生成方式参见 Readme 文档 - 数据库字段加密
	EncryptSecret string
	// 用户认证方式
	// 目前支持：BkTicket、BkToken、Taihu
	// 按顺序尝试认证，最后一个会尝试跳转登录页面
	// 如果包含 Taihu，则需要配置 TaihuAppToken
	// 示例：["Taihu", "BkToken"]
	AuthTypes []string
	// Taihu 应用 Token，用于验证用户身份
	// 可在太湖 - 应用概览 - 应用信息处获取
	TaihuAppToken string
	// Taihu 非安全模式（明文/兼容模式）
	// 必须与 Taihu 应用上设置的保持一致
	TaihuInsecure bool
	// CSRF Cookie 域名
	CSRFCookieDomain string
	// 健康探针 Token
	HealthzToken string
	// 指标 API Token
	MetricToken string

	// 缓存内存大小（单位为 MB）
	MemoryCacheSize int

	// 是否启用 swagger docs
	EnableSwagger bool
	// API 文档文件存放目录
	ApiDocFileBaseDir string
	// 静态文件存放目录
	StaticFileBaseDir string
	// 国际化文件存放目录
	I18nFileBaseDir string
	// 模板文件存放目录
	TmplFileBaseDir string
}

// BizConfig 业务相关配置
type BizConfig struct {
	// NOTE: SaaS 开发者可在此处添加业务相关配置项
}

// Config SaaS 配置
type Config struct {
	// 平台内置配置
	Platform PlatformConfig
	// 服务配置
	Service ServiceConfig
	// 业务配置
	Biz BizConfig
}
