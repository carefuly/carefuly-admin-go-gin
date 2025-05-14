/**
 * Description：
 * FileName：cache.go
 * Author：CJiaの用心
 * Create：2025/5/11 17:55:50
 * Remark：
 */

package config

type CacheConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	Db       int    `yaml:"db" json:"db"`
}
