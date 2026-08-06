package utils

import (
	"ego/src/utils/aes"
	"ego/src/utils/auditLog"
	"ego/src/utils/datetime"
	"ego/src/utils/rediscache"
	sm "ego/src/utils/sm4"
	"ego/src/utils/xls"
)

var Pack = new(pack)

type pack struct {
	RedisCache rediscache.RedisCache
	AuditLog   auditLog.AuditLog
	Xls        xls.Xls
	DateTime   datetime.DateTime
	SM4        sm.SM4
	Aes256     aes.Aes256
}
