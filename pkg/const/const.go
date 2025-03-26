/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/3/26 11:44:20
 * Remark：
 */

package _const

type GenderConst int

const (
	GenderMale    GenderConst = iota // 男
	GenderFemale                     // 女
	GenderUnknown                    // 未知
)
