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

package async

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"

	"bk.tencent.com/blueapps-go/pkg/infras/database"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	"bk.tencent.com/blueapps-go/pkg/model"
)

var srv *TaskScheduler

var initOnce sync.Once

/*
* Q：为什么不是直接接入诸如 machinery / asynq 这样的任务调度框架？
* A：开发框架作者实现过 machinery 的接入示例，但是发现有以下问题
*     - machinery 是一个比较重的框架，虽然支持多种 broker / backend，要用起来还是有一定的学习成本
*     - machinery 的维护情况不是很理想，测试覆盖率一般，且 RedisLock 实现有明显 Bug，需要慎重使用
*     - Golang 语言原生支持异步（goroutine），可以很方便地实现任务异步执行，大多数场景下不需要引入框架
*     - Cron（https://github.com/robfig/cron）库提供强大的定时功能，简单配置下即可支持常见的定时任务场景
*    综上，我们在讨论后移除了对 machinery 的引入，仅仅作为文档中的示例供有需要的开发者参考
*
* Q：目前这套基于 goroutine 实现的机制会有什么问题
* A：由于没有使用消息队列，也没有保护机制，因此如果进程重启/崩溃，会导致运行中的任务中断
*
* Q：scheduler 是如何管理周期任务的？
* A：- scheduler 首次启动时，会从 DB 中加载所有周期任务，并根据指定的 Cron 表达式执行
#    - scheduler 会根据指定的时间间隔（reloadTasksCron）周期性从 DB 中加载新增的任务
*
* Q：如果我想接入如 machinery 这样的异步框架，应该怎么做？比如怎么适配增强服务？
# A：请查阅 Readme.md 中的 `异步/定时任务` 一节
*/

// 默认每 5 分钟检查 & 重载周期任务
// NOTE: SaaS 开发者可以根据需要自行调整，但不建议过大/过小
const reloadTasksCron = "*/5 * * * *"

var tracer = otel.Tracer("task-scheduler")

// TaskScheduler 简单的定时任务调度器，依赖 robfig/cron & model.PeriodicTask
type TaskScheduler struct {
	ctx          context.Context
	cron         *cron.Cron
	taskEntryMap *taskEntryMap
}

// Run 启用调度器
func (s *TaskScheduler) Run() {
	s.cron.Run()
}

// LoadTasks 加载所有周期任务
func (s *TaskScheduler) LoadTasks() error {
	// 从数据库加载周期性任务
	periodicTasks := []model.PeriodicTask{}
	if err := database.Client(s.ctx).Find(&periodicTasks).Error; err != nil {
		return errors.Wrap(err, "load periodic tasks")
	}

	enabledTaskIDs := mapset.NewSet[int64]()
	// 根据是否启用，注册/注销周期任务
	for _, task := range periodicTasks {
		if !task.Enabled {
			srv.unregister(task.ID)
			continue
		}
		enabledTaskIDs.Add(task.ID)
		if err := srv.register(task); err != nil {
			return errors.Wrap(err, "register periodic task")
		}
	}

	// 对于已经注册但 DB 中已经删除的，需要注销
	for taskID := range s.taskEntryMap.mapping {
		if !enabledTaskIDs.Contains(taskID) {
			srv.unregister(taskID)
		}
	}
	log.Debugf(s.ctx, "%d periodic tasks loaded", len(s.taskEntryMap.mapping))
	return nil
}

// 注册单个周期任务
func (s *TaskScheduler) register(task model.PeriodicTask) error {
	// 跳过已注册的任务
	if _, ok := s.taskEntryMap.get(task.ID); ok {
		return nil
	}

	taskRepr := fmt.Sprintf("periodic task %s (id: %d)", task.Name, task.ID)
	log.Infof(s.ctx, "register %s with cron: %s args: %v", taskRepr, task.Cron, task.Args)

	entryID, err := s.cron.AddFunc(task.Cron, func() {
		ctx, span := tracer.Start(s.ctx, taskRepr)
		defer span.End()

		// 已注册任务不存在 -> 已被删除，但还没重载刷新，可以跳过，其他错误需要打印错误日志
		if err := database.Client(ctx).First(&task, task.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Infof(ctx, "%s not found in database, skip run...", taskRepr)
			} else {
				log.Errorf(ctx, "failed to reload %s from database: %s", taskRepr, err)
			}
			return
		}
		// 被禁用的已注册任务，在重载前需要跳过
		if !task.Enabled {
			log.Infof(ctx, "%s is disabled, skip run...", taskRepr)
			return
		}
		// 解析参数 & 下发异步任务
		var taskArgs []any
		if err := json.Unmarshal(task.Args, &taskArgs); err != nil {
			log.Errorf(ctx, "failed to unmarshal %s args: %s", taskRepr, err)
			return
		}
		ApplyTask(ctx, task.Name, taskArgs)
	})
	if err != nil {
		return err
	}

	s.taskEntryMap.set(task.ID, entry{id: entryID, name: task.Name})
	return nil
}

// 注销单个周期任务
func (s *TaskScheduler) unregister(taskID int64) {
	// 跳过未注册的任务
	e, ok := s.taskEntryMap.get(taskID)
	if !ok {
		return
	}

	log.Infof(s.ctx, "unregister periodic task %s (id: %d)", e.name, taskID)

	s.cron.Remove(e.id)
	s.taskEntryMap.delete(taskID)
}

// Scheduler 获取调度器
func Scheduler() *TaskScheduler {
	if srv == nil {
		log.Fatal("task server not init")
	}
	return srv
}

// newScheduler 创建调度器
func newScheduler(ctx context.Context) (*TaskScheduler, error) {
	return &TaskScheduler{
		ctx: ctx,
		// 注：如有需要，可使用 cron.WithSeconds() 支持秒级精度控制
		// cron 默认使用服务部署机器时区（time.Local），这里显式指定更清晰
		// ref: https://github.com/robfig/cron
		cron: cron.New(
			cron.WithLocation(time.Local),
		),
		// taskEntryMap 提供并发场景下的读写保护
		taskEntryMap: &taskEntryMap{
			mapping: make(map[int64]entry),
		},
	}, nil
}

// InitTaskScheduler 初始化任务调度器
func InitTaskScheduler(ctx context.Context) {
	if srv != nil {
		return
	}
	initOnce.Do(func() {
		var err error
		srv, err = newScheduler(ctx)
		if err != nil {
			log.Fatalf("failed to init task server: %s", err)
		}
		// 添加周期任务：定时从数据库加载定义的周期任务
		_, err = srv.cron.AddFunc(reloadTasksCron, func() {
			if err = srv.LoadTasks(); err != nil {
				log.Warnf(ctx, "failed to reload periodic tasks: %s", err)
			}
		})
		if err != nil {
			log.Fatalf("failed to add reload tasks periodic task: %s", err)
		}
		log.Info(ctx, "task server initialized")
	})
}
