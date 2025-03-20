/**
 * Description：
 * FileName：server.go
 * Author：CJiaの用心
 * Create：2025/3/20 22:59:56
 * Remark：
 */

package config

type ServerConfig struct {
	Host string `mapstructure:"host" yaml:"host" json:"host"`
	Port int    `mapstructure:"port" yaml:"port" json:"port"`
}
