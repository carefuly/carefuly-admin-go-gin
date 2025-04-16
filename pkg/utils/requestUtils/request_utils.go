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
)

// GetRequestUser 获取请求user
func GetRequestUser(c *gin.Context) string {
	username, ok := c.Get("username")
	if !ok {
		return "AnonymousUser"
	}
	return username.(string)
}
