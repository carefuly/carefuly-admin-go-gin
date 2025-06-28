/**
 * Description：
 * FileName：main.go
 * Author：CJiaの用心
 * Create：2025/6/28 22:54:58
 * Remark：
 */

package main

import (
	"fmt"
	"github.com/carefuly/carefuly-admin-go-gin/internal/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"time"
)

func main() {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// 导出路由
	r.GET("/api/export", func(c *gin.Context) {
		// 1. 创建Excel文件
		f := excelize.NewFile()
		defer f.Close()

		// 2. 创建工作表
		index, _ := f.NewSheet("Sheet1")

		// 3. 填充数据 (示例)
		f.SetCellValue("Sheet1", "A1", "ID")
		f.SetCellValue("Sheet1", "B1", "名称")
		f.SetCellValue("Sheet1", "C1", "创建时间")

		for i := 0; i < 100; i++ {
			row := i + 2
			f.SetCellValue("Sheet1", "A"+fmt.Sprint(row), i+1)
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(row), "项目 "+fmt.Sprint(i+1))
			f.SetCellValue("Sheet1", "C"+fmt.Sprint(row), time.Now().Format(time.RFC3339))
		}

		// 4. 设置活动工作表
		f.SetActiveSheet(index)

		// 5. 设置响应头
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=export.xlsx")
		c.Header("Pragma", "no-cache")
		c.Header("Cache-Control", "no-store")

		// 6. 流式写入响应
		if _, err := f.WriteTo(c.Writer); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "生成Excel失败"})
		}
	})

	r.Run(":8080")
}
