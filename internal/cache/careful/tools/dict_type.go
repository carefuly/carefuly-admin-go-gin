/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 20:09:46
 * Remark：
 */

package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrDictTypeNotExist = redis.Nil

type DictTypeCache interface {
	Get(ctx context.Context, id string) (*tools.DictType, error)
	Set(ctx context.Context, domain tools.DictType) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisDictTypeCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisDictTypeCache(cmd redis.Cmdable) DictTypeCache {
	return &RedisDictTypeCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisDictTypeCache) Get(ctx context.Context, id string) (*tools.DictType, error) {
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

	var domain tools.DictType
	err = json.Unmarshal([]byte(data), &domain)
	return &domain, err
}

func (c *RedisDictTypeCache) Set(ctx context.Context, domain tools.DictType) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisDictTypeCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisDictTypeCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisDictTypeCache) key(id string) string {
	return fmt.Sprintf("dictType:info:%s", id)
}
