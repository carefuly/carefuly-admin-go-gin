/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/8 22:14:32
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MenuColumn 菜单数据列表
type MenuColumn struct {
	models.CoreModels
	Status bool   `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Title  string `gorm:"type:varchar(64);not null;column:title;comment:标题" json:"title"`                                      // 标题
	Field  string `gorm:"type:varchar(64);not null;column:field;comment:字段名" json:"field"`                                     // 字段名
	Width  int    `gorm:"type:int;default:150;column:width;comment:宽度" json:"width"`                                           // 宽度
	MenuId string `gorm:"type:varchar(64);column:menu_id;comment:关联菜单" json:"menu_id"`                                         // 关联菜单
	Menu   Menu   `gorm:"foreignKey:menu_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"menu"`                        // 菜单
}

func NewMenuColumn() *MenuColumn {
	return &MenuColumn{}
}

func (m *MenuColumn) TableName() string {
	return "careful_system_menu_column"
}

func (m *MenuColumn) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='菜单数据列表'").AutoMigrate(&MenuColumn{})
	if err != nil {
		zap.L().Error("MenuColumn表模型迁移失败", zap.Error(err))
	}
}
