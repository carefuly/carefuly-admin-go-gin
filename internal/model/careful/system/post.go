/**
 * Description：
 * FileName：post.go
 * Author：CJiaの用心
 * Create：2025/6/13 17:11:36
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Post 岗位表
type Post struct {
	models.CoreModels
	Status bool   `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-在职 false-离职】" json:"status"`           // 状态
	Name   string `gorm:"type:varchar(100);not null;uniqueIndex:uni_post_name_code;index:idx_name;column:name;comment:岗位名称" json:"name"` // 部门名称
	Code   string `gorm:"type:varchar(100);not null;uniqueIndex:uni_post_name_code;index:idx_code;column:code;comment:岗位编码" json:"code"` // 部门编码
}

func NewPost() *Post {
	return &Post{}
}

func (p *Post) TableName() string {
	return "careful_system_post"
}

func (p *Post) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='岗位表'").AutoMigrate(&Post{})
	if err != nil {
		zap.L().Error("Post表模型迁移失败", zap.Error(err))
	}
}
