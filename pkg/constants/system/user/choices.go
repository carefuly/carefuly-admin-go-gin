/**
 * Description：
 * FileName：choices.go
 * Author：CJiaの用心
 * Create：2025/4/16 15:41:04
 * Remark：
 */

package user

import (
	"errors"
	"fmt"
)

// GenderMapping 性别映射
var GenderMapping = map[GenderConst]string{
	GenderMale: "男",
	GenderFemale: "女",
	GenderUnknown: "未知",
	GenderMaleNorFemale: "不男不女",
}

// ConvertGender 性别类型转换
func ConvertGender(input GenderConst) (GenderConst, error) {
	if _, exists := GenderMapping[input]; exists {
		return input, nil
	}
	return -1, errors.New(fmt.Sprintf("无效的类型值: %d，可选值：0/1/2/3", input))
}