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

package objstorage

// Object 蓝盾制品库文件数据
type Object struct {
	Path             string            `json:"path"`
	Name             string            `json:"name"`
	FullPath         string            `json:"fullPath"`
	IsDir            bool              `json:"folder"`
	Size             int64             `json:"size"`
	SHA256           string            `json:"sha256,omitempty"`
	MD5              string            `json:"md5,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	CreatedBy        string            `json:"createdBy,omitempty"`
	CreatedDate      string            `json:"createdDate,omitempty"`
	LastModifiedBy   string            `json:"lastModifiedBy,omitempty"`
	LastModifiedDate string            `json:"lastModifiedDate,omitempty"`
}

// DirInfo 蓝盾制品库目录文件数据
type DirInfo struct {
	PageNumber int      `json:"pageNumber"`
	PageSize   int      `json:"pageSize"`
	TotalPages int64    `json:"totalPages"`
	Total      int64    `json:"totalRecords"`
	Objects    []Object `json:"records"`
}

// ListDirResp 分页查询蓝盾制品库目录下文件返回结果
type ListDirResp struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    DirInfo `json:"data"`
	TraceID string  `json:"traceId"`
}

// CreateDirResp 创建目录返回结果
type CreateDirResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"traceId"`
}

// DeleteDirResp 删除目录返回结果
type DeleteDirResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"traceId"`
}

// UploadFileResp 上传文件返回结果
type UploadFileResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Object `json:"data"`
	TraceID string `json:"traceId"`
}

// PreSignedUrlData 预签名 URL 数据
type PreSignedUrlData struct {
	FullPath   string `json:"fullPath"`
	Url        string `json:"url"`
	ExpireDate string `json:"expireDate"`
}

// GenPreSignedUrlResp 生成预签名 URL（只允许下载）返回结果
type GenPreSignedUrlResp struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Data    []PreSignedUrlData `json:"data"`
	TraceID string             `json:"traceId"`
}
