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

// BoolValueImportMapping 布尔值类型映射
var BoolValueImportMapping = map[string]bool{
	"是": true,
	"否": false,
}

// ConvertBoolValueImport 布尔值类型转换
func ConvertBoolValueImport(input string) (bool, error) {
	if val, exists := BoolValueImportMapping[input]; exists {
		return val, nil
	}
	return false, errors.New(fmt.Sprintf("无效的类型值: %s，可选值：是/否", input))
}
