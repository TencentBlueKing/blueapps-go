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

package logging

import (
	"fmt"
	"log/slog"
	"strings"
)

var loggers map[string]*slog.Logger

// GetLogger 获取指定 Logger
func GetLogger(name string) *slog.Logger {
	if logger, ok := loggers[name]; ok {
		return logger
	}

	// 不存在则返回默认的
	return slog.Default()
}

// InitLogger ...
func InitLogger(name string, opts *Options) (err error) {
	if loggers == nil {
		loggers = make(map[string]*slog.Logger)
	}

	// 已存在，则忽略，不需要再初始化
	if _, ok := loggers[name]; ok {
		return nil
	}

	if loggers[name], err = newLogger(opts); err != nil {
		return err
	}
	if name == "default" {
		// SetDefault 会改变 Golang slog 的 默认 logging
		// 同时会改变 Golang log 包使用的默认 log.Logger
		slog.SetDefault(loggers[name])
	}

	return nil
}

// 根据配置生成 Logger
func newLogger(opts *Options) (*slog.Logger, error) {
	w, err := newWriter(opts.WriterName, opts.WriterConfig)
	if err != nil {
		return nil, err
	}

	level, err := toSlogLevel(opts.Level)
	if err != nil {
		return nil, err
	}
	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// groups 长度不为 0 -> 非顶层字段，无需处理
			if len(groups) != 0 {
				return a
			}
			// 替换部分字段以适配蓝鲸日志平台清洗规则
			switch a.Key {
			case slog.MessageKey:
				a.Key = "message"
			case slog.LevelKey:
				a.Key = "levelname"
			case slog.SourceKey:
				a.Key = "pathname"
			}
			return a
		},
	}

	switch opts.HandlerName {
	case "text":
		return slog.New(slog.NewTextHandler(w, handlerOpts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, handlerOpts)), nil
	}

	return nil, fmt.Errorf("[%s] handler not supported", opts.HandlerName)
}

// toSlogLevel 将配置输入的日志级别转换为 slog Level 对象
func toSlogLevel(level string) (slog.Level, error) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "WARN", "WARNING":
		return slog.LevelWarn, nil
	case "ERROR":
		return slog.LevelError, nil
	}

	return slog.LevelInfo, fmt.Errorf("[%s] level not supported", level)
}
