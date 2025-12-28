package main

import (
	"context"
	"dqq/go/basic/grpc/live/idl/service"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type U struct {
	user_service.UnimplementedUserServer // 匿名成员，“继承”
}

func (u U) Regist(ctx context.Context, request *user_service.RegistRequest) (*user_service.RegistResp, error) {
	fmt.Println(request.Name, request.Password)
	return &user_service.RegistResp{
		Status: &user_service.Status{
			Code:    1,
			Messgae: "注册成功",
		},
		Userid: 432,
	}, nil
}
func (u U) Login(ctx context.Context, request *user_service.LoginRequest) (*user_service.LoginResp, error) {
	fmt.Println(request.Name, request.Password)
	return &user_service.LoginResp{
		Status: &user_service.Status{
			Code:    1,
			Messgae: "登录成功",
		},
	}, nil
}

// 拦截器
func timer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	begin := time.Now()
	resp, err = handler(ctx, req)
	log.Printf("use time %d ms", time.Since(begin).Milliseconds())
	return
}

func fetchKey(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("获取不到context")
	}
	if values, exists := meta["api_key"]; !exists {
		return nil, errors.New("获取不到API Key")
	} else {
		if values[0] != "123456" {
			return nil, errors.New("非法的API Key")
		}
	}
	resp, err = handler(ctx, req)
	return
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	creds, err := credentials.NewServerTLSFromFile("data/server.crt", "data/rsa_private_key.pem") // client用server的公钥加密自己的AES key，此后双方互会传数据之前都用该AES key加密数据
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(timer, fetchKey), // 链式拦截器
		grpc.Creds(creds), // TLS数据加密
	)
	user_service.RegisterUserServer(server, U{})
	err = server.Serve(listener)
	if err != nil {
		panic(err)
	}
}

// go run ./grpc/live/server

/**
生成1024位的RSA私钥：
openssl genrsa -out data/rsa_private_key.pem 1024
根据私钥生成公钥：
openssl rsa -in data/rsa_private_key.pem -pubout -out data/rsa_public_key.pem
pem是一种标准格式，它通常包含页眉和页脚。const

生成自签名证书（证书的作用就是证明server的公钥是什么）：
生成证书  openssl req -x509 -new -nodes -key data/rsa_private_key.pem -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost" -days 3650 -out data/server.crt
*/
