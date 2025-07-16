/**
 * Description：
 * FileName：storage.go
 * Author：CJiaの用心
 * Create：2025/5/13 15:09:07
 * Remark：
 */

package middleware

import (
	"bytes"
	"fmt"
	loggerModel "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/logger"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/logger"
	loggerMiddleware "github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/middleware/logger"
	_import "github.com/carefuly/carefuly-admin-go-gin/pkg/utils/import"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/requestUtils"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
	"time"
)

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) GetResValue(c *gin.Context, key string) string {
	value, exists := c.Get(key)
	if !exists {
		return ""
	}
	return value.(string)
}

func (s *Storage) Logger(path string) *zap.Logger {
	l := logger.NewLogger()
	syncer := l.GetLogFileWriter(path)
	encoder := l.GetJSONEncoder()
	core := zapcore.NewCore(encoder, syncer, zapcore.DebugLevel)
	return zap.New(core, zap.AddCaller())
}

func (s *Storage) StorageLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if s.containsAnySubstring(c.Request.URL.Path, []string{"swagger", "static", "export"}) {
			c.Next()
		} else {
			// 开始时间
			start := time.Now()
			path := c.Request.URL.Path
			query := c.Request.URL.RawQuery
			contentType := c.GetHeader("Content-Type")

			// 创建自定义读取器
			buffer := &bytes.Buffer{}
			loggingReader := &loggerMiddleware.LoggingReader{
				Reader: c.Request.Body,
				Buffer: buffer,
			}
			c.Request.Body = ioutil.NopCloser(loggingReader)

			// 创建自定义响应写入器
			crw := &loggerMiddleware.CustomGinResponseWriter{
				ResponseWriter: c.Writer,
				Body:           bytes.NewBuffer(nil),
			}
			c.Writer = crw

			c.Next()

			// 结束时间
			timeStamp := time.Now()
			latency := timeStamp.Sub(start)

			// 获取响应数据
			responseBody := crw.Body.String()
			responseJson := crw.Format(responseBody)

			var record loggerModel.OperateLogger

			record.RequestUsername = requestUtils.GetRequestUser(c)
			record.RequestTime = fmt.Sprintf("%v", latency)

			record.RequestStatus = c.Writer.Status()
			record.RequestMethod = c.Request.Method
			record.RequestIp = requestUtils.NormalizeIP(c)

			record.RequestPath = path
			record.RequestQuery = query

			var body string
			if strings.Contains(contentType, "multipart/form-data") {
				body = "上传文件"
			} else {
				body = _import.CleanInput(buffer.String())
			}
			record.RequestBody = body

			ua := user_agent.New(c.Request.UserAgent())
			browserName, _ := ua.Browser()
			record.RequestOs = ua.OS()
			record.RequestBrowser = browserName
			record.UserAgent = c.Request.UserAgent()

			record.RequestCode = responseJson.Code
			record.RequestResult = responseBody

			record.Errors = c.Errors.ByType(gin.ErrorTypePrivate).String()
			record.Internal = s.GetResValue(c, "internal")

			if record.RequestMethod != "GET" { // GET不进行持久化记录
				record.Insert(c, db, record)
			}

			l := s.Logger(record.RequestPath)
			l.Info(path,
				zap.String("requestUsername", record.RequestUsername),
				zap.String("requestTime", fmt.Sprintf("%v", latency)),
				zap.Int("requestStatus", c.Writer.Status()),
				zap.String("requestMethod", c.Request.Method),
				zap.String("requestIp", requestUtils.NormalizeIP(c)),
				zap.String("requestPath", path),
				zap.Any("requestQuery", query),
				zap.Any("requestBody", loggingReader.Format()),
				zap.String("requestOs", record.RequestOs),
				zap.String("requestBrowser", record.RequestBrowser),
				zap.String("userAgent", record.UserAgent),
				zap.Int("requestCode", record.RequestCode),
				zap.Any("requestResult", responseJson),
				zap.String("requestErrors", record.Errors),
				zap.String("requestInternal", record.Internal),
			)
		}
	}
}

// 检查字符串是否包含切片中的任意一个子串
func (s *Storage) containsAnySubstring(str string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}
