/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/14 15:55:03
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Bucket 存储桶表
type Bucket struct {
	models.CoreModels
	Status bool   `gorm:"type:boolean;index:idx_status;default:false;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Name   string `gorm:"type:varchar(50);not null;uniqueIndex;index:idx_name;column:name;comment:名称" json:"name"`              // 名称
	Code   string `gorm:"type:varchar(50);not null;uniqueIndex;index:idx_code;column:code;comment:编码" json:"code"`              // 编码
	Size   int    `gorm:"type:tinyint;default:1;column:size;comment:存储桶大小(GB)" json:"size"`                                     // 存储桶大小(GB)
}

func NewBucket() *Bucket {
	return &Bucket{}
}

func (d *Bucket) TableName() string {
	return "careful_tools_bucket"
}

func (d *Bucket) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='存储桶表'").AutoMigrate(&Bucket{})
	if err != nil {
		zap.L().Error("Bucket表模型迁移失败", zap.Error(err))
	}
}
