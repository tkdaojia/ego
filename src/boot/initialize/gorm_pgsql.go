package initialize

import (
	"ego/src/boot/global"
	"ego/src/boot/initialize/internal"
	"ego/src/model/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// GormPgSql 初始化 Postgresql 数据库
func GormPgSql() *gorm.DB {
	p := global.C_CONFIG.Pgsql
	return initPgSqlDatabase(p)
}

// initPgSqlDatabase 初始化 Postgresql 数据库的辅助函数
func initPgSqlDatabase(p config.Pgsql) *gorm.DB {
	if p.Dbname == "" {
		return nil
	}
	pgsqlConfig := postgres.Config{
		DSN:                  p.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}
	// 数据库配置
	general := p.GeneralDB
	if db, err := gorm.Open(postgres.New(pgsqlConfig), internal.Gorm.Config(general)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(p.ConnMaxLifetime) * time.Second)
		registerGlobalOmitZeroId(db)
		return db
	}
}
