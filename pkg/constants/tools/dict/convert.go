/**
 * Description：
 * FileName：convert.go
 * Author：CJiaの用心
 * Create：2025/4/15 23:08:29
 * Remark：
 */

package dict

import (
	"errors"
	"fmt"
)

// TypeImportMapping 字典类型映射
var TypeImportMapping = map[string]TypeConst{
	"普通字典": TypeConst0,
	"系统字典": TypeConst1,
	"枚举字典": TypeConst2,
}

// ConvertDictTypeImport 字典类型转换
func ConvertDictTypeImport(input string) (TypeConst, error) {
	if val, exists := TypeImportMapping[input]; exists {
		return val, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %s，可选值：普通字典/系统字典/系统字典", input))
}

// TypeValueImportMapping 字典数据类型映射
var TypeValueImportMapping = map[string]TypeValueConst{
	"字符串": TypeValueConst0,
	"整型":  TypeValueConst1,
	"布尔":  TypeValueConst2,
}

// ConvertDictTypeValueImport 字典数据类型转换
func ConvertDictTypeValueImport(input string) (TypeValueConst, error) {
	if val, exists := TypeValueImportMapping[input]; exists {
		return val, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %s，可选值：字符串/整型/布尔", input))
}
