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
	"gorm.io/gorm"

	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	"bk.tencent.com/blueapps-go/pkg/model"
)

// EntryListRequest List Entries API 输入结构
type EntryListRequest struct {
	CategoryID int64  `form:"categoryID" binding:"omitempty,gt=0"`
	Keyword    string `form:"keyword" binding:"omitempty"`
}

// EntryListResponse List Entries API 返回结构
type EntryListResponse struct {
	CategoryID   int64  `json:"categoryID"`
	CategoryName string `json:"categoryName"`

	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Price     float32 `json:"price"`
	Updater   string  `json:"updater"`
	UpdatedAt string  `json:"updatedAt"`
}

// EntryCreateRequest Create Entry API 输入结构
type EntryCreateRequest struct {
	CategoryID int64   `json:"categoryID" binding:"required,gt=0"`
	Name       string  `json:"name" binding:"required,min=1,max=32"`
	Desc       string  `json:"desc" binding:"omitempty"`
	Price      float32 `json:"price" binding:"required,gt=0"`
}

// Validate ...
func (req *EntryCreateRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	tx := database.Client(ctx).Where("name = ?", req.Name).First(&model.Entry{})
	if tx.Error == nil {
		return errors.Errorf(i18n.T(ctx, "entry name `%s` already used"), req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}

// EntryCreateResponse Create Entry API 输出结构
type EntryCreateResponse struct {
	ID int64 `json:"id"`
}

// EntryRetrieveResponse Retrieve Entry API 返回结构
type EntryRetrieveResponse struct {
	CategoryID   int64  `json:"categoryID"`
	CategoryName string `json:"categoryName"`

	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	Price     float32 `json:"price"`
	Creator   string  `json:"creator"`
	Updater   string  `json:"updater"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// EntryUpdateRequest Update Entry API 输入结构
type EntryUpdateRequest struct {
	CategoryID int64 `json:"categoryID"`

	Name  string  `json:"name" binding:"required,min=1,max=32"`
	Desc  string  `json:"desc" binding:"omitempty"`
	Price float32 `json:"price" binding:"required,gt=0"`
}

// Validate ...
func (req *EntryUpdateRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	tx := database.Client(ctx).
		Not("id = ?", c.Param("id")).
		Where("name = ?", req.Name).
		First(&model.Entry{})
	if tx.Error == nil {
		return errors.Errorf(i18n.T(ctx, "entry name `%s` already used"), req.Name)
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New(tx.Error.Error())
}
