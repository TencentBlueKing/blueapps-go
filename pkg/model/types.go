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

package model

import (
	"database/sql"
	"database/sql/driver"

	"github.com/pkg/errors"
	"gorm.io/gorm/schema"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/utils/crypto"
)

// AESEncryptString 为 Gorm 自定义字段类型，用 AES-GCM 算法加密
type AESEncryptString string

// Scan 解析 driver 提供的数据
func (s *AESEncryptString) Scan(value any) error {
	if value == nil {
		*s = ""
		return nil
	}

	var data string
	switch v := value.(type) {
	case string:
		data = v
	case []byte:
		data = string(v)
	default:
		return errors.New("invalid value type")
	}

	decryptedData, err := crypto.AESDecrypt(config.G.Service.EncryptSecret, data)
	if err != nil {
		return err
	}
	*s = AESEncryptString(decryptedData)
	return nil
}

// Value 提供 driver.Value
func (s AESEncryptString) Value() (driver.Value, error) {
	encryptedData, err := crypto.AESEncrypt(config.G.Service.EncryptSecret, string(s))
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// GormDataType 提供 Gorm 需要的数据类型
func (s *AESEncryptString) GormDataType() string {
	return "varchar(128)"
}

var (
	_ driver.Valuer                = (*AESEncryptString)(nil)
	_ sql.Scanner                  = (*AESEncryptString)(nil)
	_ schema.GormDataTypeInterface = (*AESEncryptString)(nil)
)
