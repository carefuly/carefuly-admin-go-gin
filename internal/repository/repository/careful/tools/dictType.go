/**
 * Description：
 * FileName：dictType.go
 * Author：CJiaの用心
 * Create：2025/5/23 16:50:37
 * Remark：
 */

package tools

import daoTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/dao/careful/tools"

var (
	ErrDictTypeNotFound             = daoTools.ErrDictTypeNotFound
	ErrDictTypeDuplicate            = daoTools.ErrDictTypeDuplicate
	ErrDictTypeVersionInconsistency = daoTools.ErrDictTypeVersionInconsistency
)
