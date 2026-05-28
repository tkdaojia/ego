package initialize

import (
	"ego/src/boot/global"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.C_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

func registerGlobalOmitZeroId(db *gorm.DB) {
	_ = db.Callback().Create().Before("gorm:create").Register("global_omit_zero_id", func(tx *gorm.DB) {
		if tx.Statement.Schema != nil && tx.Statement.Schema.PrioritizedPrimaryField != nil {
			field := tx.Statement.Schema.PrioritizedPrimaryField

			// 只要主键是 ID，就去检查它的值
			if field.Name == "ID" {
				val, isZero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
				if isZero {
					tx.Statement.Omit(field.DBName) // 值为零（比如0），直接从 SQL 里抹除 id 字段
				} else if num, ok := val.(uint); ok && num == 0 {
					tx.Statement.Omit(field.DBName) // 针对 uint 的二次防线
				}
			}
		}
	})
}
