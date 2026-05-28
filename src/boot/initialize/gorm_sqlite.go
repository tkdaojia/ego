package initialize

import (
	"ego/src/boot/global"
	"ego/src/boot/initialize/internal"
	"ego/src/model/config"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"time"
)

// GormSqlite 初始化Sqlite数据库
func GormSqlite() *gorm.DB {
	s := global.C_CONFIG.Sqlite
	return initSqliteDatabase(s)
}

// initSqliteDatabase 初始化Sqlite数据库辅助函数
func initSqliteDatabase(s config.Sqlite) *gorm.DB {
	if s.Dbname == "" {
		return nil
	}

	// 数据库配置
	general := s.GeneralDB
	if db, err := gorm.Open(sqlite.Open(s.Dsn()), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(s.ConnMaxLifetime) * time.Second)
		return db
	}
}
