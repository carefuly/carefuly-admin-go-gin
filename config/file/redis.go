/**
 * Description：
 * FileName：redis.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:02:18
 * Remark：
 */

package config

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	Db       int    `yaml:"db" json:"db"`
}
