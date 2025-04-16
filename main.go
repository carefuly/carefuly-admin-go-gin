/**
 * Description：
 * FileName：main.go
 * Author：CJiaの用心
 * Create：2025/3/20 22:48:49
 * Remark：
 */

package main

import (
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/ioc"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// @title CarefulAdmin
// @version 1.0
// @description CarefulAdmin在线接口文档
// @termsOfService http://swagger.io/terms
// @contact.name CJiaの用心
// @contact.url http://www.swagger.io/support
// @contact.email 2224693191@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /dev-api
// @securityDefinitions.apikey  LoginToken
// @in                          header
// @name                        Authorization
func main() {
	var relyConfig config.RelyConfig
	relyConfig.Logger = ioc.InitStdoutLogger()
	initConfig := ioc.InitConfig(true)

	dbPool := ioc.NewDbPool()
	dbPool.InitDatabases(initConfig.DatabaseConfig)
	relyConfig.Db = config.DatabasesPool{
		Careful: dbPool.CarefulDB,
	}

	relyConfig.Redis = ioc.InitRedis(initConfig.RedisConfig)
	relyConfig.Token = initConfig.TokenConfig

	server := ioc.NewServer(relyConfig, "zh")
	middlewares := server.InitGinMiddlewares(relyConfig)
	relyConfig.Trans, _ = server.InitGinTrans()
	engine := server.InitWebServer(middlewares, relyConfig)

	err := engine.Run(fmt.Sprintf("%s:%d", initConfig.ServerConfig.Host, initConfig.ServerConfig.Port))
	if err != nil {
		zap.L().Error("启动失败", zap.Error(err))
	}

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
