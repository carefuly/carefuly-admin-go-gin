/**
 * Description：
 * FileName：bucket_logging_decorator.go
 * Author：CJiaの用心
 * Create：2025/7/16 09:51:30
 * Remark：
 */

package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelLogger "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/logger"
	cacheTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/careful/tools"
	cacheRecord "github.com/carefuly/carefuly-admin-go-gin/internal/repository/cache/decorator/record"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"net/http"
	"time"
)

type BucketCacheLoggingDecorator struct {
	cache  cacheTools.BucketCache
	logger cacheRecord.CacheLogger
}

func NewBucketCacheLoggingDecorator(cache cacheTools.BucketCache, logger cacheRecord.CacheLogger) BucketCacheLoggingDecorator {
	return BucketCacheLoggingDecorator{
		cache:  cache,
		logger: logger,
	}
}

func (d *BucketCacheLoggingDecorator) Get(ctx context.Context, id string) (*domainTools.Bucket, error) {
	start := time.Now()
	result, err := d.cache.Get(ctx, id)

	// 请求头
	request := ctx.Value("request").(*http.Request)

	// 创建日志条目
	entry := &modelLogger.CacheLogger{
		CoreModels: models.CoreModels{
			Creator:    ctx.Value("userId").(string),
			Modifier:   ctx.Value("userId").(string),
			BelongDept: ctx.Value("deptId").(string),
		},
		CacheHost:     request.Host,
		CacheIp:       ctx.Value("requestIp").(string),
		CacheUsername: ctx.Value("username").(string),
		CacheMethod:   request.Method,
		CachePath:     request.URL.Path,
		CacheKey:      d.key(id),
	}

	if err != nil {
		if errors.Is(err, cacheTools.ErrBucketNotExist) {
			entry.CacheValue = "not_found"
		} else {
			entry.CacheError = err.Error()
		}
	} else if result != nil {
		// 记录值摘要
		if data, err := json.Marshal(result); err == nil {
			entry.CacheValue = string(data)
		}
	}

	// 记录执行时间
	entry.CacheTime = time.Since(start).String()

	// 异步记录日志
	go d.logger.Log(ctx, entry)

	return result, err
}

func (d *BucketCacheLoggingDecorator) Set(ctx context.Context, domain domainTools.Bucket) error {
	start := time.Now()
	err := d.cache.Set(ctx, domain)

	// 请求头
	request := ctx.Value("request").(*http.Request)

	// 记录值摘要
	valueStr := ""
	if data, err := json.Marshal(domain); err == nil {
		valueStr = string(data)
	}

	// 创建日志条目
	entry := &modelLogger.CacheLogger{
		CoreModels: models.CoreModels{
			Creator:    ctx.Value("userId").(string),
			Modifier:   ctx.Value("userId").(string),
			BelongDept: ctx.Value("deptId").(string),
		},
		CacheHost:     request.Host,
		CacheIp:       ctx.Value("requestIp").(string),
		CacheUsername: ctx.Value("username").(string),
		CacheMethod:   request.Method,
		CachePath:     request.URL.Path,
		CacheKey:      d.key(domain.Id),
		CacheValue:    valueStr,
	}

	if err != nil {
		entry.CacheError = err.Error()
	}

	// 记录执行时间
	entry.CacheTime += time.Since(start).String()

	// 异步记录日志
	go d.logger.Log(ctx, entry)

	return err
}

func (d *BucketCacheLoggingDecorator) Del(ctx context.Context, id string) error {
	start := time.Now()
	err := d.cache.Del(ctx, id)

	// 请求头
	request := ctx.Value("request").(*http.Request)

	// 创建日志条目
	entry := &modelLogger.CacheLogger{
		CoreModels: models.CoreModels{
			Creator:    ctx.Value("userId").(string),
			Modifier:   ctx.Value("userId").(string),
			BelongDept: ctx.Value("deptId").(string),
		},
		CacheHost:     request.Host,
		CacheIp:       ctx.Value("requestIp").(string),
		CacheUsername: ctx.Value("username").(string),
		CacheMethod:   request.Method,
		CachePath:     request.URL.Path,
		CacheKey:      d.key(id),
	}

	if err != nil {
		entry.CacheError = err.Error()
	}

	// 记录执行时间
	entry.CacheTime = time.Since(start).String()

	// 异步记录日志
	go d.logger.Log(ctx, entry)

	return err
}

func (d *BucketCacheLoggingDecorator) SetNotFound(ctx context.Context, id string) error {
	start := time.Now()
	err := d.cache.SetNotFound(ctx, id)

	// 请求头
	request := ctx.Value("request").(*http.Request)

	// 创建日志条目
	entry := &modelLogger.CacheLogger{
		CoreModels: models.CoreModels{
			Creator:    ctx.Value("userId").(string),
			Modifier:   ctx.Value("userId").(string),
			BelongDept: ctx.Value("deptId").(string),
		},
		CacheHost:     request.Host,
		CacheIp:       ctx.Value("requestIp").(string),
		CacheUsername: ctx.Value("username").(string),
		CacheMethod:   request.Method,
		CachePath:     request.URL.Path,
		CacheKey:      d.key(id),
		CacheValue:    "not_found",
	}

	if err != nil {
		entry.CacheError = err.Error()
	}

	// 记录执行时间
	entry.CacheTime += " | " + time.Since(start).String()

	// 异步记录日志
	go d.logger.Log(ctx, entry)

	return err
}

func (d *BucketCacheLoggingDecorator) key(id string) string {
	return fmt.Sprintf("careful:tools:bucket:info:%s", id)
}
