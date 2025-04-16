/**
 * Description：
 * FileName：convert.go
 * Author：CJiaの用心
 * Create：2025/4/17 00:33:17
 * Remark：
 */

package dictType

import (
	"errors"
	"fmt"
)

// DictTagImportMapping 标签类型映射
var DictTagImportMapping = map[string]DictTagConst{
	"primary": DictTagConstPrimary,
	"success": DictTagConstSuccess,
	"warning": DictTagConstWarning,
	"danger":  DictTagConstDanger,
	"info":    DictTagConstInfo,
}

// ConvertDictTagImport 标签类型转换
func ConvertDictTagImport(input string) (DictTagConst, error) {
	if val, exists := DictTagImportMapping[input]; exists {
		return val, nil
	}
	return "", errors.New(fmt.Sprintf("无效的类型值: %s，可选值：primary/success/warning/danger/info", input))
}
