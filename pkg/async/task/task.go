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

// Package task 包含异步任务实现
package task

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
)

// Fibonacci 斐波那契数的递归实现，因为性能很差所以适合模拟需要长时间运行的后台任务
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// CalcFib 计算斐波那契数任务
func CalcFib(ctx context.Context, n float64) (int, error) {
	// 由于 json Unmarshal 会把整数 & 浮点数都解析为 float64 类型，这由任务处理类型转换
	nInt := int(n)

	task := model.Task{
		Name:      "CalcFib",
		Args:      []byte(fmt.Sprintf("{\"n\": %d}", nInt)),
		StartedAt: time.Now(),
	}
	if err := database.Client(ctx).Create(&task).Error; err != nil {
		return 0, err
	}

	// 执行计算任务
	fibN := fibonacci(nInt)

	// 回填执行结果
	task.Result = []byte(strconv.Itoa(fibN))
	task.Duration = time.Since(task.StartedAt)
	if err := database.Client(ctx).Save(&task).Error; err != nil {
		return 0, err
	}

	return fibN, nil
}
