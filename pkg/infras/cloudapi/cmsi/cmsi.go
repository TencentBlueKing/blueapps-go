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

// Package cmsi 提供邮件，短信等消息发送能力
package cmsi

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TencentBlueKing/bk-apigateway-sdks/core/bkapi"
	"github.com/TencentBlueKing/bk-apigateway-sdks/core/define"

	"bk.tencent.com/blueapps-go/pkg/config"
)

// ApiClient 蓝鲸 CMSI 组件 API Client
type ApiClient struct {
	define.BkApiClient
}

// New 创建 ApiClient
func New() (*ApiClient, error) {
	authorization, _ := json.Marshal(map[string]string{
		"bk_app_code":   config.G.Platform.AppID,
		"bk_app_secret": config.G.Platform.AppSecret,
	})
	client, err := bkapi.NewBkApiClient("cmsi", bkapi.ClientConfig{
		Endpoint: config.G.Platform.BkPlatUrl.BkCompApi,
		ClientOptions: []define.BkApiClientOption{
			bkapi.OptSetRequestHeader("x-bkapi-authorization", string(authorization)),
			bkapi.OptJsonResultProvider(),
			bkapi.OptJsonBodyProvider(),
			bkapi.OptTimeout(60 * time.Second),
		},
	})
	if err != nil {
		return nil, err
	}
	return &ApiClient{client}, nil
}

// GetMsgType 获取支持发送的消息类型
func (c *ApiClient) GetMsgType(ctx context.Context) (map[string]any, error) {
	apiOperation := c.BkApiClient.NewOperation(
		bkapi.OperationConfig{
			Name:   "get_msg_type",
			Method: "POST",
			Path:   "/api/c/compapi/v2/cmsi/get_msg_type/",
		},
	)

	var ret map[string]any
	if _, err := apiOperation.SetContext(ctx).SetResult(&ret).Request(); err != nil {
		return nil, err
	}
	return ret, nil
}

// SendMsg 发送指定类型的消息
// NOTE: cmsi.send_msg 不只支持以下参数，SaaS 开发者可查阅文档按需添加
func (c *ApiClient) SendMsg(ctx context.Context, msgType, receiver, title, content string) (map[string]any, error) {
	apiOperation := c.BkApiClient.NewOperation(
		bkapi.OperationConfig{
			Name:   "send_msg",
			Method: "POST",
			Path:   "/api/c/compapi/v2/cmsi/send_msg/",
		},
		bkapi.OptSetRequestBody(map[string]string{
			"msg_type":           msgType,
			"receiver__username": receiver,
			"title":              title,
			"content":            content,
		}),
	)

	var ret map[string]any
	if _, err := apiOperation.SetContext(ctx).SetResult(&ret).Request(); err != nil {
		return nil, err
	}
	return ret, nil
}

// SendMail 发送邮件
// NOTE: cmsi.send_mail 不只支持以下参数，SaaS 开发者可查阅文档按需添加
func (c *ApiClient) SendMail(ctx context.Context, receiver, title, content string) (map[string]any, error) {
	apiOperation := c.BkApiClient.NewOperation(
		bkapi.OperationConfig{
			Name:   "send_email",
			Method: "POST",
			Path:   "/api/c/compapi/v2/cmsi/send_mail/",
		},
		bkapi.OptSetRequestBody(map[string]string{
			"receiver__username": receiver,
			"title":              title,
			"content":            content,
		}),
	)

	var ret map[string]any
	if _, err := apiOperation.SetContext(ctx).SetResult(&ret).Request(); err != nil {
		return nil, err
	}
	return ret, nil
}
