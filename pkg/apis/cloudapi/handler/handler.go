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

	"github.com/gin-gonic/gin"

	"bk.tencent.com/blueapps-go/pkg/apis/cloudapi/serializer"
	"bk.tencent.com/blueapps-go/pkg/infras/cloudapi/cmsi"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// SendEmail ...
//
//	@Summary	发送邮件
//	@Tags		cloud-api
//	@Param		body	body		serializer.SendEmailRequest	true	"邮件配置"
//	@Success	201		{object}	ginx.Response{data=nil}
//	@Router		/api/emails [post]
func SendEmail(c *gin.Context) {
	var req serializer.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(c); err != nil {
		ginx.SetErrResp(c, http.StatusBadRequest, err.Error())
		return
	}

	client, err := cmsi.New()
	if err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, "cmsi client init failed: "+err.Error())
		return
	}
	if _, err = client.SendMail(c.Request.Context(), req.Receiver, req.Title, req.Content); err != nil {
		ginx.SetErrResp(c, http.StatusInternalServerError, "call cmsi api failed: "+err.Error())
		return
	}

	ginx.SetResp(c, http.StatusCreated, nil)
}
