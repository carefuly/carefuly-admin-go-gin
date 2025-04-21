/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/16 22:29:47
 * Remark：
 */

package tools

import (
	"database/sql"
	"errors"

	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dictType"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrDictTypeUniqueIndex = errors.New("违反唯一约束")

// DictType 字典信息表
type DictType struct {
	models.CoreModels
	Name      string                `gorm:"type:varchar(50);not null;index:idx_name;column:name;comment:字典信息名称" json:"name"`                // 字典信息名称
	StrValue  sql.NullString        `gorm:"type:varchar(50);column:strValue;comment:字符串-字典信息值" json:"strValue"`                             // 字符串-字典信息值
	IntValue  sql.NullInt64         `gorm:"type:tinyint;column:intValue;comment:整型-字典信息值" json:"intValue"`                                  // 整型-字典信息值
	BoolValue sql.NullBool          `gorm:"type:boolean;column:boolValue;comment:布尔-字典信息值" json:"boolValue"`                                // 布尔-字典信息值
	DictTag   dictType.DictTagConst `gorm:"type:varchar(10);default:primary;index:idx_dict_tag;column:dictTag;comment:标签类型" json:"dictTag"` // 标签类型
	DictColor string                `gorm:"type:varchar(50);column:dictColor;comment:标签颜色" json:"dictColor"`                                // 标签颜色
	DictName  string                `gorm:"type:varchar(100);column:dictName;comment:字典名称" json:"dictName"`                                 // 字典名称
	TypeValue dict.TypeValueConst   `gorm:"type:tinyint;default:0;index:idx_type_value;column:typeValue;comment:字典类型值" json:"typeValue"`    // 字典类型值
	DictId    string                `gorm:"type:varchar(100);index:idx_dict_id;column:dict_id;comment:字典ID" json:"dict_id"`                 // 字典ID
	Dict      *Dict                 `gorm:"foreignKey:DictId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"dict"`                    // 字典信息
}

func NewDictType() *DictType {
	return &DictType{}
}

func (d *DictType) TableName() string {
	return "careful_tools_dict_type"
}

func (d *DictType) AutoMigrate(db *gorm.DB) {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&DictType{}); err != nil {
		zap.L().Error("DictType表模型迁移失败", zap.Error(err))
	}
}

// BeforeSave 在创建/更新时校验数据一致性
func (d *DictType) BeforeSave(tx *gorm.DB) error {
	query := tx.Model(&DictType{}).
		Where("dict_id = ?", d.DictId)

	// 根据类型添加不同的条件
	switch d.TypeValue {
	case 0:
		query = query.Where("strValue = ? OR name = ?", d.StrValue.String, d.Name)
	case 1:
		query = query.Where("intValue = ? OR name = ?", d.IntValue.Int64, d.Name)
	case 2:
		query = query.Where("boolValue = ? OR name = ?", d.BoolValue.Bool, d.Name)
	}

	// 排除自身（更新时用）
	if d.Id != "" {
		query = query.Where("id <> ?", d.Id)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrDictTypeUniqueIndex
	}
	return nil
}
