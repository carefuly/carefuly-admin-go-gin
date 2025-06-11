/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/5/15 15:54:38
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Dept 部门表
type Dept struct {
	models.CoreModels
	Status   bool   `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"`                        // 状态
	Name     string `gorm:"type:varchar(100);not null;uniqueIndex:uni_dept_name_code_parent;index:idx_name;column:name;comment:部门名称" json:"name"`       // 部门名称
	Code     string `gorm:"type:varchar(100);not null;uniqueIndex:uni_dept_name_code_parent;index:idx_code;column:code;comment:部门编码" json:"code"`       // 部门编码
	Owner    string `gorm:"type:varchar(32);column:owner;comment:负责人" json:"owner"`                                                                     // 负责人
	Phone    string `gorm:"type:varchar(32);column:phone;comment:联系电话" json:"phone"`                                                                    // 联系电话
	Email    string `gorm:"type:varchar(32);column:email;comment:邮箱" json:"email"`                                                                      // 邮箱
	ParentID string `gorm:"type:varchar(100);uniqueIndex:uni_dept_name_code_parent;index:idx_parent_id;column:parent_id;comment:上级部门" json:"parent_id"` // 上级部门
}

func NewDept() *Dept {
	return &Dept{}
}

func (d *Dept) TableName() string {
	return "careful_system_dept"
}

func (d *Dept) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='部门表'").AutoMigrate(&Dept{})
	if err != nil {
		zap.L().Error("Dept表模型迁移失败", zap.Error(err))
	}
}
