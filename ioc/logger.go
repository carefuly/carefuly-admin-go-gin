/**
 * Description：
 * FileName：logger.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:08:55
 * Remark：
 */

package ioc

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitStdoutLogger 初始化终端日志记录器
func InitStdoutLogger() *zap.Logger {
	stdoutLogger := logger.NewLogger()
	writeSyncer := stdoutLogger.GetLogStdoutWriter()
	encoder := stdoutLogger.GetEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 全局日志
	var globalLogger = new(zap.Logger)
	globalLogger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(globalLogger)
	// 日志记录器
	return globalLogger
}
