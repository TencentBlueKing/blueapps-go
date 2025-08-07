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

package i18n

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"bk.tencent.com/blueapps-go/pkg/config"
)

var initOnce sync.Once

// 国际化字典
var i18nMsgMap map[string]map[Lang]string

// InitMsgMap 服务启动时初始化 i18n 配置
func InitMsgMap() {
	if i18nMsgMap != nil {
		return
	}

	initOnce.Do(func() {
		// 读取国际化配置文件
		yamlFile, err := os.ReadFile(MsgFilepath())
		if err != nil {
			log.Fatalf("failed to read i18n file: %v", err)
		}

		var rawMsgList []map[string]string
		if err = yaml.Unmarshal(yamlFile, &rawMsgList); err != nil {
			log.Fatalf("failed to unmarshal i18n file: %v", err)
		}

		// 初始化国际化字典
		i18nMsgMap = map[string]map[Lang]string{}
		for _, rawMsg := range rawMsgList {
			msgID, ok := rawMsg["id"]
			if !ok {
				continue
			}
			ms := map[Lang]string{}
			for _, lang := range supportedLanguages {
				if _, ok = rawMsg[string(lang)]; ok {
					ms[lang] = rawMsg[string(lang)]
				}
			}
			i18nMsgMap[msgID] = ms
		}
	})
}

// MsgFilepath 获取国际化配置文件路径
func MsgFilepath() string {
	return filepath.Join(config.G.Service.I18nFileBaseDir, "messages.yaml")
}
