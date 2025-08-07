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
	"github.com/samber/lo"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/infras/objstorage"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

// GetIndexPage 首页
func GetIndexPage(c *gin.Context) {
	renderHTML(c, "index.html", nil)
}

// GetHomePage 主页
func GetHomePage(c *gin.Context) {
	renderHTML(c, "home.html", nil)
}

// GetCRUDPage CRUD 示例页面
func GetCRUDPage(c *gin.Context) {
	renderHTML(c, "crud.html", nil)
}

// GetCachePage 缓存示例页面
func GetCachePage(c *gin.Context) {
	renderHTML(c, "cache.html", nil)
}

// GetCloudAPIPage 云 API 调用示例页面
func GetCloudAPIPage(c *gin.Context) {
	renderHTML(c, "cloud_api.html", nil)
}

// GetAsyncTaskPage 异步任务示例页面
func GetAsyncTaskPage(c *gin.Context) {
	renderHTML(c, "async_task.html", nil)
}

// GetObjStoragePage 对象存储示例页面
func GetObjStoragePage(c *gin.Context) {
	renderHTML(c, "obj_storage.html", gin.H{"objectStorageEnabled": objstorage.IsBkRepoAvailable()})
}

func renderHTML(c *gin.Context, name string, data gin.H) {
	data = lo.Assign(data, gin.H{
		"appID": config.G.Platform.AppID,
		"user":  ginx.GetUserID(c),
		"lang":  ginx.GetLang(c),
	})
	c.HTML(http.StatusOK, name, data)
}
