package main

import (
	"context"
	"ego/src/boot/core"
	"ego/src/boot/global"
	"ego/src/boot/initialize"
	pb "ego/src/zrgrpc"
	"fmt"
	"time"
)

func main() {
	global.C_VP = core.Viper("")
	global.C_LOG = core.Zap()
	initialize.Redis()
	global.C_DB = initialize.Gorm()

	global.C_GRPC = core.Grpc()

	// 实例化客户端
	client := pb.NewGreeterClient(global.C_GRPC)

	// 设置 1 秒超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resb, err := client.GetUser(ctx, &pb.UserRequest{Mobile: "18773410566"})
	if err != nil {
		return
	}
	fmt.Println(resb)
}
