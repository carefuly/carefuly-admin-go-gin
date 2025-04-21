/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 14:41:17
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dictType"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	_import "github.com/carefuly/carefuly-admin-go-gin/pkg/utils/import"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/jsonformat"
)

var (
	ErrNotSupportedTypeValue        = tools.ErrNotSupportedTypeValue
	ErrDictTypeRecordNotFound       = tools.ErrDictTypeRecordNotFound
	ErrDictTypeNotFound             = tools.ErrDictTypeNotFound
	ErrDuplicateDictType            = tools.ErrDuplicateDictType
	ErrDictTypeVersionInconsistency = tools.ErrDictTypeVersionInconsistency
)

type DictTypeService interface {
	Create(ctx context.Context, domain domainTools.DictType) error
	Import(ctx context.Context, userId string, listMap []map[string]string) _import.ImportResult
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, domain domainTools.DictType) error
	GetById(ctx context.Context, id string) (domainTools.DictType, error)
	GetListPage(ctx context.Context, filter domainTools.DictTypeFilter) ([]domainTools.DictType, int64, error)
	GetListAll(ctx context.Context, filter domainTools.DictTypeFilter) ([]domainTools.DictType, error)
}

type dictTypeService struct {
	repo     tools.DictTypeRepository
	dictRepo tools.DictRepository
}

func NewDictTypeService(repo tools.DictTypeRepository, dictRepo tools.DictRepository) DictTypeService {
	return &dictTypeService{
		repo:     repo,
		dictRepo: dictRepo,
	}
}

// Create 创建
func (svc *dictTypeService) Create(ctx context.Context, domain domainTools.DictType) error {
	// 获取字典详情
	dict, err := svc.dictRepo.GetById(ctx, domain.DictId)
	if errors.Is(err, tools.ErrDictRecordNotFound) {
		return tools.ErrDictRecordNotFound
	}
	if err != nil {
		return err
	}
	if dict.Id == "" {
		return tools.ErrDictNotFound
	}

	// 设置DictName和TypeValue
	domain.DictName = dict.Name
	domain.TypeValue = dict.TypeValue

	// 唯一性校验
	// 逻辑较为复杂，暂时不实现，默认使用mysql唯一性约束

	return svc.repo.Create(ctx, domain)
}

// Import 导入
func (svc *dictTypeService) Import(ctx context.Context, userId string, listMap []map[string]string) _import.ImportResult {
	result := _import.ImportResult{}

	// 遍历数据
	for index, list := range listMap {
		rowNumber := index + 2

		// 数据清洗
		name := _import.CleanInput(list["字典信息名称"])

		// 字段校验
		if name == "" {
			result.AddError(rowNumber, "【字典信息名称】不能为空")
			continue
		}

		// 判断所属字典
		dict, err := svc.dictRepo.GetByName(ctx, list["所属字典"])
		if errors.Is(err, tools.ErrDictRecordNotFound) {
			result.AddError(rowNumber, "【所属字典】不存在")
			continue
		}
		if err != nil {
			result.AddError(rowNumber, "【所属字典】不存在")
			continue
		}
		if dict.Id == "" {
			result.AddError(rowNumber, "【所属字典】不存在")
			continue
		}

		var intValue int
		var boolValue bool
		// 类型转换
		if dict.TypeValue == 1 {
			if _import.CleanInput(list["整型值"]) == "" {
				result.AddError(rowNumber, fmt.Sprintf("【整型值】不能为空"))
				continue
			}
			intValue, _ = strconv.Atoi(list["整型值"])
		} else if dict.TypeValue == 2 {
			boolValue, err = dictType.ConvertBoolValueImport(list["布尔值"])
			if err != nil {
				result.AddError(rowNumber, fmt.Sprintf("【布尔值】转换失败：%s", err.Error()))
				continue
			}
		}
		dictTag, err := dictType.ConvertDictTagImport(list["标签类型"])
		if err != nil {
			result.AddError(rowNumber, fmt.Sprintf("【标签类型】转换失败：%s", err.Error()))
			continue
		}

		// 处理其他字段
		var sort int
		if list["排序"] == "" {
			sort = 1
		} else {
			sort, _ = strconv.Atoi(list["排序"])
		}

		// 构建领域模型
		domain := domainTools.DictType{
			DictType: modelTools.DictType{
				CoreModels: models.CoreModels{
					Creator:  userId,
					Modifier: userId,
					Sort:     sort,
					Remark:   list["备注"],
				},
				Name:      name,
				DictTag:   dictTag,
				DictColor: list["标签颜色"],
				DictName:  dict.Name,
				TypeValue: dict.TypeValue,
				DictId:    dict.Id,
			},
			StrValue:  list["字符串值"],
			IntValue:  int64(intValue),
			BoolValue: boolValue,
		}

		jsonformat.FormatJsonPrint(domain)

		// 创建记录
		if err = svc.repo.Create(ctx, domain); err != nil {
			result.AddError(rowNumber, "创建失败："+err.Error())
			continue
		}

		result.SuccessCount++
	}

	return result
}

// Delete 删除
func (svc *dictTypeService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return tools.ErrDictTypeNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *dictTypeService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *dictTypeService) Update(ctx context.Context, id string, domain domainTools.DictType) error {
	// 获取字典详情
	dict, err := svc.dictRepo.GetById(ctx, domain.DictId)
	if errors.Is(err, tools.ErrDictRecordNotFound) {
		return tools.ErrDictRecordNotFound
	}
	if err != nil {
		return err
	}
	if dict.Id == "" {
		return tools.ErrDictNotFound
	}

	// 设置DictName和TypeValue
	domain.DictName = dict.Name
	domain.TypeValue = dict.TypeValue

	// 唯一性校验
	// 逻辑较为复杂，暂时不实现，默认使用mysql唯一性约束

	_, err = svc.repo.Update(ctx, id, domain)

	switch {
	case err == nil:
		return err
	case errors.Is(err, tools.ErrDictTypeNotFound):
		return tools.ErrDictTypeNotFound
	case errors.Is(err, tools.ErrDictTypeVersionInconsistency):
		return tools.ErrDictTypeVersionInconsistency
	default:
		return err
	}
}

// GetById 获取详情
func (svc *dictTypeService) GetById(ctx context.Context, id string) (domainTools.DictType, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if errors.Is(err, tools.ErrDictTypeRecordNotFound) {
		return domain, tools.ErrDictTypeRecordNotFound
	}
	if err != nil {
		return domain, err
	}
	if domain.Id == "" {
		return domain, tools.ErrDictTypeNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *dictTypeService) GetListPage(ctx context.Context, filter domainTools.DictTypeFilter) ([]domainTools.DictType, int64, error) {
	return svc.repo.GetListPage(ctx, filter)
}

// GetListAll 查询所有列表
func (svc *dictTypeService) GetListAll(ctx context.Context, filter domainTools.DictTypeFilter) ([]domainTools.DictType, error) {
	return svc.repo.GetListAll(ctx, filter)
}
