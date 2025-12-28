package main

import (
	"context"
	user_service "dqq/go/basic/grpc/live/idl/service"
	"fmt"
	"log"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func attachKey(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "api_key", "123456")
	return invoker(ctx, method, req, reply, cc, opts...)
}

func main() {
	creds, err := credentials.NewClientTLSFromFile("data/server.crt", "")
	if err != nil {
		panic(err)
	}

	conn, err := grpc.NewClient(
		"localhost:1234",                     //跟证书里的域名保持一致
		grpc.WithTransportCredentials(creds), // TLS数据加密。这会导致第1次grpc调用非常慢
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024),
			grpc.MaxCallSendMsgSize(1024),
		),
		grpc.WithChainUnaryInterceptor(attachKey), // 把API Key放到OutgoingContext里
	)
	if err != nil {
		panic(err)
	}

	const P = 1
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			client := user_service.NewUserClient(conn) // 多路复用
			ctx := context.Background()
			client.Regist(ctx, &user_service.RegistRequest{Name: "张三", Password: "12243"}, grpc.MaxCallRecvMsgSize(1024))
			resp, err := client.Login(ctx, &user_service.LoginRequest{Name: "张三", Password: "12243"})
			if err != nil {
				log.Printf("grpc调用失败:%s", err)
			} else {
				fmt.Println(resp.Status.Code, resp.Status.Messgae)
			}
		}()
	}
	wg.Wait()
}

// go run ./grpc/live/client
