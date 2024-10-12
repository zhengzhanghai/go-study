package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/user_srv/model"
	"time"
)

func main() {
	fmt.Println(genMD5("123456"))
}

func connectDb() {
	dsn := "root:root@tcp(10.0.180.168:3306)/shop_user_srv?charset=utf8"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数，不自动加复数
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// 定义一个表结构，将表结构直接生成对应的表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}

func genMD5(code string) string {
	hash := md5.Sum([]byte(code))
	return hex.EncodeToString(hash[:])
}
