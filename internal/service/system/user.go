/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/3/26 22:24:21
 * Remark：
 */

package service

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domain "github.com/carefuly/carefuly-admin-go-gin/internal/domain/auth"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/system"
	repository "github.com/carefuly/carefuly-admin-go-gin/internal/repository/system"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateUsername
	ErrInvalidUserOrPassword = errors.New("邮箱或密码错误")
	ErrGenerateTokenError    = errors.New("生成Token异常")
)

type UserService interface {
	Login(ctx *gin.Context, rely config.RelyConfig, u domain.Login) (string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
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

// Login
// 按道理service不应该有业务登录的逻辑，这里只是为了方便看得懂
func (svc *userService) Login(ctx *gin.Context, rely config.RelyConfig, u domain.Login) (string, error) {
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60 * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte(rely.Token.ApiKeyAuth))
	if err != nil {
		return "", ErrGenerateTokenError
	}

	return tokenStr, err
}
