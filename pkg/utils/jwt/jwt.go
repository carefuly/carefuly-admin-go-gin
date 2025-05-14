/**
 * Description：
 * FileName：jwt.go
 * Author：CJiaの用心
 * Create：2025/5/13 00:45:04
 * Remark：
 */

package jwt

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("无效令牌")
	ErrExpiredToken  = errors.New("令牌已过期")
	ErrTokenNotFound = errors.New("未找到令牌")
)

type Claims struct {
	jwt.RegisteredClaims
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	UserType  int    `json:"userType"`
	UserAgent string `json:"userAgent"`
}

// GenerateToken generates a new JWT token
func GenerateToken(ctx *gin.Context, userId, username string, userType int, secret string, expireHours int) (string, error) {
	// Set claims
	claims := Claims{
		UserId:    userId,
		Username:  username,
		UserType:  userType,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	return token.SignedString([]byte(secret))
}

// ParseToken parses a JWT token and returns the claims
func ParseToken(tokenString, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, ErrTokenNotFound
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Validate token
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
