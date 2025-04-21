/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/15 00:19:52
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	"fmt"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	_import "github.com/carefuly/carefuly-admin-go-gin/pkg/utils/import"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/jsonformat"
	"github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var (
	ErrDictRecordNotFound       = tools.ErrDictRecordNotFound
	ErrDictNotFound             = tools.ErrDictNotFound
	ErrDuplicateDict            = tools.ErrDuplicateDict
	ErrDuplicateDictName        = tools.ErrDuplicateDictName
	ErrDuplicateDictCode        = tools.ErrDuplicateDictCode
	ErrDictVersionInconsistency = tools.ErrDictVersionInconsistency
)

type DictService interface {
	Create(ctx context.Context, domain domainTools.Dict) error
	Import(ctx context.Context, userId string, listMap []map[string]string) _import.ImportResult
	Delete(ctx context.Context, id string) error
	BatchDelete(ctx context.Context, ids []string) error
	Update(ctx context.Context, id string, domain domainTools.Dict) error
	GetById(ctx context.Context, id string) (domainTools.Dict, error)
	GetListPage(ctx context.Context, filter domainTools.DictFilter) ([]domainTools.Dict, int64, error)
	GetListAll(ctx context.Context, filter domainTools.DictFilter) ([]domainTools.Dict, error)
}

type dictService struct {
	repo tools.DictRepository
}

func NewDictService(repo tools.DictRepository) DictService {
	return &dictService{repo: repo}
}

// Create 创建
func (svc *dictService) Create(ctx context.Context, domain domainTools.Dict) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return tools.ErrDuplicateDictName
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return tools.ErrDuplicateDictCode
	}

	err = svc.repo.Create(ctx, domain)
	if err == nil {
		return err
	}

	// 分析具体冲突字段
	if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
		switch field {
		case "name":
			return tools.ErrDuplicateDictName
		case "code":
			return tools.ErrDuplicateDictCode
		case "all":
			return tools.ErrDuplicateDict
		default:
			return err
		}
	}

	return err
}

// Import 导入
func (svc *dictService) Import(ctx context.Context, userId string, listMap []map[string]string) _import.ImportResult {
	result := _import.ImportResult{}

	// 遍历数据
	for index, list := range listMap {
		rowNumber := index + 2

		// 数据清洗
		name := _import.CleanInput(list["字典名称"])
		code := _import.CleanInput(list["字典编码"])

		// 字段校验
		if name == "" {
			result.AddError(rowNumber, "【字典名称】不能为空")
			continue
		}
		if code == "" {
			result.AddError(rowNumber, "【字典编码】不能为空")
			continue
		}

		// 唯一性校验
		exists, err := svc.repo.CheckExistByName(ctx, name, "")
		if err != nil {
			result.AddError(rowNumber, fmt.Sprintf("检查【字典名称：%s】唯一性失败：%s", name, err.Error()))
			continue
		}
		if exists {
			result.AddError(rowNumber, fmt.Sprintf("字典名称【%s】已存在", name))
			continue
		}
		exists, err = svc.repo.CheckExistByCode(ctx, code, "")
		if err != nil {
			result.AddError(rowNumber, fmt.Sprintf("检查【字典编码：%s】唯一性失败：%s", code, err.Error()))
			continue
		}
		if exists {
			result.AddError(rowNumber, fmt.Sprintf("字典编码【%s】已存在", code))
			continue
		}

		// 类型转换
		dictType, err := dict.ConvertDictTypeImport(list["字典类型"])
		if err != nil {
			result.AddError(rowNumber, fmt.Sprintf("【字典类型】转换失败：%s", err.Error()))
			continue
		}
		typeValue, err := dict.ConvertDictTypeValueImport(list["字典类型值"])
		if err != nil {
			result.AddError(rowNumber, fmt.Sprintf("【字典类型值】转换失败：%s", err.Error()))
			continue
		}

		// 处理字段
		var sort int
		if list["排序"] == "" {
			sort = 1
		} else {
			sort, _ = strconv.Atoi(list["排序"])
		}

		// 构建领域模型
		domain := domainTools.Dict{
			Dict: modelTools.Dict{
				CoreModels: models.CoreModels{
					Creator:  userId,
					Modifier: userId,
					Sort:     sort,
					Remark:   list["备注"],
				},
				Name:      name,
				Code:      code,
				Type:      dictType,
				TypeValue: typeValue,
			},
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
func (svc *dictService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return tools.ErrDictNotFound
	}
	return err
}

// BatchDelete 批量删除
func (svc *dictService) BatchDelete(ctx context.Context, ids []string) error {
	return svc.repo.BatchDelete(ctx, ids)
}

// Update 更新
func (svc *dictService) Update(ctx context.Context, id string, domain domainTools.Dict) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, id)
	if err != nil {
		return err
	}
	if exists {
		return tools.ErrDuplicateDictName
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, id)
	if err != nil {
		return err
	}
	if exists {
		return tools.ErrDuplicateDictCode
	}

	_, err = svc.repo.Update(ctx, id, domain)
	if err == nil {
		return err
	}

	// 分析具体冲突字段
	if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
		switch field {
		case "name":
			return tools.ErrDuplicateDictName
		case "code":
			return tools.ErrDuplicateDictCode
		case "all":
			return tools.ErrDuplicateDict
		default:
			switch {
			case errors.Is(err, tools.ErrDictNotFound):
				return tools.ErrDictNotFound
			case errors.Is(err, tools.ErrDictVersionInconsistency):
				return tools.ErrDictVersionInconsistency
			default:
				return err
			}
		}
	}
	return err
}

// GetById 获取详情
func (svc *dictService) GetById(ctx context.Context, id string) (domainTools.Dict, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if errors.Is(err, tools.ErrDictRecordNotFound) {
		return domain, tools.ErrDictRecordNotFound
	}
	if err != nil {
		return domain, err
	}
	if domain.Id == "" {
		return domain, tools.ErrDictNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *dictService) GetListPage(ctx context.Context, filter domainTools.DictFilter) ([]domainTools.Dict, int64, error) {
	return svc.repo.GetListPage(ctx, filter)
}

// GetListAll 获取所有列表
func (svc *dictService) GetListAll(ctx context.Context, filter domainTools.DictFilter) ([]domainTools.Dict, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 分析错误消息中的索引名
func (svc *dictService) IsDuplicateEntryError(err error) (string, bool) {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return "", false
	}

	// MySQL 错误码 1062 表示唯一冲突
	if mysqlErr.Number != 1062 {
		return "", false
	}

	// 分析错误消息中的索引名
	switch {
	case strings.Contains(mysqlErr.Message, "uni_careful_tools_dict_name"):
		return "name", true
	case strings.Contains(mysqlErr.Message, "uni_careful_tools_dict_code"):
		return "code", true
	default:
		return "all", true // 未知唯一键冲突
	}
}
