/**
 * Description：
 * FileName：choices.go
 * Author：CJiaの用心
 * Create：2025/4/17 00:29:47
 * Remark：
 */

package dictType

import (
	"errors"
	"fmt"
)

// DictTagMapping 标签类型映射
var DictTagMapping = map[DictTagConst]string{
	DictTagConstPrimary: "primary",
	DictTagConstSuccess: "success",
	DictTagConstWarning: "warning",
	DictTagConstDanger:  "danger",
	DictTagConstInfo:    "info",
}

// ConvertDictTag 标签类型转换
func ConvertDictTag(input DictTagConst) (DictTagConst, error) {
	if _, exists := DictTagMapping[input]; exists {
		return input, nil
	}
	return "", errors.New(fmt.Sprintf("无效的类型值: %s，可选值：primary/success/warning/danger/info", input))
}
