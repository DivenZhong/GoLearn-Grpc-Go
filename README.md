# GoLearn-Grpc-Go


`Grpc-Go。[Github地址: https://github.com/DivenZhong/GoLearn-Grpc-Go.git)

用于diven zhong个人学习实践,转载请说明

#protobuf简介
Protocol Buffers(protobuf)：与编程语言无关，与程序运行平台无关的数据序列化协议以及接口定义语言(IDL: interface definition language)。

#Grpc-Go的核心概念
    1.编译器protoc：
	2.Go语言的protobuf插件:
	     编程语言的protobuf插件，搭配protoc编译器，根据.proto文件生成对应编程语言的代码，实现各自语言的protobuf协议。
		 在发送请求和接受响应的时候，完成对应的编码和解码工作，将你即将发送的数据编码成gRPC能够传输的形式，又或者将即将接收到的数据解码为编程语言能够理解的数据格式
	相应的指令:	 
	     protoc --go_out=. common.proto
         protoc --go-grpc_out=. common.proto
#接口开发流程步骤	
    1.定义.proto文件，包括消息体和rpc服务接口定义
    2.使用protoc命令来编译.proto文件，用于生成xx.pb.go和xx_grpc.pb.go文件
    3.在服务端实现rpc里定义的方法
    4.客户端调用rpc方法，获取响应结果	
	
	


#consul
      1.安装运行consul,使得127.0.0.1:8500端口可以正常访问
	  2.启动grpc服务端，完成consul服务注册，把grpc服务端服务IP、端口、服务名称等注册到consul上面
	  3.启动grpc客户端，根据服务名称去consul拿到grpc服务端的服务IP、端口等信息，并建立起连接
	  4.正常走grpc之间调用的流程


