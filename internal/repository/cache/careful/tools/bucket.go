/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 16:49:29
 * Remark：
 */

package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrBucketNotExist = redis.Nil

type BucketCache interface {
	Get(ctx context.Context, id string) (*domainTools.Bucket, error)
	Set(ctx context.Context, domain domainTools.Bucket) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisBucketCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisBucketCache(cmd redis.Cmdable) BucketCache {
	return &RedisBucketCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisBucketCache) Get(ctx context.Context, id string) (*domainTools.Bucket, error) {
	key := c.key(id)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrBucketNotExist
		}
		return nil, err
	}

	// 检查是否是防穿透标记
	if data == "not_found" {
		return nil, nil
	}

	var doMain domainTools.Bucket
	err = json.Unmarshal([]byte(data), &doMain)
	return &doMain, err
}

func (c *RedisBucketCache) Set(ctx context.Context, domain domainTools.Bucket) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisBucketCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisBucketCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisBucketCache) key(id string) string {
	return fmt.Sprintf("careful:tools:bucket:info:%s", id)
}
