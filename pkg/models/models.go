/**
 * Description：
 * FileName：models.go
 * Author：CJiaの用心
 * Create：2025/3/19 22:50:16
 * Remark：
 */

package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

// CoreModels 公共模型
// 核心标准抽象模型模型,可直接继承使用
// 增加审计字段, 覆盖字段时, 字段名称请勿修改, 必须统一审计字段名称
type CoreModels struct {
	Id         string     `gorm:"type:varchar(100);primary_key;column:id;comment:主键ID" json:"id"`              // 主键ID
	Sort       int        `gorm:"type:int;default:1;column:sort;comment:显示排序" json:"sort"`                     // 显示排序
	Version    int        `gorm:"type:int;default:1;column:version;comment:版本号" json:"version"`                // 版本号
	Creator    string     `gorm:"type:varchar(100);index;column:creator;comment:创建人" json:"creator"`           // 创建人
	Modifier   string     `gorm:"type:varchar(100);index;column:modifier;comment:修改人" json:"modifier"`         // 修改人
	BelongDept string     `gorm:"type:varchar(100);index;column:belong_dept;comment:数据归属部门" json:"belongDept"` // 数据归属部门
	Status     bool       `gorm:"type:bool;index;default:true;column:status;comment:状态" json:"status"`         // 状态
	CreateTime *time.Time `gorm:"autoCreateTime;column:create_time;comment:创建时间" json:"-"`                     // 创建时间
	UpdateTime *time.Time `gorm:"autoUpdateTime;column:update_time;comment:修改时间" json:"-"`                     // 修改时间
	Remark     string     `gorm:"type:varchar(255);column:remark;comment:备注" json:"remark"`                    // 备注
}

func (c *CoreModels) BeforeCreate(tx *gorm.DB) (err error) {
	// 记录ID
	c.Id = strings.ToUpper(uuid.NewV4().String())
	// 返回异常
	return
}
