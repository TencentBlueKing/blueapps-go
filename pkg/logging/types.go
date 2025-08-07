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

// Gin, Gorm 默认日志等级为 warn，目的是避免记录过多无关日志
// 开发者可根据需求自行调整，可选值：debug、info、warn、error
const (
	// GinLogLevel gin 日志级别
	GinLogLevel = "warn"
	// GormLogLevel gorm 日志级别
	GormLogLevel = "warn"
)

// Options Logger 配置
type Options struct {
	// 日志级别
	Level string
	// 日志内容 Handler，支持 text 和 json
	HandlerName string
	// io.Writer, 支持 stdout、stderr、file
	WriterName string
	// Writer 配置
	WriterConfig map[string]string
}
