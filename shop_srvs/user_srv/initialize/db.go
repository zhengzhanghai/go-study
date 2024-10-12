package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"shop_srvs/user_srv/global"
	"time"
)

func InitDB() {
	fmt.Println("开始连接数据库")
	dbInfo := global.ServerConfig.MySqlInfo
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		dbInfo.User, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})

	// 全局模式
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数，不自动加复数
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("数据库连接完成")

	// 定义一个表结构，将表结构直接生成对应的表
	//err = global.DB.AutoMigrate(&model.User{})
	//if err != nil {
	//	panic(err)
	//}
}
