/**
 * Description：
 * FileName：logger.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:09:44
 * Remark：
 */

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Logger 日志
type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

// GetEncoder 日志记录为console格式
func (l *Logger) GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetLogStdoutWriter 日志记录到控制台
func (l *Logger) GetLogStdoutWriter() zapcore.WriteSyncer {
	// 直接返回标准输出的同步写入器
	return zapcore.AddSync(os.Stdout)
}
