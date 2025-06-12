/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/12 11:02:10
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

var ErrRoleNotExist = redis.Nil

type RoleCache interface {
	Get(ctx context.Context, id string) (*domainSystem.Role, error)
	Set(ctx context.Context, domain domainSystem.Role) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisRoleCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisRoleCache(cmd redis.Cmdable) RoleCache {
	return &RedisRoleCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisRoleCache) Get(ctx context.Context, id string) (*domainSystem.Role, error) {
	key := c.key(id)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrRoleNotExist
		}
		return nil, err
	}

	// 检查是否是防穿透标记
	if data == "not_found" {
		return nil, nil
	}

	var doMain domainSystem.Role
	err = json.Unmarshal([]byte(data), &doMain)
	return &doMain, err
}

func (c *RedisRoleCache) Set(ctx context.Context, domain domainSystem.Role) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisRoleCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisRoleCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisRoleCache) key(id string) string {
	return fmt.Sprintf("careful:system:role:info:%s", id)
}
