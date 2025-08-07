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

// Package serializer ...
package serializer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
)

// CategoryListRequest List Categories API 输入结构
type CategoryListRequest struct {
	Keyword string `form:"keyword" binding:"omitempty"`
}

// CategoryListResponse List Categories API 返回结构
type CategoryListResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Updater   string `json:"updater"`
	UpdatedAt string `json:"updatedAt"`
}

// CategoryCreateRequest Create Category API 输入结构
type CategoryCreateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

// Validate ...
func (req *CategoryCreateRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("name = ?", req.Name).First(&model.Category{})
	if tx.Error == nil {
		return errors.Errorf(i18n.T(ctx, "category name `%s` already used"), req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}

// CategoryCreateResponse Create Category API 输出结构
type CategoryCreateResponse struct {
	ID int64 `json:"id"`
}

// CategoryRetrieveResponse Retrieve Category API 返回结构
type CategoryRetrieveResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	Updater   string `json:"updater"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// CategoryUpdateRequest Update Category API 输入结构
type CategoryUpdateRequest struct {
	Name string `json:"name" binding:"required,min=1,max=32"`
}

// Validate ...
func (req *CategoryUpdateRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	tx := database.Client(ctx).
		Not("id = ?", c.Param("id")).
		Where("name = ?", req.Name).
		First(&model.Category{})
	if tx.Error == nil {
		return errors.Errorf(i18n.T(ctx, "category name `%s` already used"), req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}
