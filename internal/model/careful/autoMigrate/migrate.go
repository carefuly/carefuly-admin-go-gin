/**
 * Description：
 * FileName：migrate.go
 * Author：CJiaの用心
 * Create：2025/5/12 14:29:43
 * Remark：
 */

package autoMigrate

import (
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/logger"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	initSystem(db)
	initTools(db)
	initLogger(db)
}

func initSystem(db *gorm.DB) {
	system.NewUser().AutoMigrate(db)
	system.NewRole().AutoMigrate(db)
	system.NewMenu().AutoMigrate(db)
	system.NewMenuButton().AutoMigrate(db)
	system.NewMenuColumn().AutoMigrate(db)
	system.NewDept().AutoMigrate(db)
	system.NewPost().AutoMigrate(db)
}

func initTools(db *gorm.DB) {
	tools.NewDict().AutoMigrate(db)
	// tools.NewDictType().AutoMigrate(db)
}

func initLogger(db *gorm.DB) {
	logger.NewOperateLogger().AutoMigrate(db)
}
