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
	model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"
)

var (
	ErrDuplicateUsername = dao.ErrDuplicateUsername
)

type UserRepository interface {
	Register(ctx context.Context, u domain.Register) error

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

func (repo *userRepository) ExistsByUserName(ctx context.Context, email string) (bool, error) {
	return repo.dao.ExistsByUserName(ctx, email)
}

func (repo *userRepository) toEntity(d domain.Register) model.User {
	return model.User{
		Username: d.Username,
		Password: d.Password,
	}
}
