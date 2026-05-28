package backrun

import (
	"context"
	"ego/src/cls"
	model "ego/src/model/basic"
	"ego/src/utils"
	"fmt"
	"github.com/robfig/cron/v3"
)

// 定义一个包级别的全局变量，方便以后在 main 函数关闭服务时，能优雅地关闭定时任务
var CronInstance *cron.Cron

func BackRunQueue() {
	if CronInstance != nil {
		return
	}

	CronInstance = cron.New(cron.WithSeconds())
	_, err := CronInstance.AddFunc("0 */1 * * * *", BackRunQueueTask)
	if err != nil {
		fmt.Println("添加任务失败：", err)
		return
	}
	CronInstance.Start()

}

func BackRunQueueTask() {
	ctx := context.Background()
	db := utils.GetDB(ctx)
	var roles []model.SysRole
	db.Select("name").Find(&roles)
	for _, role := range roles {
		cls.SystemUpRoleRedis(role.Name, ctx)
	}
}
