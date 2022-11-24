package main

import (
	"Grpc/proto"
	"context"
	consulapi "github.com/hashicorp/consul/api"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
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

var cli *clientv3.Client
var interval = 5
var grpcAddr = "localhost:8800"

// listen
func main() {
	listen, _ := net.Listen("tcp", "localhost:8800") // 创建监听
	s := grpc.NewServer()                            // 创建grpc服务
	common.RegisterUserServiceServer(s, &server{})   // 注册服务
	common.RegisterGreeterServer(s, &server{})       // 注册服务
	log.Printf("now listen: %v", "localhost:8800")   // 启动监听
	registerServerConsul()
	registerServerEtcd([]string{"127.0.0.1:2379"}, "etcdTest", grpcAddr)
	log.Println("启动Grpc服务器：", grpcAddr)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
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
	registration.Port = 8800                       // 服务端口
	registration.ID = "UserService"                // 服务节点的名称
	registration.Name = "UserService"              // 服务名称
	registration.Tags = []string{"UserService-v1"} // tag，可以为空

	// 服务注册
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Println("register server consul error : ", err)
	}
}

// etcd 注册服务
func registerServerEtcd(etcdAddrs []string, serviceName string, serviceAddr string) error {
	//获取链接
	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   etcdAddrs,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			return err
		}
	}
	//注册续租
	register(serviceName, serviceAddr)
	return nil
}

// etcd服务发现时，底层解析的是一个json串，且包含Addr字段
func getValue(addr string) string {
	return "{\"Addr\":\"" + addr + "\"}"
}

func register(serviceName, serviceAddr string) error {
	//注册服务
	leaseResp, err := cli.Grant(context.Background(), int64(interval+1))
	if err != nil {
		return err
	}
	fullKey := serviceName
	_, err = cli.Put(context.Background(), fullKey, getValue(serviceAddr), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}
	keepAlive(serviceName, serviceAddr, leaseResp)
	return nil
}

// 异步续约
func keepAlive(name string, addr string, leaseResp *clientv3.LeaseGrantResponse) {
	//永久续约，续约成功后，etcd客户端和服务器会保持通讯，通讯成功会写数据到返回的通道中
	//停止进程后，服务器链接不上客户端，相应key租约到期会被服务器自动删除
	c, err := cli.KeepAlive(cli.Ctx(), leaseResp.ID)
	go func() {
		if err == nil {
			for {
				select {
				case _, ok := <-c:
					if !ok { //续约失败
						cli.Revoke(cli.Ctx(), leaseResp.ID)
						register(name, addr)
						return
					}
				}
			}
			defer cli.Revoke(cli.Ctx(), leaseResp.ID)
		}
	}()
}

type ServerInfo struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`    // 地址
	Version string `json:"version"` // 版本
	Weight  int64  `json:"weight"`  // 权重
}
