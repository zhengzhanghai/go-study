package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/handler"
	"shop_srvs/user_srv/initialize"
	"shop_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	fmt.Println("ip:", *IP, "port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserService{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// 服务注册到consul
	RegisterService()

	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc: " + err.Error())
	}
}

func RegisterService() {
	cfg := api.DefaultConfig()
	// consul部署的服务器地址
	cfg.Address = "39.102.215.201:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return
	}

	// 生成一个检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           "10.0.177.16:50051",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	// 生成一个注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Address = "10.0.177.16"
	registration.Port = 50051
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Tags = []string{"user", "srv", "shop"}
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	zap.S().Info("服务注册完成")
}
