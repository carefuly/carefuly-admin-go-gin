/**
 * Description：
 * FileName：import.go
 * Author：CJiaの用心
 * Create：2025/4/16 10:03:34
 * Remark：
 */

package _import

import "strings"

// CleanInput 清理输入中的空格
func CleanInput(input string) string {
	return strings.NewReplacer(" ", "", "\n", "", "\r", "", "\\", "").Replace(input)
}
