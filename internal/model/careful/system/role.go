/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/10 16:53:57
 * Remark：
 */

package system

import (
	"fmt"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/role"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Role 角色表
type Role struct {
	models.CoreModels
	Status        bool                `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Name          string              `gorm:"type:varchar(64);not null;index:idx_name;column:name;comment:角色名称" json:"name"`                       // 角色名称
	Code          string              `gorm:"type:varchar(64);not null;uniqueIndex;column:code;comment:角色编码" json:"code"`                          // 角色编码
	DataRange     role.DataRangeConst `gorm:"type:tinyint;default:1;column:data_range;comment:数据权限范围" json:"data_range"`                           // 数据权限范围
	DeptIDs       []string            `gorm:"-" json:"dept_ids"`                                                                                   // 忽略GORM处理，只用于接收参数
	MenuIDs       []string            `gorm:"-" json:"menu_ids"`                                                                                   // 忽略GORM处理
	MenuButtonIDs []string            `gorm:"-" json:"menu_button_ids"`                                                                            // 忽略GORM处理
	MenuColumnIDs []string            `gorm:"-" json:"menu_column_ids"`                                                                            // 忽略GORM处理
	Dept          []*Dept             `gorm:"many2many:careful_system_role_dept;" json:"dept"`                                                     // 数据权限-关联部门
	Menu          []*Menu             `gorm:"many2many:careful_system_role_menu;" json:"menu"`                                                     // 数据权限-关联菜单
	MenuButton    []*MenuButton       `gorm:"many2many:careful_system_role_menu_button;" json:"menuButton"`                                        // 数据权限-关联菜单的接口按钮
	MenuColumn    []*MenuColumn       `gorm:"many2many:careful_system_role_menu_column;" json:"menuColumn"`                                        // 数据权限-列表权限
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

	// 迁移中间表并设置备注
	migrateManyToManyTable(db, "careful_system_role_dept", "角色-部门关联表")
	migrateManyToManyTable(db, "careful_system_role_menu", "角色-菜单关联表")
	migrateManyToManyTable(db, "careful_system_role_menu_button", "角色-菜单按钮关联表")
	migrateManyToManyTable(db, "careful_system_role_menu_column", "角色-菜单列关联表")
}

// 迁移many2many中间表并设置表备注
func migrateManyToManyTable(db *gorm.DB, tableName string, comment string) {
	err := db.Exec(fmt.Sprintf(
		"ALTER TABLE %s COMMENT = '%s'",
		tableName,
		comment,
	)).Error

	if err != nil {
		zap.L().Error(fmt.Sprintf("%s表备注设置失败", tableName), zap.Error(err))
	}
}
