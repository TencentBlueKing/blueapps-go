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

// Package envx 提供环境变量相关工具
package envx

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/emmansun/gmsm/sm4"
	"github.com/fernet/fernet-go"
	"github.com/pkg/errors"
)

const (
	paasEncryptSecretKeyEnv = "BKPAAS_ENCRYPT_SECRET_KEY"
	paasEncryptedEnvKeysEnv = "BKPAAS_ENCRYPTED_ENV_KEYS"
	paasCipherPrefix        = "bkpaas_enc$"
	fernetCipherPrefix      = "bkcrypt$"
	sm4CTRCipherPrefix      = "sm4ctr$"
)

// Get 读取环境变量，支持默认值
func Get(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// MustGet 读取环境变量，若不存在则 panic
func MustGet(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	panic(fmt.Sprintf("required environment variable %s unset", key))
}

// GetBool 读取环境变量，并自动转换为 bool 类型
// true: 1, t, T, TRUE, true, True
// false: 0, f, F, FALSE, false, False
// 其他情况返回 false
func GetBool(key string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false
	}

	switch value {
	case "1", "t", "T", "TRUE", "true", "True":
		return true
	default:
		return false
	}
}

// DecryptEncryptedEnvironment restores PaaS-managed sensitive variables before
// config loading. KMS-enabled deployments use mounted credential files; other
// deployments only decrypt listed variables with the bkpaas_enc$ prefix.
func DecryptEncryptedEnvironment() error {
	if GetBool("BKPAAS_KMS_ENABLED") {
		return decryptKMSEnvironment()
	}

	secretKey := os.Getenv(paasEncryptSecretKeyEnv)
	encryptedKeys := os.Getenv(paasEncryptedEnvKeysEnv)
	if secretKey == "" || encryptedKeys == "" {
		return nil
	}

	for _, name := range strings.Split(encryptedKeys, ",") {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		payload, ok := strings.CutPrefix(os.Getenv(name), paasCipherPrefix)
		if !ok {
			continue
		}

		plaintext, err := decryptEncryptedValue(secretKey, payload)
		if err != nil {
			return errors.Wrapf(err, "decrypt environment variable %s", name)
		}
		if err := os.Setenv(name, plaintext); err != nil {
			return errors.Wrapf(err, "set decrypted environment variable %s", name)
		}
	}

	return nil
}

// decryptEncryptedValue selects the cipher by the prefix emitted by the PaaS
// platform. This mirrors blue_krill.encrypt.handler.EncryptHandler.decrypt,
// which preserves unheaded payloads for backwards compatibility.
func decryptEncryptedValue(secretKey, payload string) (string, error) {
	if token, ok := strings.CutPrefix(payload, fernetCipherPrefix); ok {
		return decryptFernet(secretKey, token)
	}
	if token, ok := strings.CutPrefix(payload, sm4CTRCipherPrefix); ok {
		return decryptSM4CTR(secretKey, token)
	}
	return payload, nil
}

// decryptFernet decrypts the bkcrypt$ Fernet payload emitted by the PaaS
// platform's EncryptHandler. Fernet keys are URL-safe base64-encoded 32-byte
// values: the first 16 bytes sign the token and the final 16 bytes encrypt it.
func decryptFernet(secretKey, payload string) (string, error) {
	key, err := fernet.DecodeKey(secretKey)
	if err != nil {
		return "", errors.Wrap(err, "decode Fernet key")
	}

	plaintext := fernet.VerifyAndDecrypt([]byte(payload), 0, []*fernet.Key{key})
	if plaintext == nil {
		return "", errors.New("invalid Fernet token")
	}

	return string(plaintext), nil
}

// decryptSM4CTR decrypts the sm4ctr$ values emitted by EncryptHandler. Its
// payload is standard base64 of a 16-byte IV followed by the ciphertext; the
// platform's 32-character hexadecimal key supplies the first 16 bytes used by
// the Python SDK's SM4 cipher.
func decryptSM4CTR(secretKey, payload string) (string, error) {
	if len(secretKey) < sm4.BlockSize {
		return "", errors.New("invalid PaaS SM4 encryption secret key")
	}

	ciphertextWithIV, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", errors.Wrap(err, "decode SM4-CTR token")
	}
	if len(ciphertextWithIV) < sm4.BlockSize {
		return "", errors.New("invalid SM4-CTR token length")
	}

	block, err := sm4.NewCipher([]byte(secretKey[:sm4.BlockSize]))
	if err != nil {
		return "", errors.Wrap(err, "initialize SM4 cipher")
	}
	plaintext := make([]byte, len(ciphertextWithIV)-sm4.BlockSize)
	cipher.NewCTR(block, ciphertextWithIV[:sm4.BlockSize]).XORKeyStream(plaintext, ciphertextWithIV[sm4.BlockSize:])
	return string(plaintext), nil
}
