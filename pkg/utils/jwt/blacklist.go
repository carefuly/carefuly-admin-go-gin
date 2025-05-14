/**
 * Description：
 * FileName：blacklist.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:44:44
 * Remark：
 */

package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// TokenBlacklistPrefix Redis中存储已登出token的前缀
	TokenBlacklistPrefix = "token:blacklist:"
)

// TokenBlacklist JWT Token黑名单实现
type TokenBlacklist struct {
	rdb redis.Cmdable
}

// NewTokenBlacklist 创建一个Token黑名单实例
func NewTokenBlacklist(rdb redis.Cmdable) *TokenBlacklist {
	return &TokenBlacklist{
		rdb: rdb,
	}
}

// Add 将Token加入黑名单
// tokenStr: JWT token字符串
// expiresIn: token的剩余有效期（秒）
func (b *TokenBlacklist) Add(ctx context.Context, tokenStr string, expiresIn time.Duration) error {
	// 使用Redis SET命令将token加入黑名单，并设置与token相同的过期时间
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, tokenStr)
	return b.rdb.Set(ctx, key, "1", expiresIn).Err()
}

// IsBlacklisted 检查Token是否在黑名单中
func (b *TokenBlacklist) IsBlacklisted(ctx context.Context, tokenStr string) (bool, error) {
	key := fmt.Sprintf("%s%s", TokenBlacklistPrefix, tokenStr)
	result, err := b.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}
