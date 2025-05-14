/**
 * Description：
 * FileName：import.go
 * Author：CJiaの用心
 * Create：2025/5/13 15:11:00
 * Remark：
 */

package _import

import "strings"

// CleanInput 清理输入中的空格
func CleanInput(input string) string {
	return strings.NewReplacer(" ", "", "\n", "", "\r", "", "\\", "").Replace(input)
}
