/**
 * Description：
 * FileName：json.go
 * Author：CJiaの用心
 * Create：2025/4/16 10:28:23
 * Remark：
 */

package jsonformat

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

func FormatJsonPrint(data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		zap.L().Error("序列化失败", zap.Error(err))
	}
	zap.L().Debug(string(jsonData))
	fmt.Println()
}
