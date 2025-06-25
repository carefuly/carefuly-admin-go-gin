/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/14 11:22:27
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Dict 字典表
type Dict struct {
	models.CoreModels
	Status    bool                `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Name      string              `gorm:"type:varchar(100);not null;uniqueIndex;index:idx_name;column:name;comment:字典名称" json:"name"`          // 字典名称
	Code      string              `gorm:"type:varchar(100);not null;uniqueIndex;index:idx_code;column:code;comment:字典编码" json:"code"`          // 字典编码
	Type      dict.TypeConst      `gorm:"type:tinyint;default:1;index:idx_type;column:type;comment:字典分类" json:"type"`                          // 字典分类
	ValueType dict.TypeValueConst `gorm:"type:tinyint;default:1;index:idx_value_type;column:valueType;comment:数据类型" json:"valueType"`          // 数据类型
}

func NewDict() *Dict {
	return &Dict{}
}

func (d *Dict) TableName() string {
	return "careful_tools_dict"
}

func (d *Dict) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='字典表'").AutoMigrate(&Dict{})
	if err != nil {
		zap.L().Error("Dict表模型迁移失败", zap.Error(err))
	}
}
