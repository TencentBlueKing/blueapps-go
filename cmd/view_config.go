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
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"bk.tencent.com/blueapps-go/pkg/config"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

// NewViewConfigCmd ...
func NewViewConfigCmd() *cobra.Command {
	var cfgFile string
	var verbose bool

	viewCfgCmd := cobra.Command{
		Use:   "view-config",
		Short: "View service configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			// 加载配置
			cfg, err := config.Load(context.Background(), cfgFile)
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			if verbose {
				spew.Dump(cfg)
				return
			}

			data, _ := yaml.Marshal(cfg)
			fmt.Println(string(data))
		},
	}

	// 配置文件路径，如果未指定，会从环境变量读取各项配置
	// 注意：目前平台未默认提供配置文件，需通过 `模块配置 - 挂载卷` 添加
	viewCfgCmd.Flags().StringVar(&cfgFile, "conf", "", "config file")
	viewCfgCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show more details")

	return &viewCfgCmd
}

func init() {
	rootCmd.AddCommand(NewViewConfigCmd())
}
