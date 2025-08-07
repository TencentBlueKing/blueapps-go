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

	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// SendEmailRequest ...
type SendEmailRequest struct {
	Receiver string `json:"receiver" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

// Validate ...
func (r *SendEmailRequest) Validate(c *gin.Context) error {
	if r.Receiver != ginx.GetUserID(c) {
		return errors.New(i18n.T(c.Request.Context(), "can only send emails to yourself currently"))
	}
	return nil
}
