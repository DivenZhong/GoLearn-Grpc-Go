package main

import (
	"Grpc/proto"
	"context"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
	"net"
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
	common.UserServiceServer
	common.GreeterServer
}

// GetUser
func (s *server) GetUser(ctx context.Context, in *common.GetUserRequest) (*common.GetUserResponse, error) {

	log.Printf("Received : %v", in.String())
	var images []*common.ImageObj

	for i := 0; i < 5; i++ {
		a := int64(i)
		images = append(images, &common.ImageObj{Id: a, Index: a})
	}

	var data = []*common.UserObj{&common.UserObj{Id: 1, Name: "weiwei1", Images: images}, &common.UserObj{Id: 2, Name: "weiwei2", Images: images}}
	resp := &common.GetUserResponse{Code: 1, Msg: "success", Data: data}
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
	common.RegisterUserServiceServer(s, &server{})   // 注册服务
	common.RegisterGreeterServer(s, &server{})       // 注册服务
	log.Printf("now listen: %v", "localhost:8080")   // 启动监听
	registerServerConsul()
	_ = s.Serve(listen)
}

// consul 服务注册
func registerServerConsul() {
	// 创建consul客户端
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Println("consul client error : ", err)
	}

	registration := new(consulapi.AgentServiceRegistration)
	registration.Address = "127.0.0.1"             // 服务 IP
	registration.Port = 8080                       // 服务端口
	registration.ID = "20221122"                   // 服务节点的名称
	registration.Name = "UserService"              // 服务名称
	registration.Tags = []string{"UserService-v1"} // tag，可以为空

	// 服务注册
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register server consul error : ", err)
	}
}
