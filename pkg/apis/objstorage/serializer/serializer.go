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
	"mime/multipart"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"bk.tencent.com/blueapps-go/pkg/i18n"
)

// DirPathRegex 对象存储目录名称正则
var DirPathRegex = regexp.MustCompile(`^(/[\p{Han}\sA-Za-z0-9_-]+)*/$`)

// ObjNameRegex 对象存储文件名称正则
var ObjNameRegex = regexp.MustCompile(`^[\p{Han}\sA-Za-z0-9._-]+(\.[A-Za-z0-9._-]+)*$`)

// CreateDirRequest ...
type CreateDirRequest struct {
	DirPath string `form:"dirPath" binding:"required"`
}

// Validate ...
func (r *CreateDirRequest) Validate(c *gin.Context) error {
	if matched := DirPathRegex.MatchString(r.DirPath); !matched {
		return errors.Errorf(i18n.T(c.Request.Context(), "invalid dir path %s"), r.DirPath)
	}
	return nil
}

// DeleteDirRequest ...
type DeleteDirRequest = CreateDirRequest

// ListObjectsRequest ...
type ListObjectsRequest = CreateDirRequest

// ListObjectsResponse ...
type ListObjectsResponse struct {
	Name      string `json:"name"`
	IsDir     bool   `json:"isDir"`
	Size      int64  `json:"size"`
	SHA256    string `json:"sha256"`
	UpdatedAt string `json:"updatedAt"`
}

// UploadObjectRequest ...
type UploadObjectRequest struct {
	File    *multipart.FileHeader `form:"file"`
	DirPath string                `form:"dirPath"`
}

// Validate ...
func (r *UploadObjectRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	if matched := DirPathRegex.MatchString(r.DirPath); !matched {
		return errors.Errorf(i18n.T(ctx, "invalid dir path %s"), r.DirPath)
	}
	if r.File == nil {
		return errors.Errorf(i18n.T(ctx, "file is required"))
	}
	if matched := ObjNameRegex.MatchString(r.File.Filename); !matched {
		return errors.Errorf(i18n.T(ctx, "invalid file name %s"), r.File.Filename)
	}
	return nil
}

// DownloadObjectRequest ...
type DownloadObjectRequest struct {
	DirPath string `form:"dirPath"`
	ObjName string `form:"objName"`
}

// Validate ...
func (r *DownloadObjectRequest) Validate(c *gin.Context) error {
	ctx := c.Request.Context()
	if matched := DirPathRegex.MatchString(r.DirPath); !matched {
		return errors.Errorf(i18n.T(ctx, "invalid dir path %s"), r.DirPath)
	}
	if matched := ObjNameRegex.MatchString(r.ObjName); !matched {
		return errors.Errorf(i18n.T(ctx, "invalid file name %s"), r.ObjName)
	}
	return nil
}

// DeleteObjectRequest ...
type DeleteObjectRequest = DownloadObjectRequest
