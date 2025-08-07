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

package crypto_test

import (
	"crypto/rand"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"bk.tencent.com/blueapps-go/pkg/utils/crypto"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		encFunc func(key, data string) (string, error)
		decFunc func(key, data string) (string, error)
		keySize int
	}{
		{
			name:    "AES",
			encFunc: crypto.AESEncrypt,
			decFunc: crypto.AESDecrypt,
			keySize: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plaintext := "hello world"
			key := make([]byte, tt.keySize)
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				t.Fatal(err)
			}
			// 第一次加密
			ciphertext1, err := tt.encFunc(string(key), plaintext)
			assert.Nil(t, err)

			// 第二次加密
			ciphertext2, err := tt.encFunc(string(key), plaintext)
			assert.Nil(t, err)

			// 确保两次加密后的结果不一致
			assert.NotEqual(t, ciphertext1, ciphertext2)

			// 解密第一次加密的结果
			decrypted1, err := tt.decFunc(string(key), ciphertext1)
			assert.Nil(t, err)
			assert.Equal(t, plaintext, decrypted1)

			// 解密第二次加密的结果
			decrypted2, err := tt.decFunc(string(key), ciphertext2)
			assert.Nil(t, err)
			assert.Equal(t, plaintext, decrypted2)
		})
	}
}
