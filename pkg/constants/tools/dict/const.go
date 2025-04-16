/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/4/15 23:06:58
 * Remark：
 */

package dict

type TypeConst int // 字典类型

const (
	TypeConst0 TypeConst = iota // 普通字典
	TypeConst1                  // 系统字典
	TypeConst2                  // 枚举字典
)

type TypeValueConst int // 字典数据类型

const (
	TypeValueConst0 TypeValueConst = iota // 字符串
	TypeValueConst1                       // 整型
	TypeValueConst2                       // 布尔
)
