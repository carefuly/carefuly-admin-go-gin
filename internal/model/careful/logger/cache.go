/**
 * Description：
 * FileName：cache.go
 * Author：CJiaの用心
 * Create：2025/7/15 20:34:29
 * Remark：
 */

package logger

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CacheLogger 缓存日志表
type CacheLogger struct {
	models.CoreModels
	CacheHost     string `gorm:"type:varchar(100);column:cacheHost;comment:当前主机地址" json:"cacheHost"`       // 当前主机地址
	CacheIp       string `gorm:"type:varchar(100);column:cacheIp;comment:缓存者IP" json:"cacheIp"`            // 缓存者IP
	CacheUsername string `gorm:"type:varchar(40);column:cacheUsername;comment:缓存用户名" json:"cacheUsername"` // 缓存用户名
	CacheMethod   string `gorm:"type:varchar(10);column:cacheMethod;comment:缓存请求方式" json:"cacheMethod"`    // 缓存请求方式
	CachePath     string `gorm:"type:varchar(255);column:cachePath;comment:缓存请求地址" json:"cachePath"`       // 缓存请求地址
	CacheTime     string `gorm:"type:varchar(255);column:cacheTime;comment:缓存记录时间" json:"cacheTime"`       // 缓存记录时间
	CacheKey      string `gorm:"type:varchar(255);column:cacheKey;comment:缓存key键" json:"cacheKey"`         // 缓存请求地址
	CacheValue    string `gorm:"type:text;column:cacheValue;comment:缓存value值" json:"cacheValue"`           // 缓存value值
	CacheError    string `gorm:"type:varchar(255);column:cacheError;comment:缓存Error错误" json:"cacheError"`  // 缓存Error错误
}

func NewCacheLogger() *CacheLogger {
	return &CacheLogger{}
}

func (l *CacheLogger) TableName() string {
	return "careful_logger_cache_log"
}

func (l *CacheLogger) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='缓存日志表'").AutoMigrate(&CacheLogger{})
	if err != nil {
		zap.L().Error("CacheLogger表模型迁移失败", zap.Error(err))
	}
}
