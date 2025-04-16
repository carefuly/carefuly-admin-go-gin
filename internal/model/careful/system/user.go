/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:15:18
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/system/user"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	models.CoreModels
	Username string           `gorm:"type:varchar(50);not null;unique;column:username;comment:用户名" json:"username"` // 用户名
	Password string           `gorm:"type:varchar(255);not null;column:password;comment:密码" json:"-"`               // 密码
	Name     string           `gorm:"type:varchar(50);index:idx_search;column:name;comment:姓名" json:"name"`         // 姓名
	Email    string           `gorm:"type:varchar(50);index:idx_search;column:email;comment:邮箱" json:"email"`       // 邮箱
	Mobile   string           `gorm:"type:varchar(20);index:idx_search;column:mobile;comment:电话" json:"mobile"`     // 电话
	Avatar   string           `gorm:"type:text;column:avatar;comment:头像" json:"avatar"`                             // 头像
	Gender   user.GenderConst `gorm:"type:tinyint;default:0;column:gender;comment:性别" json:"gender"`                // 性别
}

func NewUser() *User {
	return &User{}
}

func (u *User) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{})
	if err != nil {
		zap.L().Error("User表模型迁移失败", zap.Error(err))
	}
}

func (u *User) TableName() string {
	return "careful_system_users"
}

// UserPassword 用户密码表
type UserPassword struct {
	models.CoreModels
	Username string `gorm:"type:varchar(50);not null;unique;column:username;comment:用户账号" json:"username"` // 用户账号
	Password string `gorm:"type:varchar(255);not null;column:password;comment:密码" json:"password"`         // 密码
}

func NewUserPassword() *UserPassword {
	return &UserPassword{}
}

func (u *UserPassword) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&UserPassword{})
	if err != nil {
		zap.L().Error("UserPassword表模型迁移失败", zap.Error(err))
	}
}

func (u *UserPassword) TableName() string {
	return "careful_system_users_password"
}
