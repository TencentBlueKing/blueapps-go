# 数据库迁移（Migration）指南

Migration 可用于管理数据库表结构变化，它允许开发者定义数据库的每次变化并应用到数据库中，以确保表结构与应用代码数据结构保持一致。

在使用 ORM 工具时，Migration 变得尤为重要，因为它帮助我们以编程方式管理数据库变化，而且也方便我们对数据库进行版本控制。

## 什么是 Migration ？

由于 Golang 开发框架中默认使用 [gorm](https://github.com/go-gorm/gorm) 作为基座，因此我们通过简单封装 [gormigrate](https://github.com/go-gormigrate/gormigrate) 来帮助我们更好地管理数据库版本变化。

在 `pkg/migration` 目录下，每个 Go 文件代表一次数据库版本变化，开发者可以通过 `migrate` 命令来将数据库应用到指定的版本。

## 如何编写 Migration ？

开发者可以通过执行 `make-migration` 命令来创建新的 Migration 文件：

```shell
go run main.go make-migration
```

该命令将会通过模板渲染出一个空的 Migration 文件，例如：

```go
package migration

import (
    "github.com/go-gormigrate/gormigrate/v2"
    "gorm.io/gorm"

    "bk.tencent.com/blueapps-go/pkg/infras/database"
    "bk.tencent.com/blueapps-go/pkg/model"
)

func init() {
	// Do Not Edit Migration ID!
	migrationID := "20241022_123456"

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
```

在创建新的 Migration 文件后，开发者需要手动实现 `Migrate` 和 `Rollback` 两个方法。

需要注意的是：`Migrate` 和 `Rollback` 两个方法必须是对等的，即调用 `Migrate` 后再调用 `Rollback` 数据库结构不会有变化。

新建表的简单示例如下：

```go
Migrate: func(tx *gorm.DB) error {
    return tx.AutoMigrate(&model.Task{}, &model.PeriodicTask{})
}

Rollback: func(tx *gorm.DB) error {
    return tx.Migrator().DropTable(&model.Task{}, &model.PeriodicTask{})
}
```

在这个例子中，我们使用 [AutoMigrate](https://gorm.io/docs/migration.html#Auto-Migration) 来自动创建 Task，PeriodicTask 表，并且使用 DropTable 来删除这两个表以确保对等。

注意(重要)：
1. GORM 的 `AutoMigrate` 可以自动创建、更新表，字段和索引，但不会删除未使用的列 
2. `AutoMigrate` 并非总能感知到字段类型的变化，如 `string` -> `sql.NullString` 的变更并不会应用到数据库中 
3. 针对 2 这种场景，需要手动调用 `db.Migrator().AlterColumn(&Model{}, "Field")` 强制更新字段
4. 更多参考：[GORM Migration](https://gorm.io/docs/migration.html)

## 如何应用 Migration ？

开发者可以通过执行 `migrate` 命令来将数据库应用到指定的版本。

```shell
# 更新数据库版本
go run main.go migrate --conf=configs/config.yaml

# 注：`migrate` 命令还支持通过指定 `migration` 参数来迁移 / 回滚数据库到指定版本
go run main.go migrate --conf=configs/config.yaml --migration=20241022_105518
```

## 目前的设计不满足项目需求？

正如 Gormigrate [Readme](https://github.com/go-gormigrate/gormigrate?tab=readme-ov-file#who-is-gormigrate-for) 所说，其主要的使用场景是小型项目（如 SaaS），它简单但足够可靠，在绝大多数场景下是够用的。

如果需要更强大的版本管理功能（如：基于其他 ORM / 使用 SQL 文件管理版本 / NoSQL 数据库...），可以了解下  [golang-migrate/migrate](https://github.com/golang-migrate/migrate) 或 [pressly/goose](https://github.com/pressly/goose)。
