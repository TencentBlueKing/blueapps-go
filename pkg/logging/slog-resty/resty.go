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

// Package slogresty 实现 resty.Logger 接口
package slogresty

import (
	"context"
	"log/slog"

	"github.com/go-resty/resty/v2"

	"bk.tencent.com/blueapps-go/pkg/logging"
)

// Logger 用于实现 resty.Logger
type Logger struct {
	ctx context.Context
}

// New 实例化 Logger
func New(ctx context.Context) *Logger {
	return &Logger{ctx: ctx}
}

// Errorf ...
func (l *Logger) Errorf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelError, format, v...)
}

// Warnf ...
func (l *Logger) Warnf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelWarn, format, v...)
}

// Debugf ...
func (l *Logger) Debugf(format string, v ...any) {
	logging.Logf(l.ctx, slog.LevelDebug, format, v...)
}

var _ resty.Logger = (*Logger)(nil)
