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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyKMSEnvironment(t *testing.T) {
	t.Setenv("KMS_TEST_PASSWORD", "old-value")
	t.Setenv("KMS_TEST_EMPTY", "old-value")

	err := applyKMSEnvironment(`{"schema_version":1,"config_version_id":"v1","env":{"KMS_TEST_PASSWORD":"kms-value","KMS_TEST_EMPTY":""}}`)
	assert.NoError(t, err)
	assert.Equal(t, "kms-value", os.Getenv("KMS_TEST_PASSWORD"))
	assert.Equal(t, "", os.Getenv("KMS_TEST_EMPTY"))
}

func TestApplyKMSEnvironmentRejectsInvalidPayloadBeforeChangingEnvironment(t *testing.T) {
	for _, payload := range []string{
		`{"schema_version":2,"env":{"KMS_TEST_PASSWORD":"new-value"}}`,
		`{"schema_version":1,"env":{"KMS_TEST_PASSWORD":"new-value","INVALID":42}}`,
		`{"schema_version":1,"env":{"KMS_TEST_PASSWORD":"new-value","INVALID=NAME":"value"}}`,
	} {
		t.Run(payload, func(t *testing.T) {
			t.Setenv("KMS_TEST_PASSWORD", "old-value")
			assert.Error(t, applyKMSEnvironment(payload))
			assert.Equal(t, "old-value", os.Getenv("KMS_TEST_PASSWORD"))
		})
	}
}
