/**
 * Description：
 * FileName：logger.go
 * Author：CJiaの用心
 * Create：2025/5/13 15:07:15
 * Remark：
 */

package middleware

import (
	"bytes"
	"fmt"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/middleware/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"time"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger(zap *zap.Logger) *Logger {
	return &Logger{
		zap: zap,
	}
}

func (l *Logger) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "swagger") {
			c.Next()
		} else {
			// 开始时间
			start := time.Now()
			path := c.Request.URL.Path
			query := c.Request.URL.Query()

			// 创建自定义读取器
			buffer := &bytes.Buffer{}
			loggingReader := &logger.LoggingReader{
				Reader: c.Request.Body,
				Buffer: buffer,
			}
			c.Request.Body = ioutil.NopCloser(loggingReader)

			c.Next()

			// 结束时间
			timeStamp := time.Now()
			latency := timeStamp.Sub(start)

			l.zap.Debug(path,
				zap.String("time", fmt.Sprintf("%v", latency)),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("ip", c.ClientIP()),
				zap.String("path", path),
				zap.Any("query", query),
				zap.Any("body", loggingReader.Format()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			)
		}
	}
}
