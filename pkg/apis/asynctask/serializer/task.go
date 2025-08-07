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

package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bk.tencent.com/blueapps-go/pkg/async"
	"bk.tencent.com/blueapps-go/pkg/i18n"
)

// TaskListResponse List Task API 返回结构
type TaskListResponse struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Args      string  `json:"args"`
	Result    string  `json:"result"`
	Creator   string  `json:"creator"`
	StartedAt string  `json:"startedAt"`
	Duration  float64 `json:"duration"`
}

// TaskCreateRequest Create Task API 请求结构
type TaskCreateRequest struct {
	Name string `json:"name"`
	Args []any  `json:"args"`
}

// Validate ...`
func (r *TaskCreateRequest) Validate(c *gin.Context) error {
	if r.Name == "" {
		return errors.New(i18n.T(c.Request.Context(), "Task name required"))
	}
	if _, ok := async.RegisteredTasks[r.Name]; !ok {
		return errors.Errorf("Task name %s invalid", r.Name)
	}
	return nil
}
