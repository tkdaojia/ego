package cls

import (
	"context"
	model "ego/src/model/basic"
	"ego/src/utils"
	"time"
)

// 更新某个权限组的redis
func SystemUpRoleRedis(name string, c context.Context) error {
	db := utils.GetDB(c)
	var one model.SysRole
	if err := db.Where("name = ?", name).First(&one).Error; err != nil {
		return err
	}
	//默认2小时
	if err := utils.Pack.RedisCache.RedisSet(c, "app:role:"+name, one.Rolelist, 2*time.Hour); err != nil {
		return err
	}

	return nil
}

func MapSetIDSysGroup(ctx context.Context) map[int]interface{} {
	db := utils.GetDB(ctx)
	var any []model.SysGroup
	db.Select("id,gname").Order("ordnum ASC").Find(&any)

	var arr = make(map[int]interface{})
	for _, v := range any {
		tid := v.ID
		arr[tid] = v.Gname
	}
	return arr
}
