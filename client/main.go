package main

import (
	"Grpc/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, _ := grpc.Dial("localhost:8080", grpc.WithInsecure()) //建立链接
	defer conn.Close()

	c := common.NewUserServiceClient(conn) // 创建一个客户端
	c2 := common.NewGreeterClient(conn)    // 创建一个客户端

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用GetUser
	resp1, _ := c.GetUser(ctx, &common.GetUserRequest{Id: 100})
	fmt.Println(resp1.Code)
	fmt.Println(resp1.Msg)
	fmt.Println(resp1.Data)

	// 调用GetNames
	resp2, _ := c.GetNames(ctx, &common.GetNamesRequest{})
	fmt.Println(resp2.Data[0])

	// 调用Greeter客户端SayHelloAgain
	r1, err1 := c2.SayHello(ctx, &common.HelloRequest{Name: "diven zhong"})
	if err1 != nil {
		log.Fatalf("could not greet: %v", err1)
	}
	log.Printf("Greeting: %s", r1.GetMessage())

	// 调用Greeter客户端SayHelloAgain
	r2, err2 := c2.SayHelloAgain(ctx, &common.HelloRequest{Name: "diven zhong"})
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Printf("Greeting: %s", r2.GetMessage())

}
