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

package model

import (
	"time"

	"gorm.io/datatypes"
)

// Task 后台任务
type Task struct {
	BaseModel
	ID        int64          `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(128);not null"`
	Args      datatypes.JSON `json:"args" gorm:"type:json"`
	Result    datatypes.JSON `json:"result" gorm:"type:json"`
	StartedAt time.Time      `json:"startedAt" gorm:"type:datetime;default:null"`
	Duration  time.Duration  `json:"duration" gorm:"type:bigint;default:null"`
}

// PeriodicTask 周期任务
type PeriodicTask struct {
	BaseModel
	ID      int64          `json:"id" gorm:"primaryKey"`
	Cron    string         `json:"cron" gorm:"type:varchar(32);not null"`
	Name    string         `json:"name" gorm:"type:varchar(128);not null"`
	Args    datatypes.JSON `json:"args" gorm:"type:json"`
	Enabled bool           `json:"enabled" gorm:"not null;default:true"`
}
