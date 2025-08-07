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

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

// NewExtractI18nMsgsCmd ...
func NewExtractI18nMsgsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "extract-i18n-msgs",
		Short: "Extract i18n messages from go source code and template",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			// 加载配置（不需要配置文件）
			_, err := config.Load(ctx, "")
			if err != nil {
				log.Fatalf("failed to load config: %s", err)
			}

			// 初始化 i18n
			i18n.InitMsgMap()
			// 提取 i18n 消息
			if err = i18n.ExtractMessages(); err != nil {
				log.Fatalf("failed to extract i18n messages: %s", err)
			}
			log.Infof(ctx, "extract i18n messages successfully")
			log.Infof(ctx, "placeholder `%s` in %s requires manual replacement", i18n.Placeholder, i18n.MsgFilepath())
		},
	}
}

func init() {
	rootCmd.AddCommand(NewExtractI18nMsgsCmd())
}
