package core

import (
	"context"
	"ego/src/boot/global"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 1. 定义一个自定义结构体用于实现 PerRPCCredentials 接口
type TokenAuth struct {
	Token string
}

// GetRequestMetadata 会在每次 RPC 请求时被调用，将 Token 添加到 metadata (Header) 中
func (t TokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.Token, // 或者自定义 key，如 "x-api-key": t.Token
	}, nil
}

// RequireTransportSecurity 指示该 Token 是否仅能在 TLS 安全连接传输
func (t TokenAuth) RequireTransportSecurity() bool {
	// 因为你当前使用的是 insecure 连接，这里必须返回 false；
	// 如果生产环境开启了 TLS 证书加密，请改为 true
	return false
}

func Grpc() *grpc.ClientConn {
	config := global.C_CONFIG.Grpc
	if config.Open == false {
		return nil
	}
	auth := TokenAuth{
		Token: config.Token,
	}

	path := fmt.Sprintf("%s:%d", config.Host, config.Addr)
	// 建立与服务端的非加密连接
	conn, err := grpc.NewClient(path,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(auth))
	if err != nil {
		panic(err)
	}
	return conn

}
