/**
 * Description：
 * FileName：userPassword.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:39:15
 * Remark：
 */

package system

import (
	"context"
	"github.com/carefuly/carefuly-admin-go-gin/internal/dao/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/auth"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
)

type UserPassWordRepository interface {
	Create(ctx context.Context, u auth.Register) error

	ExistsByUserName(ctx context.Context, username string) (bool, error)
}

type userPassWordRepository struct {
	dao system.UserPasswordDAO
}

func NewUserPassWordRepository(dao system.UserPasswordDAO) UserPassWordRepository {
	return &userPassWordRepository{
		dao: dao,
	}
}

func (repo *userPassWordRepository) Create(ctx context.Context, u auth.Register) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *userPassWordRepository) ExistsByUserName(ctx context.Context, username string) (bool, error) {
	return repo.dao.ExistsByUserName(ctx, username)
}

func (repo *userPassWordRepository) toEntity(d auth.Register) modelSystem.UserPassword {
	return modelSystem.UserPassword{
		Username: d.Username,
		Password: d.Password,
	}
}
