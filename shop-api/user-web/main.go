package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
	validator2 "shop-api/user-web/validator"
)

func main() {
	initialize.InitLogger()
	initialize.InitConfig()
	// 初始化router
	Router := initialize.Routers()

	// 初始化Redis
	initialize.InitRedis()

	// 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	// 初始化srv链接
	initialize.InitSrvConn()

	// 自定义gin验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册验证器
		v.RegisterValidation("mobile", validator2.ValidateMobile)
		// 注册验证器错误翻译
		v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	port := global.ServerConfig.Port
	zap.S().Info("启动服务器，端口: ", port)
	if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}

}
