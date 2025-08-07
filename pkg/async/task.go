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

// Package async 提供一个简单的异步 / 定时任务封装：
// 1. 使用 cron 支持定时任务（cmd: scheduler）
// 2. 简单封装 goroutine 以支持异步任务
package async

import (
	"context"
	"reflect"

	"bk.tencent.com/blueapps-go/pkg/async/task"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

// RegisteredTasks 已注册的任务
// 注意：任务函数最后一个返回值推荐为 error 类型
var RegisteredTasks = map[string]any{
	"CalcFib": task.CalcFib,
	// NOTE: SaaS 开发者可根据需求添加自定义任务
}

// ApplyTask 下发异步任务
func ApplyTask(ctx context.Context, name string, args []any) {
	go func() {
		taskFunc, ok := RegisteredTasks[name]
		if !ok {
			log.Errorf(ctx, "task func %s not found", name)
			return
		}

		taskArgs := []reflect.Value{reflect.ValueOf(ctx)}
		for _, arg := range args {
			taskArgs = append(taskArgs, reflect.ValueOf(arg))
		}
		values := reflect.ValueOf(taskFunc).Call(taskArgs)

		// 若任务执行有返回值，且最后一个返回值类型是 error 且不为 nil，需打印错误日志
		if length := len(values); length != 0 {
			if err, ok := values[length-1].Interface().(error); ok && err != nil {
				log.Errorf(ctx, "apply task %s with args %v error: %s", name, args, err)
			}
		}
	}()
}
