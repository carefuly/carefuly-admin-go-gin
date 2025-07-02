/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/5/23 14:48:09
 * Remark：
 */

package dict

type TypeConst int // 字典类型

const (
	TypeConstOrdinary TypeConst = iota + 1 // 普通字典
	TypeConstSystem                        // 系统字典
	TypeConstEnum                          // 枚举字典
)

// TypeMapping 字典类型映射
var TypeMapping = map[TypeConst]string{
	TypeConstOrdinary: "普通字典",
	TypeConstSystem:   "系统字典",
	TypeConstEnum:     "枚举字典",
}

// TypeImportMapping 字典类型映射
var TypeImportMapping = map[string]TypeConst{
	"普通字典": TypeConstOrdinary,
	"系统字典": TypeConstSystem,
	"枚举字典": TypeConstEnum,
}

type TypeValueConst int // 数据类型

const (
	TypeValueConstStr  TypeValueConst = iota + 1 // 字符串
	TypeValueConstInt                            // 整型
	TypeValueConstBool                           // 布尔
)

// TypeValueMapping 数据类型映射
var TypeValueMapping = map[TypeValueConst]string{
	TypeValueConstStr:  "字符串",
	TypeValueConstInt:  "整型",
	TypeValueConstBool: "布尔",
}

// TypeValueImportMapping 数据类型映射
var TypeValueImportMapping = map[string]TypeValueConst{
	"字符串": TypeValueConstStr,
	"整型":  TypeValueConstInt,
	"布尔":  TypeValueConstBool,
}
