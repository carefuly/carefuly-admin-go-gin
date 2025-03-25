/**
 * Description：
 * FileName：response.go
 * Author：CJiaの用心
 * Create：2025/2/19 16:05:51
 * Remark：
 */

package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`    // 状态码
	Data    interface{} `json:"data"`    // 数据
	Msg     interface{} `json:"msg"`     // 提示信息
	Success bool        `json:"success"` // 是否成功
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) SuccessResponse(ctx *gin.Context, msg any, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"msg":     msg,
		"success": true,
		"data":    data,
	})
}

func (r *Response) ErrorResponse(ctx *gin.Context, code int, msg any, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    code,
		"msg":     msg,
		"success": false,
		"data":    data,
	})
}
