/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/14 16:25:52
 * Remark：
 */

package tools

import (
	"context"
	"errors"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	repositoryTools "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/tools"
	"github.com/go-sql-driver/mysql"
	"strings"
)

var (
	ErrDictNotFound             = repositoryTools.ErrDictNotFound
	ErrDictNameDuplicate        = repositoryTools.ErrDictNameDuplicate
	ErrDictCodeDuplicate        = repositoryTools.ErrDictCodeDuplicate
	ErrDictDuplicate            = repositoryTools.ErrDictDuplicate
	ErrDictVersionInconsistency = repositoryTools.ErrDictVersionInconsistency
)

type DictService interface {
	Create(ctx context.Context, domain domainTools.Dict) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, domain domainTools.Dict) error

	GetById(ctx context.Context, id string) (domainTools.Dict, error)
	GetByName(ctx context.Context, name string) (domainTools.Dict, error)
	GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error)
	GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error)
}

type dictService struct {
	repo repositoryTools.DictRepository
}

func NewDictService(repo repositoryTools.DictRepository) DictService {
	return &dictService{
		repo: repo,
	}
}

// Create 创建
func (svc *dictService) Create(ctx context.Context, domain domainTools.Dict) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, "")
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrDictNameDuplicate
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, "")
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrDictCodeDuplicate
	}

	if err := svc.repo.Create(ctx, domain); err != nil {
		// 分析具体冲突字段
		if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
			switch field {
			case "name":
				return repositoryTools.ErrDictNameDuplicate
			case "code":
				return repositoryTools.ErrDictCodeDuplicate
			case "all":
				return repositoryTools.ErrDictDuplicate
			default:
				return err
			}
		}
	}

	return nil
}

// Delete 删除
func (svc *dictService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositoryTools.ErrDictNotFound
	}
	return err
}

// Update 更新
func (svc *dictService) Update(ctx context.Context, domain domainTools.Dict) error {
	exists, err := svc.repo.CheckExistByName(ctx, domain.Name, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrDictNameDuplicate
	}

	exists, err = svc.repo.CheckExistByCode(ctx, domain.Code, domain.Id)
	if err != nil {
		return err
	}
	if exists {
		return repositoryTools.ErrDictCodeDuplicate
	}

	if err := svc.repo.Update(ctx, domain); err != nil {
		// 分析具体冲突字段
		if field, isDuplicate := svc.IsDuplicateEntryError(err); isDuplicate {
			switch field {
			case "name":
				return repositoryTools.ErrDictNameDuplicate
			case "code":
				return repositoryTools.ErrDictCodeDuplicate
			case "all":
				return repositoryTools.ErrDictDuplicate
			default:
				return err
			}
		}
	}

	return nil
}

// GetById 获取详情
func (svc *dictService) GetById(ctx context.Context, id string) (domainTools.Dict, error) {
	domain, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositoryTools.ErrDictNotFound) {
			return domain, repositoryTools.ErrDictNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositoryTools.ErrDictNotFound
	}
	return domain, err
}

// GetByName 根据name获取详情
func (svc *dictService) GetByName(ctx context.Context, name string) (domainTools.Dict, error) {
	domain, err := svc.repo.GetByName(ctx, name)
	if err != nil {
		if errors.Is(err, repositoryTools.ErrDictNotFound) {
			return domain, repositoryTools.ErrDictNotFound
		}
		return domain, err
	}
	if domain.Id == "" {
		return domain, repositoryTools.ErrDictNotFound
	}
	return domain, err
}

// GetListPage 分页查询列表
func (svc *dictService) GetListPage(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, int64, error) {
	return svc.repo.GetListPage(ctx, filters)
}

// GetListAll 查询所有列表
func (svc *dictService) GetListAll(ctx context.Context, filters domainTools.DictFilter) ([]domainTools.Dict, error) {
	return svc.repo.GetListAll(ctx, filters)
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
