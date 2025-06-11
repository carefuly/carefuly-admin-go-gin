/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/10 16:53:57
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/role"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Role 角色表
type Role struct {
	models.CoreModels
	Status     bool                `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Name       string              `gorm:"type:varchar(64);not null;column:name;comment:角色名称" json:"name"`                                      // 角色名称
	Code       string              `gorm:"type:varchar(64);not null;uniqueIndex;column:code;comment:角色编码" json:"code"`                          // 角色编码
	DataRange  role.DataRangeConst `gorm:"type:tinyint;default:1;column:data_range;comment:数据权限范围" json:"data_range"`                           // 数据权限范围
	Dept       []Dept              `gorm:"many2many:role_dept;"`                                                                                // 数据权限-关联部门
	Menu       []Menu              `gorm:"many2many:role_menu;"`                                                                                // 数据权限-关联菜单
	permission []MenuButton        `gorm:"many2many:role_menu_button;"`                                                                         // 数据权限-关联菜单的接口按钮
	column     []MenuColumn        `gorm:"many2many:role_menu_column;"`                                                                         // 数据权限-列表权限
}

func NewRole() *Role {
	return &Role{}
}

func (r *Role) TableName() string {
	return "careful_system_role"
}

func (r *Role) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='角色表'").AutoMigrate(&Role{})
	if err != nil {
		zap.L().Error("Role表模型迁移失败", zap.Error(err))
	}
}
