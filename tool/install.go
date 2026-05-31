package main

import (
	"ego/src/boot/core"
	"ego/src/boot/global"
	"ego/src/boot/initialize"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
)

func main() {
	global.C_VP = core.Viper("../")
	db := initialize.Gorm()

	if global.C_CONFIG.System.DbType == "mysql" {
		if !strings.Contains(global.C_CONFIG.Mysql.Config, "&multiStatements=true") {
			panic("Mysql导入,请先修改配置文件Mysql-config:charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true")
		}
	}

	sqlFilePath := "db.sql"
	sqlBytes, err := os.ReadFile(sqlFilePath)
	if err != nil {
		panic("读取 SQL 文件失败")
	}

	sqlStr := string(sqlBytes)
	if strings.TrimSpace(sqlStr) == "" {
		panic("读取 SQL 文件失败")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if execErr := tx.Exec(sqlStr).Error; execErr != nil {
			return execErr
		}
		return nil
	})

	if err != nil {
		panic("执行 SQL 导入失败" + err.Error())
	}

	fmt.Println("✅  成功导入数据")
	fmt.Println("请改回Mysql-config:charset=utf8mb4&parseTime=True&loc=Local")
}
