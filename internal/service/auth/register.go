/**
 * Description：
 * FileName：register.go
 * Author：CJiaの用心
 * Create：2025/3/26 13:06:05
 * Remark：
 */

package service

import (
	"context"
	"errors"
	domain "github.com/carefuly/carefuly-admin-go-gin/internal/domain/auth"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/system"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = repository.ErrDuplicateUsername
)

type RegisterService interface {
	Register(ctx context.Context, u domain.Register) error
}

type registerService struct {
	repo         repository.UserRepository
	userPassRepo repository.UserPassWordRepository
}

func NewRegisterService(repo repository.UserRepository, userPassRepo repository.UserPassWordRepository) RegisterService {
	return &registerService{
		repo:         repo,
		userPassRepo: userPassRepo,
	}
}

func (svc *registerService) Register(ctx context.Context, u domain.Register) error {
	exists, err := svc.repo.ExistsByUserName(ctx, u.Username)
	if err != nil {
		return err
	}
	if exists {
		return repository.ErrDuplicateUsername
	}

	text := u.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	err = svc.repo.Register(ctx, u)
	if svc.IsDuplicateEntryError(err) {
		return repository.ErrDuplicateUsername
	}
	if err != nil {
		return err
	}

	// 保存用户和密码
	exists, err = svc.userPassRepo.ExistsByUserName(ctx, u.Username)
	if err != nil {
		return err
	}
	if exists {
		return repository.ErrDuplicateUsername
	}

	err = svc.userPassRepo.Create(ctx, domain.Register{
		Username: u.Username,
		Password: text,
	})
	if svc.IsDuplicateEntryError(err) {
		return repository.ErrDuplicateUsername
	}

	return err
}

func (svc *registerService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
