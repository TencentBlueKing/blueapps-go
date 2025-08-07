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
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bk.tencent.com/blueapps-go/pkg/apis/crud/serializer"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// ListEntries ...
//
//	@Summary	获取条目列表
//	@Tags		crud
//	@Success	200	{object}	ginx.Response{data=ginx.PaginatedResp{results=[]serializer.EntryListResponse}}
//	@Router		/api/entries [get]
func ListEntries(c *gin.Context) {
	var req serializer.EntryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	tx := database.Client(c.Request.Context()).Model(&model.Entry{}).Preload("Category")
	if req.CategoryID != 0 {
		tx = tx.Where("category_id = ?", req.CategoryID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		tx = tx.Where("LOWER(name) LIKE ?", keyword).
			Or("LOWER(`desc`) LIKE ?", keyword).
			Or("LOWER(updater) LIKE ?", keyword)
	}

	// 总条目数量
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
	}

	// 分页对应数据
	var entries []model.Entry
	if err := tx.Offset(ginx.GetOffset(c)).Limit(ginx.GetLimit(c)).Find(&entries).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.EntryListResponse{}
	for _, entry := range entries {
		respData = append(respData, serializer.EntryListResponse{
			CategoryID:   entry.CategoryID,
			CategoryName: entry.Category.Name,
			ID:           entry.ID,
			Name:         entry.Name,
			Desc:         entry.Desc,
			Price:        entry.Price,
			Updater:      entry.Updater,
			UpdatedAt:    entry.UpdatedAt.Format(time.RFC3339),
		})
	}
	ginx.SetResp(c, http.StatusOK, ginx.NewPaginatedRespData(total, respData))
}

// CreateEntry ...
//
//	@Summary	创建条目
//	@Tags		crud
//	@Param		body	body		serializer.EntryCreateRequest	true	"创建条目请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/entries [post]
func CreateEntry(c *gin.Context) {
	var req serializer.EntryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	// 检查分类是否存在
	var category model.Category
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", req.CategoryID).First(&category)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			ginx.SetErrResp(
				c,
				http.StatusNotFound,
				fmt.Sprintf(i18n.T(ctx, "category %d not found"), req.CategoryID),
			)
			return
		}
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	entry := model.Entry{
		Name:       req.Name,
		Desc:       req.Desc,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		BaseModel: model.BaseModel{
			Creator: ginx.GetUserID(c),
			Updater: ginx.GetUserID(c),
		},
	}
	if err := database.Client(ctx).Create(&entry).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginx.SetResp(c, http.StatusCreated, serializer.EntryCreateResponse{ID: entry.ID})
}

// RetrieveEntry ...
//
//	@Summary	获取单个条目
//	@Tags		crud
//	@Param		id	path		int	true	"条目 ID"
//	@Success	200	{object}	ginx.Response{data=serializer.EntryRetrieveResponse}
//	@Router		/api/entries/{id} [get]
func RetrieveEntry(c *gin.Context) {
	var entry model.Entry

	tx := database.Client(c.Request.Context()).Preload("Category").Where("id = ?", c.Param("id")).First(&entry)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	respData := serializer.EntryRetrieveResponse{
		// 分类属性
		CategoryID:   entry.CategoryID,
		CategoryName: entry.Category.Name,
		// 条目属性
		ID:        entry.ID,
		Name:      entry.Name,
		Desc:      entry.Desc,
		Price:     entry.Price,
		Creator:   entry.Creator,
		Updater:   entry.Updater,
		CreatedAt: entry.CreatedAt.Format(time.RFC3339),
		UpdatedAt: entry.UpdatedAt.Format(time.RFC3339),
	}
	ginx.SetResp(c, http.StatusOK, respData)
}

// UpdateEntry ...
//
//	@Summary	更新条目
//	@Tags		crud
//	@Param		id		path	int								true	"条目 ID"
//	@Param		body	body	serializer.EntryUpdateRequest	true	"更新条目请求体"
//	@Success	204		"No Content"
//	@Router		/api/entries/{id} [put]
func UpdateEntry(c *gin.Context) {
	var req serializer.EntryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	var entry model.Entry
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", c.Param("id")).First(&entry)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	// 更新 DB 模型字段
	entry.CategoryID = req.CategoryID
	entry.Name = req.Name
	entry.Desc = req.Desc
	entry.Price = req.Price
	entry.Updater = ginx.GetUserID(c)
	tx = database.Client(ctx).Save(&entry)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	ginx.SetResp(c, http.StatusNoContent, nil)
}

// DestroyEntry ...
//
//	@Summary	删除条目
//	@Tags		crud
//	@Param		id	path	int	true	"条目 ID"
//	@Success	204	"No Content"
//	@Router		/api/entries/{id} [delete]
func DestroyEntry(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&model.Entry{})
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx.SetResp(c, http.StatusNoContent, nil)
}
