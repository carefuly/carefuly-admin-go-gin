/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/4/17 00:24:19
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
