/**
 * Description：
 * FileName：types.go
 * Author：CJiaの用心
 * Create：2025/5/22 16:23:00
 * Remark：
 */

package _import

type ImportResult struct {
	SuccessCount int           // 成功导入的条数
	FailCount    int           // 失败导入的条数
	Errors       []ImportError // 错误信息
}

type ImportError struct {
	Row     int    `json:"row"`     // 数据行号
	Message string `json:"message"` // 错误信息
}

func (r *ImportResult) AddError(row int, message string) {
	r.FailCount++
	r.Errors = append(r.Errors, ImportError{
		Row:     row,
		Message: message,
	})
}
