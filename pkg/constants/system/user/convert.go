/**
 * Description：
 * FileName：convert.go
 * Author：CJiaの用心
 * Create：2025/4/16 15:43:06
 * Remark：
 */

package user

import (
	"errors"
	"fmt"
)

// GenderImportMapping 性别类型映射
var GenderImportMapping = map[string]GenderConst{
	"男":    GenderMale,
	"女":    GenderFemale,
	"未知":   GenderUnknown,
	"不男不女": GenderMaleNorFemale,
}

// ConvertGenderImport 性别类型转换
func ConvertGenderImport(input string) (GenderConst, error) {
	if val, exists := GenderImportMapping[input]; exists {
		return val, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %s，可选值：男/女/未知/不男不女", input))
}
