package core

import (
	"ego/src/boot/global"
	"fmt"
	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file : %s \n", err))
	}

	if err := v.Unmarshal(&global.C_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
