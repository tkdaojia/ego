package global

import (
	model "ego/src/model/basic"
	"ego/src/model/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	C_VP     *viper.Viper
	C_DB     *gorm.DB
	C_REDIS  model.C_RedisCache
	C_CONFIG config.Server
	C_LOG    *zap.Logger
)
