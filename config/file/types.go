/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/3/20 22:44:48
 * Remark：
 */

package config

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	ServerConfig   `yaml:"server" json:"server"`
	DatabaseConfig `yaml:"database" json:"database"`
	RedisConfig    `yaml:"redis" json:"redis"`
	TokenConfig    `yaml:"token" json:"token"`
}

type RelyConfig struct {
	Logger *zap.Logger
	Db     *gorm.DB
	Redis  redis.Cmdable
	Trans  ut.Translator
	Token TokenConfig
}
