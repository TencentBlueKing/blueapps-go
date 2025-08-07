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

package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Masterminds/sprig/v3"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/account"
	"bk.tencent.com/blueapps-go/pkg/common"
	"bk.tencent.com/blueapps-go/pkg/i18n"
	"bk.tencent.com/blueapps-go/pkg/middleware"
	"bk.tencent.com/blueapps-go/pkg/utils/ginx"
)

func TestUserAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("unittest-session", store))
	router.Use(middleware.UserAuth([]account.AuthBackend{
		account.NewStubAuthBackend(account.AllowSpecialTokenStubType),
	}))
	fm := sprig.FuncMap()
	fm["i18n"] = i18n.TranslateWithLang
	router.SetFuncMap(fm)
	router.LoadHTMLGlob("../../templates/web/*")

	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get(common.UserIDCtxKey)
		if exists {
			c.String(http.StatusOK, "user id: %s", userID)
		} else {
			c.String(http.StatusInternalServerError, "user id not found")
		}
	})

	t.Run("No user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: "InvalidToken"})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Valid user token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: account.SpecialTokenForTest})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user id: admin", rr.Body.String())
	})
}

func TestUserAuthWithSession(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("unittest-session", store))
	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(common.UserTokenKey, account.SpecialTokenForTest)
		session.Set(common.UserIDKey, "session-admin")
		_ = session.Save()
		c.Next()
	})
	router.Use(middleware.UserAuth([]account.AuthBackend{
		account.NewStubAuthBackend(account.AllowSpecialTokenStubType),
	}))

	router.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get(common.UserIDCtxKey)
		if exists {
			c.String(http.StatusOK, "user id: %s", userID)
		} else {
			c.String(http.StatusInternalServerError, "user id not found")
		}
	})

	t.Run("Use username in session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: account.SpecialTokenForTest})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "user id: session-admin", rr.Body.String())
	})
}

func TestUserAuthWithoutAuthBackend(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.Use(middleware.UserAuth([]account.AuthBackend{}))

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	t.Run("AuthBackend required", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Equal(t, "Auth backend required!", rr.Body.String())
	})
}

func TestUserAuthMultipleAuthBackends(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("unittest-session", store))
	router.Use(middleware.UserAuth([]account.AuthBackend{
		account.NewStubAuthBackend(account.ForbidAllStubType),
		account.NewStubAuthBackend(account.AllowSpecialTokenStubType),
		account.NewStubAuthBackend(account.AllowAllStubType),
	}))

	router.GET("/test", func(c *gin.Context) {
		session := sessions.Default(c)
		authSource := session.Get(common.UserAuthSourceKey)
		c.String(http.StatusOK, fmt.Sprintf("%s %s", ginx.GetUserID(c), authSource))
	})

	t.Run("Hit second AuthBackend", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: account.SpecialTokenForTest})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "admin allow_special_token", rr.Body.String())
	})

	t.Run("Hit third AuthBackend", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: common.UserTokenKey, Value: "hello world"})
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "anonymous allow_all", rr.Body.String())
	})
}
