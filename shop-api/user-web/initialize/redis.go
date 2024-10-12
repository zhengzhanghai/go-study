package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"shop-api/user-web/global"
)

func InitRedis() {
	global.RDB = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.ServerConfig.RedisInfo.Host,
			global.ServerConfig.RedisInfo.Port),
	})
}
