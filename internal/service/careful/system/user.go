/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/5/12 17:17:29
 * Remark：
 */

package system

import (
	"context"
	"errors"
	"fmt"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	repositorySystem "github.com/carefuly/carefuly-admin-go-gin/internal/repository/repository/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/bcrypt"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotFound          = repositorySystem.ErrUserNotFound
	ErrUsernameDuplicate     = repositorySystem.ErrUsernameDuplicate
	ErrUserDuplicate         = repositorySystem.ErrUserDuplicate
	ErrUserInvalidCredential = errors.New("用户名或密码错误")
	ErrUserTypeDoesNotMatch  = errors.New("用户类型不匹配")
)

type UserService interface {
	Register(ctx context.Context, user domainSystem.User) error
	Login(ctx context.Context, username, password string) (domainSystem.User, error)
	LoginWithType(ctx context.Context, username, password string, userType int) (domainSystem.User, error)
	ChangePassword(ctx context.Context, userId string, oldPassword, newPassword string) error

	Create(ctx context.Context, domain domainSystem.User) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, domain domainSystem.User) error

	GetById(ctx context.Context, id string) (domainSystem.User, error)
	GetByUsername(ctx context.Context, username string) (domainSystem.User, error)
	GetListPage(ctx context.Context, filter domainSystem.UserFilter) ([]domainSystem.User, int64, error)
	GetListAll(ctx context.Context, filter domainSystem.UserFilter) ([]domainSystem.User, error)
}

type userService struct {
	repo repositorySystem.UserRepository
}

func NewUserService(repo repositorySystem.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// Register 注册
func (svc *userService) Register(ctx context.Context, user domainSystem.User) error {
	// 检查用户名是否已存在
	exists, err := svc.repo.CheckExistByUsername(ctx, user.Username, "")
	if err != nil {
		return fmt.Errorf("检查用户名是否存在失败: %w", err)
	}
	if exists {
		return repositorySystem.ErrUsernameDuplicate
	}

	// 加密密码
	hashedPassword, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 保存明文密码（注意：实际生产环境不应存储明文密码）
	user.PasswordStr = user.Password
	user.Password = hashedPassword

	// 创建用户
	if err := svc.repo.Create(ctx, user); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrUsernameDuplicate
		}
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// Login 登录
func (svc *userService) Login(ctx context.Context, username, password string) (domainSystem.User, error) {
	// 根据用户名获取用户
	user, err := svc.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return domainSystem.User{}, ErrUserInvalidCredential
		}
		return domainSystem.User{}, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证密码
	if err := bcrypt.ComparePasswords(user.Password, password); !err {
		return domainSystem.User{}, ErrUserInvalidCredential
	}

	return user, nil
}

// LoginWithType 按用户类型登录
func (svc *userService) LoginWithType(ctx context.Context, username, password string, userType int) (domainSystem.User, error) {
	// 根据用户名获取用户
	user, err := svc.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return domainSystem.User{}, ErrUserInvalidCredential
		}
		return domainSystem.User{}, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证用户类型
	if user.UserType != userType {
		return domainSystem.User{}, ErrUserTypeDoesNotMatch
	}

	// 验证密码
	if err := bcrypt.ComparePasswords(user.Password, password); !err {
		return domainSystem.User{}, ErrUserInvalidCredential
	}

	return user, nil
}

// ChangePassword 修改密码
func (svc *userService) ChangePassword(ctx context.Context, userId string, oldPassword, newPassword string) error {
	// 获取用户信息
	main, err := svc.repo.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrUserNotFound) {
			return repositorySystem.ErrUserNotFound
		}
		return err
	}

	// 验证旧密码
	if !bcrypt.ComparePasswords(main.Password, oldPassword) {
		return ErrUserInvalidCredential
	}

	// 加密新密码
	hashedPassword, err := bcrypt.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := svc.repo.UpdatePassword(ctx, userId, hashedPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// Create 创建
func (svc *userService) Create(ctx context.Context, domain domainSystem.User) error {
	// 检查用户名是否存在
	exists, err := svc.repo.CheckExistByUsername(ctx, domain.Username, "")
	if err != nil {
		return fmt.Errorf("检查用户名是否存在失败: %w", err)
	}
	if exists {
		return repositorySystem.ErrUsernameDuplicate
	}

	// 加密密码
	hashedPassword, err := bcrypt.HashPassword(domain.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 保存明文密码（注意：实际生产环境不应存储明文密码）
	domain.PasswordStr = domain.Password
	domain.Password = hashedPassword

	// 创建用户
	if err := svc.repo.Create(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrUsernameDuplicate
		}
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// Delete 删除
func (svc *userService) Delete(ctx context.Context, id string) error {
	rowsAffected, err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repositorySystem.ErrUserNotFound
	}
	return nil
}

// Update 更新
func (svc *userService) Update(ctx context.Context, domain domainSystem.User) error {
	// 检查用户是否存在
	exists, err := svc.repo.CheckExistByUsername(ctx, domain.Username, domain.Id)
	if err != nil {
		return fmt.Errorf("检查用户名是否存在失败: %w", err)
	}
	if exists {
		return repositorySystem.ErrUsernameDuplicate
	}

	if err := svc.repo.Update(ctx, domain); err != nil {
		if svc.IsDuplicateEntryError(err) {
			return repositorySystem.ErrUsernameDuplicate
		}
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

// GetById 获取详情
func (svc *userService) GetById(ctx context.Context, id string) (domainSystem.User, error) {
	main, err := svc.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repositorySystem.ErrUserNotFound) {
			return main, repositorySystem.ErrUserNotFound
		}
		return main, err
	}
	if main.Id == "" {
		return main, repositorySystem.ErrUserNotFound
	}
	return main, nil
}

// GetByUsername 根据用户名获取详情
func (svc *userService) GetByUsername(ctx context.Context, username string) (domainSystem.User, error) {
	// 调用仓库层获取用户信息
	return svc.repo.GetByUsername(ctx, username)
}

// GetListPage 分页查询列表
func (svc *userService) GetListPage(ctx context.Context, filter domainSystem.UserFilter) ([]domainSystem.User, int64, error) {
	return svc.repo.GetListPage(ctx, filter)
}

// GetListAll 查询所有列表
func (svc *userService) GetListAll(ctx context.Context, filter domainSystem.UserFilter) ([]domainSystem.User, error) {
	return svc.repo.GetListAll(ctx, filter)
}

// IsDuplicateEntryError 判断是否是唯一冲突错误
func (svc *userService) IsDuplicateEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL 错误码 1062 表示唯一冲突
		return mysqlErr.Number == 1062
	}
	return false
}
