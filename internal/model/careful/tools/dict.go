/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/14 16:45:52
 * Remark：
 */

package tools

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Dict 字典表
type Dict struct {
	models.CoreModels
	Name      string              `gorm:"type:varchar(100);not null;unique;index:idx_name;column:name;comment:字典名称" json:"name"`       // 字典名称
	Code      string              `gorm:"type:varchar(100);not null;unique;index:idx_code;column:code;comment:字典编码" json:"code"`       // 字典编码
	Type      dict.TypeConst      `gorm:"type:tinyint;default:0;index:idx_type;column:type;comment:字典类型" json:"type"`                  // 字典类型
	TypeValue dict.TypeValueConst `gorm:"type:tinyint;default:0;index:idx_type_value;column:typeValue;comment:字典类型值" json:"typeValue"` // 字典类型值
}

func NewDict() *Dict {
	return &Dict{}
}

func (d *Dict) TableName() string {
	return "careful_tools_dict"
}

func (d *Dict) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Dict{})
	if err != nil {
		zap.L().Error("Dict表模型迁移失败", zap.Error(err))
	}
}
