package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // GRPC负载均衡
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
	"shop-api/user-web/utils"
	validator2 "shop-api/user-web/validator"
	"syscall"
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

	// 如果是本地开发环境，端口号固定，生产环境自动获取
	viper.AutomaticEnv()
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

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

	client, serviceId := RegisterService()

	port := global.ServerConfig.Port
	zap.S().Info("启动服务器，端口: ", port)
	go func() {
		if err := Router.Run(fmt.Sprintf(":%d", port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()

	// 接收终止信号，在consul中注销
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Error("注销失败")
	}
	zap.S().Info("注销成功")
}

func RegisterService() (*api.Client, string) {
	return Register("10.0.177.16", 8021, "user-web", []string{"user", "web", "api"}, "user-web")
}

// 注册服务
func Register(address string, port int, name string, tags []string, id string) (*api.Client, string) {
	cfg := api.DefaultConfig()
	cfg.Address = "39.102.215.201:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return client, ""
	}

	// 生成一个检查对象
	//check := &api.AgentServiceCheck{
	//	HTTP:                           "http://10.0.177.16:8021/health",
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "10s",
	//}
	// 生成一个注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Address = address
	registration.Port = port
	registration.Name = name
	registration.ID = id
	registration.Tags = tags
	//registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return client, id
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "39.102.215.201:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
		return
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}

func FilterServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "39.102.215.201:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
		return
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
		return
	}
	for key, _ := range data {
		fmt.Println(key)
	}
}
