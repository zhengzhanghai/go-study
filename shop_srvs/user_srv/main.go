package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/handler"
	"shop_srvs/user_srv/initialize"
	"shop_srvs/user_srv/proto"
	"shop_srvs/user_srv/utils"
	"syscall"
)

func main() {
	// 获取传进来的IP和端口
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()

	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

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
	client, serviceId := RegisterService(*Port)

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc: " + err.Error())
		}
	}()

	// 接收终止信号，在consul中注销
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Error("注销失败")
	}
	zap.S().Info("注销成功")
}

func RegisterService(port int) (*api.Client, string) {
	cfg := api.DefaultConfig()
	// consul部署的服务器地址
	cfg.Address = "39.102.215.201:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return client, ""
	}

	// 生成一个检查对象
	//check := &api.AgentServiceCheck{
	//	GRPC:                           "10.0.177.16:50051",
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "10s",
	//}
	// 生成一个注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Address = "10.0.177.16"
	registration.Port = port
	registration.Name = global.ServerConfig.Name
	// 防止consul多实例覆盖
	id, _ := uuid.GenerateUUID()
	registration.ID = id
	registration.Tags = []string{"user", "srv", "shop"}
	//registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	zap.S().Info("服务注册完成")
	return client, id
}
