/**
 * Description：
 * FileName：menu_button.go
 * Author：CJiaの用心
 * Create：2025/6/8 22:14:07
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/menu"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MenuButton 菜单权限表
type MenuButton struct {
	models.CoreModels
	Status bool             `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Name   string           `gorm:"type:varchar(64);not null;column:name;comment:名称" json:"name"`                                        // 名称
	Code   string           `gorm:"type:varchar(64);not null;column:code;comment:权限值" json:"code"`                                       // 权限值
	Api    string           `gorm:"type:varchar(255);not null;column:api;comment:接口地址" json:"api"`                                       // 接口地址
	Method menu.MethodConst `gorm:"type:varchar(16);column:method;comment:请求方式" json:"method"`                                           // 请求方式
	MenuId string           `gorm:"type:varchar(64);column:menu_id;comment:关联菜单" json:"menu_id"`                                         // 关联菜单
	Menu   Menu             `gorm:"foreignKey:menu_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"menu"`                        // 菜单
}

func NewMenuButton() *MenuButton {
	return &MenuButton{}
}

func (m *MenuButton) TableName() string {
	return "careful_system_menu_button"
}

func (m *MenuButton) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='菜单权限表'").AutoMigrate(&MenuButton{})
	if err != nil {
		zap.L().Error("MenuButton表模型迁移失败", zap.Error(err))
	}
}
