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

package envx

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/TencentBlueKing/bk-kms-sdk/go/decrypt"
	"github.com/pkg/errors"
)

// Both files are mounted by the platform.
const (
	kmsEnvelopePath   = "/etc/secrets/envelope"
	kmsPrivateKeyPath = "/etc/secrets/private-key"
)

// decryptKMSEnvironment restores the sensitive variables frozen by PaaS for
// this deployment before config loading.
func decryptKMSEnvironment() error {
	envelope, err := os.ReadFile(kmsEnvelopePath)
	if err != nil {
		return errors.Wrap(err, "read KMS credential envelope")
	}
	privateKey, err := os.ReadFile(kmsPrivateKeyPath)
	if err != nil {
		return errors.Wrap(err, "read KMS private key")
	}

	plaintext, err := decrypt.Decrypt(strings.TrimSpace(string(envelope)), strings.TrimSpace(string(privateKey)))
	if err != nil {
		// Do not expose SDK errors that may include credential contents.
		return errors.New("failed to decrypt KMS credential envelope")
	}

	return applyKMSEnvironment(plaintext)
}

func applyKMSEnvironment(plaintext string) error {
	var payload struct {
		SchemaVersion int            `json:"schema_version"`
		Env           map[string]any `json:"env"`
	}
	if err := json.Unmarshal([]byte(plaintext), &payload); err != nil {
		return errors.New("invalid KMS credential payload")
	}
	if payload.SchemaVersion != 1 || len(payload.Env) == 0 {
		return errors.New("KMS credential payload requires schema_version 1 and non-empty env")
	}

	// Validate all variables before changing the current process environment.
	for name, raw := range payload.Env {
		value, ok := raw.(string)
		if !ok || name == "" || strings.ContainsAny(name, "=\x00") || strings.ContainsRune(value, '\x00') {
			return errors.New("KMS env must contain valid string environment variable names and values")
		}
	}
	for name, raw := range payload.Env {
		if err := os.Setenv(name, raw.(string)); err != nil {
			return errors.Wrapf(err, "set KMS environment variable %s", name)
		}
	}
	return nil
}
