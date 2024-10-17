package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // GRPC负载均衡
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop-api/user-web/proto"
)

func main() {
	fmt.Println("start ......")
	testUserSrvLoaderBalance()
	fmt.Println("end ......")
}

func testUserSrvLoaderBalance() {
	conn, err := grpc.NewClient(
		"consul://39.102.215.201:8500/user-srv?wait=14s",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))

	if err != nil {
		zap.S().Errorw("[GetUserList] 连接用户服务失败", "msg", err.Error())
	}
	client := proto.NewUserClient(conn)
	resp, err := client.GetUserList(context.Background(), &proto.PageInfo{
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		panic(err)
	}
	if resp != nil && len(resp.Users) > 0 {
		for user := range resp.Users {
			fmt.Println(user)
		}
	}
	fmt.Println("获取完成 ......")
}
