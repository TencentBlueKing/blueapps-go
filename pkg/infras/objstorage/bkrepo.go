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

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/config"
)

const (
	// ObjectAlreadyExistsErrorCode 对象已存在错误码
	ObjectAlreadyExistsErrorCode = 250107
	// DirAlreadyExistsErrorCode 目录已存在
	DirAlreadyExistsErrorCode = 251012
)

// BkRepo（蓝鲸制品库）是蓝鲸提供的对象存储服务，支持通过 HTTP API 上传/下载图片、文件等数据
// api doc ref: https://github.com/TencentBlueKing/bk-repo/tree/v1.5.1-rc.10/docs/apidoc

// BkGenericRepoClient 蓝盾通用制品仓库客户端
type BkGenericRepoClient struct {
	cfg    *config.BkRepoConfig
	client *resty.Client
}

// ListDir 分页列出制品库目录（path）下的文件（path 需要是某个目录）
func (c *BkGenericRepoClient) ListDir(ctx context.Context, path string, curPage, pageSize int) (*DirInfo, error) {
	url := fmt.Sprintf("/repository/api/node/page/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)
	params := map[string]string{
		"pageNumber":    strconv.Itoa(curPage),
		"pageSize":      strconv.Itoa(pageSize),
		"includeFolder": "true",
		// 指定排序字段 -> 文件类型排前面
		"sort":         "true",
		"sortProperty": "folder",
	}

	var respData ListDirResp
	resp, err := c.client.R().SetContext(ctx).SetResult(&respData).SetQueryParams(params).Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("list dir objects failed, status code: %d", resp.StatusCode())
	}

	if respData.Code != 0 {
		return nil, errors.Errorf(
			"list dir objects failed, traceID: %s, code: %d, message: %s",
			respData.TraceID, respData.Code, respData.Message,
		)
	}
	return &respData.Data, nil
}

// CreateDir 在制品库中创建目录（path）
func (c *BkGenericRepoClient) CreateDir(ctx context.Context, path string) error {
	url := fmt.Sprintf("/repository/api/node/mkdir/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)

	var respData CreateDirResp
	resp, err := c.client.R().SetContext(ctx).SetResult(&respData).Post(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("create dir failed, status code: %d", resp.StatusCode())
	}

	if respData.Code != 0 {
		return errors.Errorf(
			"create dir failed, traceID: %s, code: %d, message: %s",
			respData.TraceID, respData.Code, respData.Message,
		)
	}
	return nil
}

// DeleteDir 删除制品库中的目录（path）
func (c *BkGenericRepoClient) DeleteDir(ctx context.Context, path string) error {
	url := fmt.Sprintf("/repository/api/node/delete/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)

	var respData DeleteDirResp
	resp, err := c.client.R().SetContext(ctx).SetResult(&respData).Delete(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("delete dir failed, status code: %d %s", resp.StatusCode(), resp.String())
	}

	if respData.Code != 0 {
		return errors.Errorf(
			"delete dir failed, traceID: %s, code: %d, message: %s",
			respData.TraceID, respData.Code, respData.Message,
		)
	}
	return nil
}

// UploadFile 上传文件 file 到制品库 path 路径上，允许通过 allowOverwrite 参数决定是否覆盖已有文件
func (c *BkGenericRepoClient) UploadFile(ctx context.Context, file io.Reader, path string, allowOverwrite bool) error {
	url := fmt.Sprintf("/generic/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)
	var respData UploadFileResp

	// 注：不允许使用 SetFileReader，会导致上传数据附带额外的元数据
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(file).
		SetHeader("Content-Type", "application/octet-stream").
		SetHeader("X-BKREPO-OVERWRITE", strconv.FormatBool(allowOverwrite)).
		SetResult(&respData).
		Put(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("upload file failed, status code: %d", resp.StatusCode())
	}
	if respData.Code != 0 {
		baseErrMsg := "upload file failed"
		if respData.Code == ObjectAlreadyExistsErrorCode || respData.Code == DirAlreadyExistsErrorCode {
			baseErrMsg = "object already exists"
		}
		return errors.Errorf("%s, traceID: %s, code: %d, message: %s",
			baseErrMsg, respData.TraceID, respData.Code, respData.Message,
		)
	}
	return nil
}

// DownloadFile 下载文件
func (c *BkGenericRepoClient) DownloadFile(ctx context.Context, path string) (io.ReadCloser, error) {
	url := fmt.Sprintf("/generic/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("download", "true").
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		_ = resp.RawBody().Close()
		return nil, errors.Errorf("download file failed, status code: %d", resp.StatusCode())
	}
	return resp.RawBody(), nil
}

// DeleteFile 删除文件
func (c *BkGenericRepoClient) DeleteFile(ctx context.Context, path string) error {
	url := fmt.Sprintf("/generic/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)

	resp, err := c.client.R().SetContext(ctx).Delete(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.Errorf("delete file failed, status code: %d", resp.StatusCode())
	}
	return nil
}

// GetFileMetadata 获取文件元数据
func (c *BkGenericRepoClient) GetFileMetadata(ctx context.Context, path string) (map[string]string, error) {
	url := fmt.Sprintf("/generic/%s/%s/%s", c.cfg.Project, c.cfg.Bucket, path)

	resp, err := c.client.R().SetContext(ctx).Head(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("get file metadata failed, status code: %d", resp.StatusCode())
	}

	// 读取元数据（只取第一个值）
	metadata := map[string]string{}
	for k, v := range resp.Header() {
		metadata[k] = v[0]
	}
	return metadata, nil
}

// GenPreSignedUrl 生成预签名 URL（只允许下载）
func (c *BkGenericRepoClient) GenPreSignedUrl(
	ctx context.Context, path string, expireSeconds int,
) (*PreSignedUrlData, error) {
	url := "/generic/temporary/url/create"

	body := map[string]any{
		"projectId":     c.cfg.Project,
		"repoName":      c.cfg.Bucket,
		"fullPathSet":   []string{path},
		"expireSeconds": expireSeconds,
		"type":          "DOWNLOAD",
	}
	var respData GenPreSignedUrlResp

	resp, err := c.client.R().SetContext(ctx).SetBody(body).SetResult(&respData).Post(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.Errorf("gen pre-signed url failed, status code: %d", resp.StatusCode())
	}
	if respData.Code != 0 {
		return nil, errors.Errorf(
			"gen pre-signed url failed, traceID: %s, code: %d, message: %s",
			respData.TraceID, respData.Code, respData.Message,
		)
	}
	return lo.ToPtr(respData.Data[0]), nil
}
