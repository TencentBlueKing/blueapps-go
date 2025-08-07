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
	"sync"

	"github.com/robfig/cron/v3"
)

type entry struct {
	id   cron.EntryID
	name string
}

// 任务 ID 与 cron.EntryID + 任务名称的映射表
type taskEntryMap struct {
	mapping map[int64]entry
	sync.RWMutex
}

func (m *taskEntryMap) get(taskID int64) (entry, bool) {
	m.RLock()
	defer m.RUnlock()
	e, ok := m.mapping[taskID]
	return e, ok
}

func (m *taskEntryMap) set(taskID int64, entry entry) {
	m.Lock()
	defer m.Unlock()
	m.mapping[taskID] = entry
}

func (m *taskEntryMap) delete(taskID int64) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, taskID)
}
