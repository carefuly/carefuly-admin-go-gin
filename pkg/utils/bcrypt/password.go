/**
 * Description：
 * FileName：password.go
 * Author：CJiaの用心
 * Create：2025/5/11 19:39:20
 * Remark：
 */

package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePasswords 比较密码是否匹配
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
