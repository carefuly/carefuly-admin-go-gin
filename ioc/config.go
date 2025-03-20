/**
 * Description：
 * FileName：config.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:03:55
 * Remark：
 */

package ioc

import (
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(debug bool) *config.Config {
	// 从配置文件中读取配置信息
	configFilePrefix := "config"
	configFileName := fmt.Sprintf(".\\config\\%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf(".\\config\\%s-dev.yaml", configFilePrefix)
	}
	conf := new(config.NaCosConfig)

	// 配置
	v := viper.New()
	// 文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.L().Error("配置文件读取失败", zap.Error(err))
	}
	// 全局变量
	if err := v.Unmarshal(conf); err != nil {
		panic(err)
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("[监听配置文件修改]")
		zap.L().Debug("配置文件发生变动")
		if err := v.Unmarshal(conf); err != nil {
			zap.L().Error("配置文件发生变动，解析失败", zap.Error(err))
		}
		zap.L().Debug("配置文件发生变动", zap.Any("conf", conf))
	})

	zap.L().Debug("[监听配置文件修改]", zap.Any("conf", conf))

	var globalConfig = new(config.Config)
	// 将配置信息返回
	return globalConfig
}
