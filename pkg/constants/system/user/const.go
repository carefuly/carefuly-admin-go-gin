/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/4/16 15:38:24
 * Remark：
 */

package user

type GenderConst int // 性别

const (
	GenderMale          GenderConst = iota // 男
	GenderFemale                           // 女
	GenderUnknown                          // 未知
	GenderMaleNorFemale                    // 不男不女
)
