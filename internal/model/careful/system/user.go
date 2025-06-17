/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/5/12 14:25:26
 * Remark：
 */

package system

import (
	"fmt"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/user"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	models.CoreModels
	Status      bool             `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"` // 状态
	Username    string           `gorm:"type:varchar(50);not null;unique;column:username;comment:用户名" json:"username"`                        // 用户名
	Password    string           `gorm:"type:varchar(255);not null;column:password;comment:密码" json:"-"`                                      // 密码
	PasswordStr string           `gorm:"type:varchar(255);not null;column:password_str;comment:明文密码" json:"-"`                                // 明文密码
	UserType    user.TypeConst   `gorm:"type:tinyint;default:1;column:user_type;comment:用户类型" json:"userType"`                                // 用户类型
	Name        string           `gorm:"type:varchar(50);index:idx_search;column:name;comment:姓名" json:"name"`                                // 姓名
	Gender      user.GenderConst `gorm:"type:tinyint;default:1;column:gender;comment:性别" json:"gender"`                                       // 性别
	Email       string           `gorm:"type:varchar(50);index:idx_search;column:email;comment:邮箱" json:"email"`                              // 邮箱
	Mobile      string           `gorm:"type:varchar(20);index:idx_search;column:mobile;comment:电话" json:"mobile"`                            // 电话
	Avatar      string           `gorm:"type:text;column:avatar;comment:头像" json:"avatar"`                                                    // 头像
	PostIDs     []string         `gorm:"-" json:"post_ids"`                                                                                   // 忽略GORM处理
	RoleIDs     []string         `gorm:"-" json:"role_ids"`                                                                                   // 忽略GORM处理
	DeptId      string           `gorm:"type:varchar(100);index;column:dept_id;comment:部门ID" json:"dept_id"`                                  // 部门ID
	Dept        *Dept            `gorm:"foreignKey:DeptId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"dept"`                         // 部门
	Post        []*Post          `gorm:"many2many:careful_system_users_post;"`                                                                // 关联岗位
	Role        []*Role          `gorm:"many2many:careful_system_users_role;"`                                                                // 关联角色
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "careful_system_users"
}

func (u *User) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='用户表'").AutoMigrate(&User{})
	if err != nil {
		zap.L().Error("User表模型迁移失败", zap.Error(err))
	}

	// 迁移中间表并设置备注
	migrateManyToManyTable(db, "careful_system_users_post", "用户-关联岗位表")
	migrateManyToManyTable(db, "careful_system_users_role", "用户-关联角色表")
}

// 迁移many2many中间表并设置表备注
func (u *User) migrateManyToManyTable(db *gorm.DB, tableName string, comment string) {
	err := db.Exec(fmt.Sprintf(
		"ALTER TABLE %s COMMENT = '%s'",
		tableName,
		comment,
	)).Error

	if err != nil {
		zap.L().Error(fmt.Sprintf("%s表备注设置失败", tableName), zap.Error(err))
	}
}
