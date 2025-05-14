/**
 * Description：
 * FileName：models_test.go.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:05:42
 * Remark：
 */

package models

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
	"testing"
)

func TestCoreModels_BeforeCreate(t *testing.T) {
	for i := 0; i < 20; i++ {
		fmt.Println(strings.ToUpper(uuid.NewV4().String()))
	}
}
