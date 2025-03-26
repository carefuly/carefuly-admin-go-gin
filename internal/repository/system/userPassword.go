/**
 * Description：
 * FileName：userPassword.go
 * Author：CJiaの用心
 * Create：2025/3/26 13:13:42
 * Remark：
 */

package repository

import (
	"context"
	dao "github.com/carefuly/carefuly-admin-go-gin/internal/dao/system"
	domain "github.com/carefuly/carefuly-admin-go-gin/internal/domain/auth"
	model "github.com/carefuly/carefuly-admin-go-gin/internal/model/system"
)

type UserPassWordRepository interface {
	Create(ctx context.Context, u domain.Register) error

	ExistsByUserName(ctx context.Context, email string) (bool, error)
}

type userPassWordRepository struct {
	dao dao.UserPasswordDAO
}

func NewUserPassWordRepository(dao dao.UserPasswordDAO) UserPassWordRepository {
	return &userPassWordRepository{
		dao: dao,
	}
}

func (repo *userPassWordRepository) Create(ctx context.Context, u domain.Register) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *userPassWordRepository) ExistsByUserName(ctx context.Context, email string) (bool, error) {
	return repo.dao.ExistsByUserName(ctx, email)
}

func (repo *userPassWordRepository) toEntity(d domain.Register) model.UserPassword {
	return model.UserPassword{
		Username: d.Username,
		Password: d.Password,
	}
}
