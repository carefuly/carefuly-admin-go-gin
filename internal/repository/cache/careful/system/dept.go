/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/6/11 15:45:53
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

var ErrDeptNotExist = redis.Nil

type DeptCache interface {
	Get(ctx context.Context, id string) (*domainSystem.Dept, error)
	Set(ctx context.Context, domain domainSystem.Dept) error
	Del(ctx context.Context, id string) error
	SetNotFound(ctx context.Context, id string) error // 防止缓存穿透
}

type RedisDeptCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func NewRedisDeptCache(cmd redis.Cmdable) DeptCache {
	return &RedisDeptCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *RedisDeptCache) Get(ctx context.Context, id string) (*domainSystem.Dept, error) {
	key := c.key(id)

	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrDeptNotExist
		}
		return nil, err
	}

	// 检查是否是防穿透标记
	if data == "not_found" {
		return nil, nil
	}

	var doMain domainSystem.Dept
	err = json.Unmarshal([]byte(data), &doMain)
	return &doMain, err
}

func (c *RedisDeptCache) Set(ctx context.Context, domain domainSystem.Dept) error {
	key := c.key(domain.Id)
	data, err := json.Marshal(domain)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

func (c *RedisDeptCache) Del(ctx context.Context, id string) error {
	key := c.key(id)
	return c.cmd.Del(ctx, key).Err()
}

func (c *RedisDeptCache) SetNotFound(ctx context.Context, id string) error {
	key := c.key(id)
	// 设置短暂的有效期防止缓存穿透
	return c.cmd.Set(ctx, key, "not_found", time.Minute).Err()
}

func (c *RedisDeptCache) key(id string) string {
	return fmt.Sprintf("careful:system:dept:info:%s", id)
}
