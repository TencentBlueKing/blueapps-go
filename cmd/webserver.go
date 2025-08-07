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
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/otel"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	"bk.tencent.com/blueapps-go/pkg/router"
)

// NewWebServerCmd ...
func NewWebServerCmd() *cobra.Command {
	var cfgFile string

	wsCmd := cobra.Command{
		Use:   "webserver",
		Short: "Start the HTTP server.",
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
				shutdown, sErr := otel.InitTracer(ctx, cfg.Platform.Addons.BkOtel, otel.GenServiceName("web"))
				if sErr != nil {
					log.Fatalf("failed to init OpenTelemetry: %s", sErr)
				}
				defer func() {
					if err = shutdown(ctx); err != nil {
						log.Fatalf("failed to shutdown OpenTelemetry: %s", err)
					}
				}()
			}

			// 启动 Web 服务
			log.Infof(ctx, "Starting server at http://0.0.0.0:%d", config.G.Service.Server.Port)
			srv := &http.Server{
				Addr:    ":" + strconv.Itoa(cfg.Service.Server.Port),
				Handler: router.New(log.GetLogger("gin")),
			}
			go func() {
				if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Start server failed: %s", err)
				}
			}()

			// 等待中断信号以优雅地关闭服务器
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit

			srvCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Service.Server.GraceTimeout)*time.Second)
			defer cancel()

			log.Info(ctx, "Shutdown server ...")
			if err = srv.Shutdown(srvCtx); err != nil {
				log.Fatalf("Shutdown server failed: %s", err)
			}
			log.Info(ctx, "Server exiting")
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	wsCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")

	return &wsCmd
}

func init() {
	rootCmd.AddCommand(NewWebServerCmd())
}
