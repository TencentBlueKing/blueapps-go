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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"

	"bk.tencent.com/blueapps-go/pkg/apis/objstorage/serializer"
	"bk.tencent.com/blueapps-go/pkg/infras/objstorage"
	log "bk.tencent.com/blueapps-go/pkg/logging"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// CreateDir ...
//
//	@Summary	创建目录
//	@Tags		object-storage
//	@Param		body	body		serializer.CreateDirRequest	true	"创建目录请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/obj-storage/dirs [post]
func CreateDir(c *gin.Context) {
	var req serializer.CreateDirRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	if err := objstorage.NewClient(ctx).CreateDir(ctx, req.DirPath); err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	ginx.SetResp(c, http.StatusCreated, nil)

	log.Infof(ctx, "user %s create dir %s", ginx.GetUserID(c), req.DirPath)
}

// DeleteDir ...
//
//	@Summary	删除目录
//	@Tags		object-storage
//	@Param		query	query	serializer.DeleteDirRequest	true	"删除目录请求体"
//	@Success	204		"No Content"
//	@Router		/api/obj-storage/dirs [delete]
func DeleteDir(c *gin.Context) {
	var req serializer.DeleteDirRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	if err := objstorage.NewClient(ctx).DeleteDir(ctx, req.DirPath); err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	ginx.SetResp(c, http.StatusNoContent, nil)

	log.Infof(ctx, "user %s delete dir %s", ginx.GetUserID(c), req.DirPath)
}

// ListObjects ...
//
//	@Summary	获取已上传对象列表
//	@Tags		object-storage
//	@Param		query	query		serializer.ListObjectsRequest	true	"获取对象列表请求体"
//	@Success	200		{object}	ginx.Response{data=ginx.PaginatedResp{results=[]serializer.CategoryListResponse}}
//	@Router		/api/obj-storage/objects [get]
func ListObjects(c *gin.Context) {
	var req serializer.ListObjectsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	// 获取指定对象存储目录下属对象数据
	ctx := c.Request.Context()
	dirInfo, err := objstorage.NewClient(ctx).ListDir(ctx, req.DirPath, ginx.GetPage(c), ginx.GetLimit(c))
	if err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	respData := []serializer.ListObjectsResponse{}
	for _, r := range dirInfo.Objects {
		respData = append(respData, serializer.ListObjectsResponse{
			Name:      r.Name,
			IsDir:     r.IsDir,
			Size:      r.Size,
			SHA256:    r.SHA256,
			UpdatedAt: r.LastModifiedDate,
		})
	}
	ginx.SetResp(c, http.StatusOK, ginx.NewPaginatedRespData(dirInfo.Total, respData))
}

// UploadObject ...
//
//	@Summary	上传对象
//	@Tags		object-storage
//	@Param		body	body		serializer.UploadObjectRequest	true	"上传对象请求体"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/obj-storage/objects [post]
func UploadObject(c *gin.Context) {
	var req serializer.UploadObjectRequest
	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	rawFile, err := req.File.Open()
	if err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	ctx := c.Request.Context()
	if err = objstorage.NewClient(ctx).UploadFile(
		ctx, rawFile, fmt.Sprintf("%s%s", req.DirPath, req.File.Filename), true,
	); err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	ginx.SetResp(c, http.StatusCreated, nil)

	log.Infof(ctx, "user %s upload object %s%s", ginx.GetUserID(c), req.DirPath, req.File.Filename)
}

// DownloadObject ...
//
//	@Summary	下载对象
//	@Tags		object-storage
//	@Param		query	query	serializer.DownloadObjectRequest	true	"下载对象请求体"
//	@Produces	octet-stream
//	@Success	200	{file}	file	"文件内容"
//	@Router		/api/obj-storage/objects/download [get]
func DownloadObject(c *gin.Context) {
	var req serializer.DownloadObjectRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	path := fmt.Sprintf("%s%s", req.DirPath, req.ObjName)
	// 获取存储对象元数据
	ctx := c.Request.Context()
	cli := objstorage.NewClient(ctx)
	metadata, err := cli.GetFileMetadata(ctx, path)
	if err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	// 获取存储对象数据（io.ReadCloser）
	reader, err := cli.DownloadFile(ctx, path)
	if err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	log.Infof(ctx, "user %s download object %s%s", ginx.GetUserID(c), req.DirPath, req.ObjName)

	c.Header("Content-Type", metadata["Content-Type"])
	c.Header("Content-Disposition", metadata["Content-Disposition"])
	c.Header("Content-Length", metadata["Content-Length"])
	c.DataFromReader(http.StatusOK, cast.ToInt64(metadata["Content-Length"]), metadata["Content-Type"], reader, nil)
}

// DeleteObject ...
//
//	@Summary	删除已上传对象
//	@Tags		object-storage
//	@Param		query	query	serializer.DeleteObjectRequest	true	"删除对象请求体"
//	@Success	204		"No Content"
//	@Router		/api/obj-storage/objects [delete]
func DeleteObject(c *gin.Context) {
	var req serializer.DeleteObjectRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := c.Request.Context()
	path := fmt.Sprintf("%s%s", req.DirPath, req.ObjName)
	if err := objstorage.NewClient(ctx).DeleteFile(ctx, path); err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, err.Error())
		return
	}
	ginx.SetResp(c, http.StatusNoContent, nil)

	log.Infof(ctx, "user %s delete object %s%s", ginx.GetUserID(c), req.DirPath, req.ObjName)
}
