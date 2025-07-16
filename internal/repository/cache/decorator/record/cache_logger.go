/**
 * Description：
 * FileName：cache_logger.go
 * Author：CJiaの用心
 * Create：2025/7/16 01:07:07
 * Remark：
 */

package record

import (
	"context"
	modelLogger "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// CacheLogger 缓存日志记录器
type CacheLogger struct {
	db *gorm.DB
}

func NewCacheLogger(db *gorm.DB) CacheLogger {
	return CacheLogger{db: db}
}

// Log 异步记录缓存操作日志
func (l *CacheLogger) Log(ctx context.Context, entry *modelLogger.CacheLogger) {
	// 使用goroutine异步记录日志，不影响主流程
	go func() {
		// 设置上下文超时防止日志写入阻塞
		logCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// 使用数据库连接池
		tx := l.db.WithContext(logCtx).Begin()
		if err := tx.Create(entry).Error; err != nil {
			zap.L().Error("缓存日志记录失败",
				zap.String("key", entry.CacheKey),
				zap.String("method", entry.CacheMethod),
				zap.Error(err),
			)
		}

		// 确保提交或回滚
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				zap.L().Error("缓存日志事务异常", zap.Any("recover", r))
			}
		}()

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
		}
	}()
}
