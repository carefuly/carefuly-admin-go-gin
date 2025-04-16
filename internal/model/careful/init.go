/**
 * Description：
 * FileName：init.go
 * Author：CJiaの用心
 * Create：2025/4/16 16:05:58
 * Remark：
 */

package careful

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/logger"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"gorm.io/gorm"
)

func InitAutoMigrate(db *gorm.DB) {
	// initSystem(db)
	initTools(db)
	// initLogger(db)
}

func initSystem(db *gorm.DB) {
	system.NewUser().AutoMigrate(db)
	system.NewUserPassword().AutoMigrate(db)
}

func initTools(db *gorm.DB) {
	tools.NewDict().AutoMigrate(db)
	tools.NewDictType().AutoMigrate(db)
}

func initLogger(db *gorm.DB) {
	logger.NewOperateLogger().AutoMigrate(db)
}
