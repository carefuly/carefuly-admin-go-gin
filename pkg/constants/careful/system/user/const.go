/**
 * Description：
 * FileName：const.go
 * Author：CJiaの用心
 * Create：2025/6/17 14:17:20
 * Remark：
 */

package user

type TypeConst int

const (
	TypeConstAdminUser TypeConst = iota + 1 // 后台用户
	TypeConstFrontUser                      // 前台用户
)

type GenderConst int

const (
	GenderConstMale   GenderConst = iota + 1 // 男
	GenderConstFemale                        // 女
	GenderConstSecret                        // 保密
)
