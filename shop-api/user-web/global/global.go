package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"shop-api/user-web/config"
	"shop-api/user-web/proto"
)

var (
	ServerConfig  *config.ServerConfig
	Trans         ut.Translator
	RDB           *redis.Client
	UserSrvClient proto.UserClient
)
