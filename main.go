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

	server := ioc.NewServer(relyConfig)
	middlewares := server.InitGinMiddlewares()
	relyConfig.Trans, _ = server.InitGinTrans()
	engine := server.InitWebServer(middlewares)

	err := engine.Run(fmt.Sprintf("%s:%d", initConfig.ServerConfig.Host, initConfig.ServerConfig.Port))
	if err != nil {
		zap.L().Error("启动失败", zap.Error(err))
	}

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
