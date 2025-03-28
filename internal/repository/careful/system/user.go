/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:28:29
 * Remark：
 */

package system

import (
	"context"
	"github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/auth"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
)

var (
	ErrDuplicateUsername = system.ErrDuplicateUsername
	ErrUserNotFound      = system.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u auth.Register) error

	FindByUserName(ctx context.Context, username string) (domainSystem.User, error)

	ExistsByUserName(ctx context.Context, username string) (bool, error)
}

type userRepository struct {
	dao system.UserDAO
}

func NewUserRepository(dao system.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (repo *userRepository) Create(ctx context.Context, u auth.Register) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *userRepository) FindByUserName(ctx context.Context, username string) (domainSystem.User, error) {
	user, err := repo.dao.FindByUserName(ctx, username)
	if err != nil {
		return domainSystem.User{}, err
	}
	return repo.toDomain(user), err
}

func (repo *userRepository) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	return repo.dao.ExistsByUserName(ctx, username)
}

func (repo *userRepository) toEntity(d auth.Register) modelSystem.User {
	return modelSystem.User{
		Username: d.Username,
		Password: d.Password,
	}
}

func (repo *userRepository) toDomain(u *modelSystem.User) domainSystem.User {
	return domainSystem.User{
		User:       *u,
		CreateTime: u.CoreModels.CreateTime.Format("2006-01-02 15:04:05.000"),
		UpdateTime: u.CoreModels.UpdateTime.Format("2006-01-02 15:04:05.000"),
	}
}
