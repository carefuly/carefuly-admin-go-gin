/**
 * Description：
 * FileName：request_utils.go
 * Author：CJiaの用心
 * Create：2025/4/8 11:19:13
 * Remark：
 */

package requestUtils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// GetRequestUser 获取请求user
func GetRequestUser(c *gin.Context) string {
	username, ok := c.Get("username")
	if !ok {
		return "AnonymousUser"
	}
	return username.(string)
}

// NormalizeIP IPv6 本地地址处理
func NormalizeIP(c *gin.Context) string {
	ip := c.ClientIP()

	// 如果地址是 IPv6 本地回环
	if ip == "::1" {
		return "127.0.0.1"
	}

	// 如果地址是 IPv4 本地回环
	if ip == "0:0:0:0:0:0:0:1" {
		return "127.0.0.1"
	}

	// 处理IPv6地址的方括号
	if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
		return ip[1 : len(ip)-1]
	}

	return ip
}
