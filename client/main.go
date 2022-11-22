package main

import (
	"Grpc/proto"
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	//#region consul连接测试
	conn, _ := grpc.Dial(
		// consul://127.0.0.1:8500 consul地址
		// UserService 拉取的服务名
		// 底层就是利用grpc-consul-resolver将参数解析成HTTP请求获取对应的服务
		"consul://127.0.0.1:8500/UserService",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	defer conn.Close()

	userSer := common.NewUserServiceClient(conn) // 创建一个UserService客户端

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用GetUser
	resp1, _ := userSer.GetUser(ctx, &common.GetUserRequest{Id: 100, Name: "diven zhong"})
	fmt.Println(resp1.Code)
	fmt.Println(resp1.Msg)
	fmt.Println(resp1.Data)

	// 调用GetNames
	resp2, _ := userSer.GetNames(ctx, &common.GetNamesRequest{})
	fmt.Println(resp2.Data[0])
	//#endregion

	//#region grpc连接测试
	conn_grpc, _ := grpc.Dial("localhost:8080", grpc.WithInsecure()) //建立grpc链接
	greeterSer := common.NewGreeterClient(conn_grpc)                 // 创建一个Greeter客户端
	ctx_grpc, _ := context.WithTimeout(context.Background(), time.Second)
	// 调用Greeter客户端SayHello
	r1, err1 := greeterSer.SayHello(ctx_grpc, &common.HelloRequest{Name: "diven zhong"})
	if err1 != nil {
		log.Fatalf("could not greet: %v", err1)
	}
	log.Printf("Greeting: %s", r1.GetMessage())

	// 调用Greeter客户端SayHelloAgain
	r2, err2 := greeterSer.SayHelloAgain(ctx_grpc, &common.HelloRequest{Name: "diven zhong"})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Printf("Greeting: %s", r2.GetMessage())
	//#endregion
}
