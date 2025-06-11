/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/6/10 17:00:28
 * Remark：
 */

package role

type DataRangeConst int

const (
	DataRangeConstOnly      DataRangeConst = iota + 1 // 仅本人数据权限
	DataRangeConstDept                                // 本部门数据权限
	DataRangeConstDeptBelow                           // 本部门及以下数据权限
	DataRangeConstAll                                 // 全部数据权限
	DataRangeConstCustom                              // 自定数据权限
)

// DataRangeMapping 数据权限范围映射
var DataRangeMapping = map[DataRangeConst]string{
	DataRangeConstOnly:      "仅本人数据权限",
	DataRangeConstDept:      "本部门数据权限",
	DataRangeConstDeptBelow: "本部门及以下数据权限",
	DataRangeConstAll:       "全部数据权限",
	DataRangeConstCustom:    "自定数据权限",
}

// DataRangeImportMapping 数据权限范围映射
var DataRangeImportMapping = map[string]DataRangeConst{
	"仅本人数据权限":    DataRangeConstOnly,
	"本部门数据权限":    DataRangeConstDept,
	"本部门及以下数据权限": DataRangeConstDeptBelow,
	"全部数据权限":     DataRangeConstAll,
	"自定数据权限":     DataRangeConstCustom,
}
