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

package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"bk.tencent.com/blueapps-go/pkg/config"
	"bk.tencent.com/blueapps-go/pkg/infras/database"
	log "bk.tencent.com/blueapps-go/pkg/logging"
)

var migrationTmpl = `
// Package migration stores all database migrations
package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"

	"bk.tencent.com/blueapps-go/pkg/infras/database"
)


func init() {
	// Do Not Edit Migration ID!
	migrationID := "{{ .id }}"

	database.RegisterMigration(&gormigrate.Migration{
		ID: migrationID,
		Migrate: func(tx *gorm.DB) error {
			logApplying(migrationID)

			// TODO implement migrate code
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			logRollingBack(migrationID)

			// TODO implement rollback code
			return nil
		},
	})
}
`

var makeMigrationCmd = &cobra.Command{
	Use:   "make-migration",
	Short: "Generate an empty migration file.",
	Run: func(cmd *cobra.Command, args []string) {
		migrationID := database.GenMigrationID()

		// 文件
		fileName := fmt.Sprintf("%s.go", migrationID)
		filePath := path.Join(config.BaseDir, "pkg/migration", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("failed to create migration file with path: %s, err: %s", filePath, err)
		}
		defer file.Close()

		// 模板
		tmpl, err := template.New("migration").
			Parse(strings.TrimLeft(migrationTmpl, "\n"))
		if err != nil {
			log.Fatal("failed to initialize migration template")
		}
		if err = tmpl.Execute(file, map[string]string{"id": migrationID}); err != nil {
			log.Fatal("failed to render migration file from template")
		}

		log.Infof(
			context.Background(),
			"migration file %s generated, you must edit it and "+
				"implement the migration logic and then run `migrate` to apply",
			fileName,
		)
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
}
