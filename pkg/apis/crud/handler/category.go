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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/apis/crud/serializer"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// ListCategories ...
//
//	@Summary	获取分类列表
//	@Tags		crud
//	@Success	200	{object}	ginx.Response{data=[]serializer.CategoryListResponse}
//	@Router		/api/categories [get]
func ListCategories(c *gin.Context) {
	var req serializer.CategoryListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	tx := database.Client(c.Request.Context()).Model(&model.Category{})
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		tx = tx.Where("LOWER(name) LIKE ?", keyword).Or("LOWER(updater) LIKE ?", keyword)
	}

	var categories []model.Category
	if err := tx.Find(&categories).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.CategoryListResponse{}
	for _, category := range categories {
		respData = append(respData, serializer.CategoryListResponse{
			ID:        category.ID,
			Name:      category.Name,
			Updater:   category.Updater,
			UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
		})
	}
	ginx.SetResp(c, http.StatusOK, respData)
}

// CreateCategory ...
//
//	@Summary	创建分类
//	@Tags		crud
//	@Param		body	body		serializer.CategoryCreateRequest	true	"创建分类请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/categories [post]
func CreateCategory(c *gin.Context) {
	var req serializer.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	category := model.Category{
		Name: req.Name,
		BaseModel: model.BaseModel{
			Creator: ginx.GetUserID(c),
			Updater: ginx.GetUserID(c),
		},
	}
	if err := database.Client(c.Request.Context()).Create(&category).Error; err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	ginx.SetResp(c, http.StatusCreated, serializer.CategoryCreateResponse{ID: category.ID})
}

// RetrieveCategory ...
//
//	@Summary	获取单个分类
//	@Tags		crud
//	@Param		id	path		int	true	"分类 ID"
//	@Success	200	{object}	ginx.Response{data=serializer.CategoryRetrieveResponse}
//	@Router		/api/categories/{id} [get]
func RetrieveCategory(c *gin.Context) {
	var category model.Category

	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).First(&category)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	respData := serializer.CategoryRetrieveResponse{
		ID:        category.ID,
		Name:      category.Name,
		Creator:   category.Creator,
		Updater:   category.Updater,
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
	}
	ginx.SetResp(c, http.StatusOK, respData)
}

// UpdateCategory ...
//
//	@Summary	更新分类
//	@Tags		crud
//	@Param		id		path	int									true	"分类 ID"
//	@Param		body	body	serializer.CategoryUpdateRequest	true	"更新分类请求体"
//	@Success	204		"No Content"
//	@Router		/api/categories/{id} [put]
func UpdateCategory(c *gin.Context) {
	var req serializer.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	var category model.Category
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("id = ?", c.Param("id")).First(&category)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusNotFound, tx.Error.Error())
		return
	}

	// 更新 DB 模型字段
	category.Name = req.Name
	category.Updater = ginx.GetUserID(c)
	tx = database.Client(ctx).Save(&category)
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	ginx.SetResp(c, http.StatusNoContent, nil)
}

// DestroyCategory ...
//
//	@Summary	删除分类
//	@Tags		crud
//	@Param		id	path	int	true	"分类 ID"
//	@Success	204	"No Content"
//	@Router		/api/categories/{id} [delete]
func DestroyCategory(c *gin.Context) {
	tx := database.Client(c.Request.Context()).Where("id = ?", c.Param("id")).Delete(&model.Category{})
	if tx.Error != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, tx.Error.Error())
		return
	}
	ginx.SetResp(c, http.StatusNoContent, nil)
}
