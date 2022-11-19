package main

import (
	"Grpc/proto"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func (s *server) SayHello(ctx context.Context, in *common.HelloRequest) (*common.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &common.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *common.HelloRequest) (*common.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &common.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

type server struct {
	Name string
	common.UserServiceServer
	common.GreeterServer
}

// GetUser
func (s *server) GetUser(ctx context.Context, in *common.GetUserRequest) (*common.GetUserResponse, error) {
	var images []*common.ImageObj

	for i := 0; i < 5; i++ {
		a := int64(i)
		images = append(images, &common.ImageObj{Id: a, Index: a})
	}

	var data = []*common.UserObj{&common.UserObj{Id: 1, Name: "weiwei1", Images: images}, &common.UserObj{Id: 2, Name: "weiwei2", Images: images}}
	resp := &common.GetUserResponse{Code: 1, Msg: "suc", Data: data}
	return resp, nil
}

// GetNames
func (s *server) GetNames(ctx context.Context, in *common.GetNamesRequest) (*common.GetNamesResponse, error) {
	var names = []string{"weiwei1", "weiwei2", "weiwei3"}
	return &common.GetNamesResponse{Code: 1, Msg: "suc", Data: names}, nil
}

// listen
func main() {
	listen, _ := net.Listen("tcp", "localhost:8080") // 创建监听
	s := grpc.NewServer()                            // 创建grpc服务
	common.RegisterUserServiceServer(s, &server{})   // 注册服务                          // 创建grpc服务
	common.RegisterGreeterServer(s, &server{})       // 注册服务
	log.Printf("now listen: %v", "localhost:8080")   // 启动监听
	_ = s.Serve(listen)
}
