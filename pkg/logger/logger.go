/**
 * Description：
 * FileName：logger.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:09:44
 * Remark：
 */

package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
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

// GetJSONEncoder 日志记录为json格式
func (l *Logger) GetJSONEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.CallerKey = "caller"
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GetLogFileWriter 日志记录到文件
func (l *Logger) GetLogFileWriter(path string) zapcore.WriteSyncer {
	serve := strings.Split(path, "/")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./tmp/admin/%s/%s/%s/%s.log", serve[2], time.Now().Format("2006-01-02"), serve[3], serve[4]),
		MaxSize:    5,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
