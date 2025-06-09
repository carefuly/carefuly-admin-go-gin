/**
 * Description：
 * FileName：menu_button.go
 * Author：CJiaの用心
 * Create：2025/6/9 11:52:30
 * Remark：
 */

package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrMenuButtonNotExist = redis.Nil

type MenuButtonCache interface {
	Get(ctx context.Context, id string) (*domainSystem.MenuButton, error)
	Set(ctx context.Context, domain domainSystem.MenuButton) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisMenuButtonCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisMenuButtonCache(cmd redis.Cmdable) MenuButtonCache {
	return &RedisMenuButtonCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisMenuButtonCache) Get(ctx context.Context, id string) (*domainSystem.MenuButton, error) {
	key := c.key(id)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrMenuButtonNotExist
		}
		return nil, err
	}

	// 检查是否是防穿透标记
	if data == "not_found" {
		return nil, nil
	}

	var doMain domainSystem.MenuButton
	err = json.Unmarshal([]byte(data), &doMain)
	return &doMain, err
}

func (c *RedisMenuButtonCache) Set(ctx context.Context, domain domainSystem.MenuButton) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisMenuButtonCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisMenuButtonCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisMenuButtonCache) key(id string) string {
	return fmt.Sprintf("careful:system:menu_button:info:%s", id)
}
