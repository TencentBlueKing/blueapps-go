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

package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/apis/asynctask/serializer"
	"bk.tencent.com/blueapps-go/pkg/async"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// ListTasks ...
//
//	@Summary	获取任务列表
//	@Tags		async-task
//	@Success	200	{object}	ginx.Response{data=ginx.PaginatedResp{results=[]serializer.TaskListResponse}}
//	@Router		/api/tasks [get]
func ListTasks(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Order("created_at desc").Model(&model.Task{})

	// 总条目数量
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
	}

	var executedTasks []model.Task
	if err := tx.Offset(ginx.GetOffset(c)).Limit(ginx.GetLimit(c)).Find(&executedTasks).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.TaskListResponse{}
	for _, task := range executedTasks {
		respData = append(respData, serializer.TaskListResponse{
			ID:        task.ID,
			Name:      task.Name,
			Args:      string(task.Args),
			Result:    string(task.Result),
			Creator:   task.Creator,
			StartedAt: lo.Ternary(task.StartedAt.IsZero(), "", task.StartedAt.Format(time.RFC3339)),
			Duration:  task.Duration.Seconds(),
		})
	}
	ginx.SetResp(c, http.StatusOK, ginx.NewPaginatedRespData(total, respData))
}

// CreateTask ...
//
//	@Summary	创建异步任务
//	@Tags		async-task
//	@Param		body	body		serializer.TaskCreateRequest	true	"异步任务配置"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/tasks [post]
func CreateTask(c *gin.Context) {
	var req serializer.TaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	// 异步任务执行，不使用 c.Request.Context() 以避免提前 cancel
	async.ApplyTask(context.Background(), req.Name, req.Args)
	ginx.SetResp(c, http.StatusCreated, nil)
}
