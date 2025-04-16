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

// DictType 字典信息表
type DictType struct {
	models.CoreModels
	Name      string                `gorm:"type:varchar(50);not null;index:idx_name;uniqueIndex:uni_dict_name_value;column:name;comment:字典信息名称" json:"name"` // 字典信息名称
	StrValue  sql.NullString        `gorm:"type:varchar(50);uniqueIndex:uni_dict_name_value;column:strValue;comment:字符串-字典信息值" json:"strValue"`              // 字符串-字典信息值
	IntValue  sql.NullInt64         `gorm:"type:tinyint;uniqueIndex:uni_dict_name_value;column:intValue;comment:整型-字典信息值" json:"intValue"`                   // 整型-字典信息值
	BoolValue sql.NullBool          `gorm:"type:boolean;column:boolValue;comment:布尔-字典信息值" json:"boolValue"`                                                 // 布尔-字典信息值
	DictTag   dictType.DictTagConst `gorm:"type:varchar(10);default:primary;index:idx_dict_tag;column:dictType;comment:标签类型" json:"dictTag"`                 // 标签类型
	DictColor string                `gorm:"type:varchar(50);column:dictColor;comment:标签颜色" json:"dictColor"`                                                 // 标签颜色
	DictName  string                `gorm:"type:varchar(100);column:dictName;comment:字典名称" json:"dictName"`                                                  // 字典名称
	TypeValue dict.TypeValueConst   `gorm:"type:tinyint;default:0;index:idx_type_value;column:typeValue;comment:字典类型值" json:"typeValue"`                     // 字典类型值
	DictId    string                `gorm:"type:varchar(100);index:idx_dict_id;uniqueIndex:uni_dict_name_value;column:dictId;comment:字典ID" json:"dictId"`    // 字典ID
	Dict      *Dict                 `gorm:"foreignKey:DictId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"dict"`                                     // 字典信息
}

func NewDictType() *DictType {
	return &DictType{}
}

func (d *DictType) TableName() string {
	return "careful_tools_dict_type"
}

func (d *DictType) AutoMigrate(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&DictType{})
	if err != nil {
		zap.L().Error("DictType表模型迁移失败", zap.Error(err))
	}
}

// BeforeSave 在创建/更新时校验数据一致性
func (d *DictType) BeforeSave(tx *gorm.DB) error {
	switch d.TypeValue {
	case dict.TypeValueConst0:
		if !d.StrValue.Valid {
			return errors.New("字符串类型必须设置StrValue")
		}
	case dict.TypeValueConst1:
		if !d.IntValue.Valid {
			return errors.New("整型类型必须设置IntValue")
		}
	case dict.TypeValueConst2:
		if !d.BoolValue.Valid {
			return errors.New("布尔类型必须设置BoolValue")
		}
	}
	return nil
}
