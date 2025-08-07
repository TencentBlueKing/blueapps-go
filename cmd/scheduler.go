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

	"github.com/spf13/cobra"

	"bk.tencent.com/blueapps-go/pkg/async"
	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/otel"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

// NewSchedulerCmd 用于创建定时任务调度器启动命令
// 需要注意的是：为避免重复执行定时任务，需要确保同时只有一个 scheduler 正在运行
// 如果希望同时启动多个 scheduler 启动，则需要添加诸如 redis / zk 这样的分布式锁
func NewSchedulerCmd() *cobra.Command {
	var cfgFile string

	schedulerCmd := cobra.Command{
		Use:   "scheduler",
		Short: "Execute tasks based on cron expressions, please ensure only one running scheduler.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			// 初始化 i18n
			i18n.InitMsgMap()
			// 初始化 Logger
			if err = initLogger(&cfg.Service.Log); err != nil {
				log.Fatalf("failed to init logging: %s", err)
			}
			// 初始化增强服务客户端
			if err = initAddons(ctx, cfg); err != nil {
				log.Fatalf("failed to init addons: %s", err)
			}
			// 初始化 OpenTelemetry
			if cfg.Platform.Addons.BkOtel != nil {
				shutdown, sErr := otel.InitTracer(ctx, cfg.Platform.Addons.BkOtel, otel.GenServiceName("scheduler"))
				if sErr != nil {
					log.Fatalf("failed to init OpenTelemetry: %s", sErr)
				}
				defer func() {
					if err = shutdown(ctx); err != nil {
						log.Fatalf("failed to shutdown OpenTelemetry: %s", err)
					}
				}()
			}

			// 初始化 task server
			async.InitTaskScheduler(ctx)

			srv := async.Scheduler()
			// 加载周期任务
			if err = srv.LoadTasks(); err != nil {
				log.Fatal(err.Error())
			}
			// 启用调度服务器
			srv.Run()
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	schedulerCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")

	return &schedulerCmd
}

func init() {
	rootCmd.AddCommand(NewSchedulerCmd())
}
