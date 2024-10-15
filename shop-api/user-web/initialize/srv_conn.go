package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop-api/user-web/global"
	"shop-api/user-web/proto"
)

func InitSrvConn() {
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return
	}

	data, err := client.Agent().ServicesWithFilter(
		fmt.Sprintf("Service == \"%s\"", global.ServerConfig.UserSrvConfig.Name))
	if err != nil {
		panic(err)
		return
	}

	if len(data) == 0 {
		panic("没有发现user-srv服务")
		return
	}

	ip := ""
	port := 0
	for key, value := range data {
		fmt.Println(key)
		ip = value.Address
		port = value.Port
		break
	}

	zap.S().Info(fmt.Sprintf("用户微服务 ip: %s, port: %d", ip, port))

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", ip, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg", err.Error())
	}
	global.UserSrvClient = proto.NewUserClient(conn)
}
