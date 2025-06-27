/**
 * Description：
 * FileName：xlsx.go
 * Author：CJiaの用心
 * Create：2025/5/22 16:20:51
 * Remark：
 */

package xlsx

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

type Xlsx struct {
	FilePath string
	file     *excelize.File // 内部文件引用
}

// NewXlsxFile 创建新的Xlsx实例
func NewXlsxFile(filePath string) *Xlsx {
	return &Xlsx{
		FilePath: filePath,
	}
}

// ReadFirstSheet 读取第一个Sheet
func (x *Xlsx) ReadFirstSheet() ([]map[string]string, error) {
	if err := x.openFile(); err != nil {
		return nil, err
	}

	sheets := x.file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("文件没有sheet")
	}
	return x.processSheet(sheets[0])
}

// ReadSheetByName 读取指定名称的Sheet
func (x *Xlsx) ReadSheetByName(sheetName string) ([]map[string]string, error) {
	if err := x.openFile(); err != nil {
		return nil, err
	}

	sheets := x.file.GetSheetList()
	exists := false
	for _, s := range sheets {
		if s == sheetName {
			exists = true
			break
		}
	}

	if !exists {
		return nil, fmt.Errorf("sheet[%s]不存在", sheetName)
	}

	return x.processSheet(sheetName)
}

// ReadAllSheets 读取所有Sheet
func (x *Xlsx) ReadAllSheets() (map[string][]map[string]string, error) {
	if err := x.openFile(); err != nil {
		return nil, err
	}

	result := make(map[string][]map[string]string)
	sheets := x.file.GetSheetList()

	for _, sheet := range sheets {
		data, err := x.processSheet(sheet)
		if err != nil {
			return nil, err
		}
		result[sheet] = data
	}

	return result, nil
}

// openFile 打开Excel文件（内部方法）
func (x *Xlsx) openFile() error {
	if x.file != nil {
		return nil // 文件已打开
	}

	f, err := excelize.OpenFile(x.FilePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	x.file = f
	return nil
}

// Close 关闭文件（使用完后必须调用）
func (x *Xlsx) Close() error {
	if x.file == nil {
		return nil
	}
	if err := x.file.Close(); err != nil {
		return fmt.Errorf("关闭文件失败: %w", err)
	}
	x.file = nil
	return nil
}

// processSheet 处理单个Sheet数据（内部方法）
func (x *Xlsx) processSheet(sheetName string) ([]map[string]string, error) {
	if err := x.openFile(); err != nil {
		return nil, err
	}

	rows, err := x.file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("获取sheet[%s]失败: %w", sheetName, err)
	}

	if len(rows) == 0 {
		return []map[string]string{}, nil
	}

	// 处理表头
	headers := rows[0]
	headerMap := make(map[string]int)
	uniqueHeaders := make([]string, len(headers))

	for i, h := range headers {
		// 去除首尾空白并处理重复表头
		trimmed := strings.TrimSpace(h)
		if trimmed == "" {
			trimmed = fmt.Sprintf("column_%d", i+1) // 空表头处理
		}

		count := headerMap[trimmed] + 1
		headerMap[trimmed] = count

		uniqueHeader := trimmed
		if count > 1 {
			uniqueHeader = fmt.Sprintf("%s_%d", trimmed, count)
		}
		uniqueHeaders[i] = uniqueHeader
	}

	// 处理数据行
	result := make([]map[string]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		rowMap := make(map[string]string, len(uniqueHeaders))

		for colIdx := 0; colIdx < len(uniqueHeaders); colIdx++ {
			var value string
			if colIdx < len(row) {
				value = strings.TrimSpace(row[colIdx])
			}
			rowMap[uniqueHeaders[colIdx]] = value
		}

		result = append(result, rowMap)
	}

	return result, nil
}
