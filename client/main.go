package main

import (
	"Grpc/proto"
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	clientv3 "go.etcd.io/etcd/client/v3"
	etcdResolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcResolver "google.golang.org/grpc/resolver"
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
	conn_grpc, _ := grpc.Dial("localhost:8800", grpc.WithInsecure()) //建立grpc链接
	greeterSer := common.NewGreeterClient(conn_grpc)                 // 创建一个Greeter客户端
	ctx_grpc, cancel_grpc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel_grpc()
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

	//初始化etcd客户端
	var addr = "127.0.0.1:2379"
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})
	//新建builder，etcd官方实现的Builder对象
	r, _ := etcdResolver.NewBuilder(cli)
	//向grpc注册builder，这样Dial时，就可以按照Scheme查找到此Builder，
	grpcResolver.Register(r)
	//注意：project_test为服务名称，要和服务器注册服务名称相匹配。round_robin表示以轮询方式访问grpc服务
	conn, err := grpc.Dial(r.Scheme()+"://"+addr+"/etcdTest",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)

	if err != nil {
		panic(err)
	}
	//下面每隔1s对服务器调用一次
	etcdClient := common.NewGreeterClient(conn)
	for {
		resp, err_etcd := etcdClient.SayHello(context.Background(), &common.HelloRequest{Name: "etcd test"})
		if err_etcd != nil {
			log.Println(err_etcd)
		} else {
			log.Println(resp)
		}
		resp2, err_etcd := etcdClient.SayHelloAgain(context.Background(), &common.HelloRequest{Name: "etcd test"})
		if err_etcd != nil {
			log.Println(err_etcd)
		} else {
			log.Println(resp2)
		}
		<-time.After(time.Second * 2)
	}

}
