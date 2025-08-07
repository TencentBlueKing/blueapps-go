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
	"bk.tencent.com/blueapps-go/pkg/model"
	"bk.tencent.com/blueapps-go/pkg/version"
)

// NewInitDataCmd ...
func NewInitDataCmd() *cobra.Command {
	var cfgFile string

	migrateCmd := cobra.Command{
		Use:   "init-data",
		Short: "Initialize the database with initial data after the first deployment.",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置
			cfg, err := config.Load(ctx, cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			if cfg.Platform.Addons.Mysql == nil {
				log.Fatal("mysql config not found, skip init data...")
			}

			database.InitDBClient(ctx, cfg.Platform.Addons.Mysql, slog.Default())

			if err = initDemoData(ctx); err != nil {
				log.Fatalf("failed to init demo data: %s", err)
			}
			log.Infof(ctx, "demo data initialized %s", version.Version())
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	migrateCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")

	return &migrateCmd
}

// NOTE: SaaS 开发者可以自定义需要初始化的数据
func initDemoData(ctx context.Context) error { // nolint
	baseModel := model.BaseModel{
		Creator: "admin",
		Updater: "admin",
	}

	categories := []model.Category{
		{
			Name:      "fruit",
			BaseModel: baseModel,
			Entries: []model.Entry{
				{
					Name: "Apple",
					Desc: "Apple is a sweet, edible fruit produced by an apple tree, " +
						"typically red, green, or yellow in color.",
					Price:     6.99,
					BaseModel: baseModel,
				},
				{
					Name: "Banana",
					Desc: "Banana is a long, curved fruit with a yellow peel and soft, " +
						"sweet flesh inside, produced by the banana plant.",
					Price:     3.49,
					BaseModel: baseModel,
				},
				{
					Name: "Orange",
					Desc: "Orange is a round, juicy citrus fruit with " +
						"a tough bright orange rind and a sweet-tart flavor.",
					Price:     4.69,
					BaseModel: baseModel,
				},
				{
					Name: "Peach",
					Desc: "Peach is a round, juicy fruit with a fuzzy skin " +
						"and sweet flesh, typically yellow or white in color.",
					Price:     5.79,
					BaseModel: baseModel,
				},
			},
		},
		{
			Name:      "book",
			BaseModel: baseModel,
			Entries: []model.Entry{
				{
					Name: "The Origin of Species",
					Desc: "\"On the Origin of Species\" overturned creationism and the fixity of species " +
						"with a revolutionary theory of evolution, establishing biology on a scientific foundation.",
					Price:     1859.02,
					BaseModel: baseModel,
				},
				{
					Name: "The Influence of Sea Power Upon History",
					Desc: "\"The Influence of Sea Power upon History\" summarizes and studies the strategies " +
						"and tactics of naval warfare throughout history and their impacts, " +
						"proposing that control of the sea determines the rise and fall of a nation's fortunes.",
					Price:     1890.09,
					BaseModel: baseModel,
				},
				{
					Name: "Relativity: The Special and General Theory",
					Desc: "\"Relativity\" is a groundbreaking work written by the scientist Albert Einstein, " +
						"which completely overturned the concepts of classical physics.",
					Price:     1916.03,
					BaseModel: baseModel,
				},
				{
					Name: "Introduction to Interstellar Travel",
					Desc: "\"Introduction to Interstellar Travel\" provides a comprehensive introduction to " +
						"the complexity and challenges of interstellar travel technology and practice.",
					Price:     1963.12,
					BaseModel: baseModel,
				},
			},
		},
		{
			Name:      "ball",
			BaseModel: baseModel,
			Entries: []model.Entry{
				{
					Name: "Football",
					Desc: "Football is a team sport where two teams of eleven players each " +
						"try to score goals by kicking a ball into the opposing team’s net.",
					Price:     599.99,
					BaseModel: baseModel,
				},
				{
					Name: "Basketball",
					Desc: "Basketball is a team sport where two teams of five players each " +
						"try to score points by shooting a ball through the opposing team’s hoop.",
					Price:     499.99,
					BaseModel: baseModel,
				},
				{
					Name: "Volleyball",
					Desc: "Volleyball is a team sport where two teams of six players each " +
						"try to score points by hitting a ball over a net into the opposing team’s court.",
					Price:     399.99,
					BaseModel: baseModel,
				},
			},
		},
	}
	return database.Client(ctx).Create(categories).Error
}

func init() {
	rootCmd.AddCommand(NewInitDataCmd())
}
