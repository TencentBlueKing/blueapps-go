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

// Package handler ...
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/apis/asynctask/serializer"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// ListPeriodicTasks ...
//
//	@Summary	获取定时任务列表
//	@Tags		async-task
//	@Success	200	{object}	ginx.Response{data=[]serializer.PeriodicTaskListResponse}
//	@Router		/api/periodic-tasks [get]
func ListPeriodicTasks(c *gin.Context) {
	var periodicTasks []model.PeriodicTask
	if err := database.Client(c.Request.Context()).Order("id DESC").Find(&periodicTasks).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.PeriodicTaskListResponse{}
	for _, task := range periodicTasks {
		respData = append(respData, serializer.PeriodicTaskListResponse{
			ID:      task.ID,
			Cron:    task.Cron,
			Name:    task.Name,
			Args:    string(task.Args),
			Enabled: task.Enabled,
			Creator: task.Creator,
		})
	}
	ginx.SetResp(c, http.StatusOK, respData)
}

// CreatePeriodicTask ...
//
//	@Summary	创建定时任务
//	@Tags		async-task
//	@Param		body	body		serializer.PeriodicTaskCreateRequest	true	"定时任务配置"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/periodic-tasks [post]
func CreatePeriodicTask(c *gin.Context) {
	var req serializer.PeriodicTaskCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	args, _ := json.Marshal(req.Args)
	periodicTask := model.PeriodicTask{
		Cron: req.Cron,
		Name: req.Name,
		Args: args,
		BaseModel: model.BaseModel{
			Creator: ginx.GetUserID(c),
			Updater: ginx.GetUserID(c),
		},
	}
	if err := database.Client(c.Request.Context()).Create(&periodicTask).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginx.SetResp(c, http.StatusCreated, nil)
}

// DeletePeriodicTask ...
//
//	@Summary	删除定时任务
//	@Tags		async-task
//	@Param		id	path	int	true	"定时任务 ID"
//	@Success	204	"No Content"
//	@Router		/api/periodic-tasks/{id} [delete]
func DeletePeriodicTask(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&model.PeriodicTask{})
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx.SetResp(c, http.StatusNoContent, nil)
}

// TogglePeriodicTaskEnabled ...
//
//	@Summary	切换定时任务启用状态
//	@Tags		async-task
//	@Param		id	path	int	true	"定时任务 ID"
//	@Success	204	"No Content"
//	@Router		/api/periodic-tasks/{id}/enabled [put]
func TogglePeriodicTaskEnabled(c *gin.Context) {
	var periodicTask model.PeriodicTask
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", c.Param("id")).First(&periodicTask)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	periodicTask.Enabled = !periodicTask.Enabled
	periodicTask.Updater = ginx.GetUserID(c)

	tx = database.Client(ctx).Save(&periodicTask)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx.SetResp(c, http.StatusOK, serializer.TogglePeriodicTaskEnabledResponse{Enabled: periodicTask.Enabled})
}
