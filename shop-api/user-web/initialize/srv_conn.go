package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop-api/user-web/global"
	"shop-api/user-web/proto"
)

func InitSrvConn() {
	ip := global.ServerConfig.UserSrvConfig.Host
	port := global.ServerConfig.UserSrvConfig.Port
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", ip, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg", err.Error())
	}
	global.UserSrvClient = proto.NewUserClient(conn)
}
