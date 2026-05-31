package main

import (
	"ego/src/boot/core"
	"ego/src/boot/global"
	"ego/src/boot/initialize"
	"ego/src/utils/backrun"
)

func main() {
	global.C_VP = core.Viper("")
	global.C_LOG = core.Zap()
	initialize.Redis()
	global.C_DB = initialize.Gorm()
	backrun.BackRunQueueTask()
	core.RunServer()
}
