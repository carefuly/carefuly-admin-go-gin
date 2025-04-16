/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/28 11:42:12
 * Remark：
 */

package system

import (
	"context"
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	"github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/auth"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/repository/careful/system"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateUsername     = system.ErrDuplicateUsername
	ErrInvalidUserOrPassword = errors.New("用户账号或密码错误")
	ErrGenerateTokenError    = errors.New("生成Token异常")
)

type UserService interface {
	UsernameAndPasswordCreate(ctx context.Context, u auth.Register) error
	FindByUserName(ctx *gin.Context, rely config.RelyConfig, u auth.Login) (string, error)
}

type userService struct {
	repo         system.UserRepository
	userPassRepo system.UserPassWordRepository
}

func NewUserService(repo system.UserRepository, userPassRepo system.UserPassWordRepository) UserService {
	return &userService{
		repo:         repo,
		userPassRepo: userPassRepo,
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UId       string
	Username  string
	Name      string
	Email     string
	Mobile    string
	UserAgent string
}

// UsernameAndPasswordCreate 用户账号和密码注册
func (svc *userService) UsernameAndPasswordCreate(ctx context.Context, u auth.Register) error {
	exists, err := svc.repo.ExistsByUserName(ctx, u.Username)
	if err != nil {
		return err
	}
	if exists {
		return system.ErrDuplicateUsername
	}

	text := u.Password

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	err = svc.repo.Create(ctx, u)
	if svc.IsDuplicateEntryError(err) {
		return system.ErrDuplicateUsername
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
		return system.ErrDuplicateUsername
	}

	err = svc.userPassRepo.Create(ctx, auth.Register{
		Username: u.Username,
		Password: text,
	})
	if svc.IsDuplicateEntryError(err) {
		return system.ErrDuplicateUsername
	}

	return err
}

// FindByUserName 根据用户名查找用户
func (svc *userService) FindByUserName(ctx *gin.Context, rely config.RelyConfig, u auth.Login) (string, error) {
	user, err := svc.repo.FindByUserName(ctx, u.Username)
	if err != nil {
		return "", ErrInvalidUserOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return "", ErrInvalidUserOrPassword
	}

	return svc.setJWTToken(ctx, rely, user)
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

// setJWTToken 设置JWT Token
func (svc *userService) setJWTToken(ctx *gin.Context, rely config.RelyConfig, u domainSystem.User) (string, error) {
	uc := UserClaims{
		UId:       u.Id,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Mobile:    u.Mobile,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 一小时过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60 * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte(rely.Token.ApiKeyAuth))
	if err != nil {
		return "", ErrGenerateTokenError
	}

	return tokenStr, err
}
