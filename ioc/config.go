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
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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

	// 从NaCos中读取配置信息
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: conf.NaCosConfig.Host,
			Port:   conf.NaCosConfig.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         conf.NaCosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            conf.NaCosConfig.User,
		Password:            conf.NaCosConfig.Password,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})

	if err != nil {
		zap.L().Error("创建NaCos配置失败", zap.Error(err))
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: conf.NaCosConfig.DataId,
		Group:  conf.NaCosConfig.Group,
	})
	if err != nil {
		zap.L().Debug("读取NaCos配置失败", zap.Error(err))
	}

	var globalConfig = new(config.Config)
	// 将配置信息写入到全局变量中
	err = yaml.Unmarshal([]byte(content), &globalConfig)
	if err != nil {
		zap.L().Error("解析配置文件失败", zap.Error(err))
	}

	globalConfig.ServerConfig.Host = conf.ServerConfig.Host
	globalConfig.ServerConfig.Port = conf.ServerConfig.Port

	zap.L().Debug("[配置文件信息]", zap.Any("globalConfig", globalConfig))

	// 将配置信息返回
	return globalConfig
}
