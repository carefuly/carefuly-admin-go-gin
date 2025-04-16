/**
 * Description：
 * FileName：constants.go
 * Author：CJiaの用心
 * Create：2025/4/14 16:48:15
 * Remark：
 */

package utils

import (
	"errors"
	"fmt"
)

// DictTypeMapping 字典类型映射
var DictTypeMapping = map[int]string{
	0: "普通字典",
}

// ConvertDictType 字典类型类型
func ConvertDictType(input int) (int, error) {
	if _, exists := DictTypeMapping[input]; exists {
		return input, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %d，可选值：0", input))
}

// DictValueTypeMapping 字典数据类型映射
var DictValueTypeMapping = map[int]string{
	0: "字符串",
	1: "整型",
	2: "布尔类型",
}

// ConvertDictValueType 字典数据类型类型
func ConvertDictValueType(input int) (int, error) {
	if _, exists := DictValueTypeMapping[input]; exists {
		return input, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %d，可选值：0, 1, 2", input))
}
