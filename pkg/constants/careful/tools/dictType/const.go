/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/5/23 16:02:30
 * Remark：
 */

package dictType

type DictTagConst string // 标签类型

const (
	DictTagConstPrimary DictTagConst = "primary" // primary
	DictTagConstSuccess DictTagConst = "success" // success
	DictTagConstWarning DictTagConst = "warning" // warning
	DictTagConstDanger  DictTagConst = "danger"  // danger
	DictTagConstInfo    DictTagConst = "info"    // info
)

// DictTagMapping 标签类型映射
var DictTagMapping = map[DictTagConst]string{
	DictTagConstPrimary: "primary",
	DictTagConstSuccess: "success",
	DictTagConstWarning: "warning",
	DictTagConstDanger:  "danger",
	DictTagConstInfo:    "info",
}

// DictTagImportMapping 标签类型映射
var DictTagImportMapping = map[string]DictTagConst{
	"primary": DictTagConstPrimary,
	"success": DictTagConstSuccess,
	"warning": DictTagConstWarning,
	"danger":  DictTagConstDanger,
	"info":    DictTagConstInfo,
}

// BoolValueImportMapping 布尔值类型映射
var BoolValueImportMapping = map[string]bool{
	"是": true,
	"否": false,
}
