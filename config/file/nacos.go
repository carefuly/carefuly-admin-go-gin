/**
 * Description：
 * FileName：nacos.go
 * Author：CJiaの用心
 * Create：2025/3/20 22:57:16
 * Remark：
 */

package config

type naCosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}

type NaCosConfig struct {
	ServerConfig ServerConfig `mapstructure:"server"`
	NaCosConfig  naCosConfig  `mapstructure:"nacos"`
}
