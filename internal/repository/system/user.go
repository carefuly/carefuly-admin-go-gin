/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/26 13:13:21
 * Remark：
 */

package repository

import (
	"context"
	dao "github.com/carefuly/carefuly-admin-go-gin/internal/dao/system"
	domain "github.com/carefuly/carefuly-admin-go-gin/internal/domain/auth"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/system"
	model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"
)

var (
	ErrDuplicateUsername = dao.ErrDuplicateUsername
	ErrUserNotFound      = dao.ErrUserNotFound
)

type UserRepository interface {
	Register(ctx context.Context, u domain.Register) error

	FindByUserName(ctx context.Context, username string) (domainSystem.User, error)

	ExistsByUserName(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

func (repo *userRepository) Register(ctx context.Context, u domain.Register) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *userRepository) FindByUserName(ctx context.Context, username string) (domainSystem.User, error) {
	user, err := repo.dao.FindByUserName(ctx, username)
	if err != nil {
		return domainSystem.User{}, err
	}
	return repo.toDomain(user), err
}

func (repo *userRepository) ExistsByUserName(ctx context.Context, email string) (bool, error) {
	return repo.dao.ExistsByUserName(ctx, email)
}

func (repo *userRepository) toEntity(d domain.Register) model.User {
	return model.User{
		Username: d.Username,
		Password: d.Password,
	}
}

func (repo *userRepository) toDomain(u *model.User) domainSystem.User {
	return domainSystem.User{
		User:       *u,
		CreateTime: u.CoreModels.CreateTime.Format("2006-01-02 15:04:05.000"),
		UpdateTime: u.CoreModels.UpdateTime.Format("2006-01-02 15:04:05.000"),
	}
}
