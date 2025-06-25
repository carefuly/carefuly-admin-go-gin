/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 15:49:28
 * Remark：
 */

package system

import (
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/menu"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Menu 菜单表
type Menu struct {
	models.CoreModels
	Status      bool           `gorm:"type:boolean;index:idx_status;default:true;column:status;comment:状态【true-启用 false-停用】" json:"status"`                // 状态
	Type        menu.TypeConst `gorm:"type:tinyint;default:2;uniqueIndex:uni_menu_title_unique;column:type;comment:菜单类型" json:"type"`                      // 菜单类型
	Icon        string         `gorm:"type:varchar(64);column:icon;comment:菜单图标" json:"icon"`                                                              // 菜单图标
	Title       string         `gorm:"type:varchar(64);not null;uniqueIndex:uni_menu_title_unique;index:idx_title;column:title;comment:菜单标题" json:"title"` // 菜单标题
	Name        string         `gorm:"type:varchar(64);not null;column:name;comment:组件名称" json:"name"`                                                     // 组件名称
	Component   string         `gorm:"type:varchar(128);column:component;comment:组件地址" json:"component"`                                                   // 组件地址
	Path        string         `gorm:"type:varchar(128);not null;column:path;comment:路由地址" json:"path"`                                                    // 路由地址
	Redirect    string         `gorm:"type:varchar(128);column:redirect;comment:重定向地址" json:"redirect"`                                                    // 重定向地址
	IsHide      bool           `gorm:"type:boolean;default:false;column:isHide;comment:是否隐藏" json:"isHide"`                                                // 是否隐藏
	IsLink      string         `gorm:"type:varchar(255);column:isLink;comment:是否外链【不填写默认没有外链】" json:"isLink"`                                              // 是否外链
	IsKeepAlive bool           `gorm:"type:boolean;default:false;column:isKeepAlive;comment:是否页面缓存" json:"isKeepAlive"`                                    // 是否页面缓存
	IsFull      bool           `gorm:"type:boolean;default:false;column:isFull;comment:是否缓存全屏" json:"isFull"`                                              // 是否缓存全屏
	IsAffix     bool           `gorm:"type:boolean;default:false;column:isAffix;comment:是否缓存固定路由" json:"isAffix"`                                          // 是否缓存固定路由
	ParentID    string         `gorm:"type:varchar(100);uniqueIndex:uni_menu_title_unique;column:parent_id;comment:上级菜单" json:"parent_id"`                 // 上级菜单
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) TableName() string {
	return "careful_system_menu"
}

func (m *Menu) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB,COMMENT='菜单表'").AutoMigrate(&Menu{})
	if err != nil {
		zap.L().Error("Menu表模型迁移失败", zap.Error(err))
	}
}
