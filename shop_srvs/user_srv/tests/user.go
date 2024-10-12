package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.NewClient(
		"127.0.0.1:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func main() {
	Init()
	//TestCreateUser()
	//TestGetUserList()
	//TestGetUserByMobile()
	//TestGetUserById()
	//TestUpdateUser()
	conn.Close()
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		resp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Mobile:   fmt.Sprintf("1359876567%d", i),
			Nickname: fmt.Sprintf("章海_%d", i),
			Password: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Id)
	}
}

func TestGetUserList() {
	resp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range resp.Users {
		fmt.Println(user.Id, user.Mobile, user.Nickname, user.Password)
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestGetUserByMobile() {
	user, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: "13598765670",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Id, user.Mobile, user.Nickname, user.Password)
	checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          "admin123",
		EncryptedPassword: user.Password,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(checkRsp.Success)
}

func TestGetUserById() {
	user, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: "20",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user.Id, user.Mobile, user.Nickname, user.Password)
	checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          "admin123",
		EncryptedPassword: user.Password,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(checkRsp.Success)
}

func TestUpdateUser() {
	user, err := userClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: "9",
	})
	if err != nil {
		panic(err)
	}
	_, err = userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       user.Id,
		Gender:   "female",
		Nickname: "哈哈",
		Birthday: 1724833158,
	})
	if err != nil {
		panic(err)
	}
}
