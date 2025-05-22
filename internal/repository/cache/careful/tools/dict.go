/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/12 15:06:44
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

var ErrDictNotExist = redis.Nil

type DictCache interface {
	Get(ctx context.Context, id string) (*domainTools.Dict, error)
	Set(ctx context.Context, domain domainTools.Dict) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisDictCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisDictCache(cmd redis.Cmdable) DictCache {
	return &RedisDictCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisDictCache) Get(ctx context.Context, id string) (*domainTools.Dict, error) {
	key := c.key(id)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrDictNotExist
		}
		return nil, err
	}

	// 检查是否是防穿透标记
	if data == "not_found" {
		return nil, nil
	}

	var doMain domainTools.Dict
	err = json.Unmarshal([]byte(data), &doMain)
	return &doMain, err
}

func (c *RedisDictCache) Set(ctx context.Context, domain domainTools.Dict) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisDictCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisDictCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisDictCache) key(id string) string {
	return fmt.Sprintf("careful:tools:dict:info:%s", id)
}
