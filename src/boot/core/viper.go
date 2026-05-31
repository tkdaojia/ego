package core

import (
	"ego/src/boot/global"
	"fmt"
	"github.com/spf13/viper"
)

func Viper(configDir string) *viper.Viper {
	v := viper.New()
	// 如果传入的目录为空，默认使用当前同级目录
	if configDir == "" {
		configDir = "."
	}
	// 设置配置文件名和类型
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath(configDir)

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file : %s \n", err))
	}

	if err := v.Unmarshal(&global.C_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
