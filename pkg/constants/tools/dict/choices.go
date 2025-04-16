/**
 * Description：
 * FileName：choices.go
 * Author：CJiaの用心
 * Create：2025/4/15 23:08:19
 * Remark：
 */

package dict

import (
	"errors"
	"fmt"
)

// TypeMapping 字典类型映射
var TypeMapping = map[TypeConst]string{
	TypeConst0: "普通字典",
	TypeConst1: "系统字典",
	TypeConst2: "枚举字典",
}

// ConvertDictType 字典类型转换
func ConvertDictType(input TypeConst) (TypeConst, error) {
	if _, exists := TypeMapping[input]; exists {
		return input, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %d，可选值：0/1/2", input))
}

// TypeValueMapping 字典数据类型映射
var TypeValueMapping = map[TypeValueConst]string{
	TypeValueConst0: "字符串",
	TypeValueConst1: "整型",
	TypeValueConst2: "布尔",
}

// ConvertDictTypeValue 字典数据类型转换
func ConvertDictTypeValue(input TypeValueConst) (TypeValueConst, error) {
	if _, exists := TypeValueMapping[input]; exists {
		return input, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %d，可选值：0/1/2", input))
}