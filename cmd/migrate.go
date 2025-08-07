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
	"log/slog"

	"github.com/spf13/cobra"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	// load migration package to register migrations
	_ "bk.tencent.com/blueapps-go/pkg/migration"
	"bk.tencent.com/blueapps-go/pkg/version"
)

// NewMigrateCmd ...
func NewMigrateCmd() *cobra.Command {
	var cfgFile string
	var migrationID string

	migrateCmd := cobra.Command{
		Use:   "migrate",
		Short: "Apply migrations to the database tables.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			if cfg.Platform.Addons.Mysql == nil {
				log.Fatal("mysql config not found, skip migrate...")
			}

			database.InitDBClient(ctx, cfg.Platform.Addons.Mysql, slog.Default())

			if err = database.RunMigrate(ctx, migrationID); err != nil {
				log.Fatalf("failed to run migrate: %s", err)
			}
			dbVersion, err := database.Version(ctx)
			if err != nil {
				log.Fatalf("failed to get database version: %s", err)
			}
			log.Infof(ctx, "migrate success %s\nDatabaseVersion: %s", version.Version(), dbVersion)
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	migrateCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")
	migrateCmd.Flags().StringVar(&migrationID, "migration", "", "migration to apply, blank means latest version")

	return &migrateCmd
}

func init() {
	rootCmd.AddCommand(NewMigrateCmd())
}
