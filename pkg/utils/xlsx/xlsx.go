/**
 * Description：
 * FileName：xlsx.go
 * Author：CJiaの用心
 * Create：2025/4/16 00:01:01
 * Remark：
 */

package xlsx

import (
	"errors"
	"github.com/tealeg/xlsx"
	"strings"
)

var (
	ErrOpenFile      = errors.New("打开文件失败")
	ErrSheetNotFound = errors.New("指定的工作表不存在")
)

type Xlsx struct {
	FilePath string
}

func NewXlsxFile(filePath string) *Xlsx {
	return &Xlsx{
		FilePath: filePath,
	}
}

func (x *Xlsx) Read() ([]map[string]any, error) {
	// 解析Excel文件
	var result []map[string]any // 存储所有行的数据
	// 打开上传的Excel文件
	xlFile, err := xlsx.OpenFile(x.FilePath)
	if err != nil {
		// 打开文件失败
		return nil, ErrOpenFile
	}

	for _, sheet := range xlFile.Sheets {
		// 获取表头（第一行）
		header := sheet.Rows[0]

		var keys []string
		for _, cell := range header.Cells {
			keys = append(keys, cell.String())
		}

		// 遍历数据行（从第二行开始）
		for _, row := range sheet.Rows[1:] {
			// 跳过空行
			if row == nil || len(row.Cells) == 0 {
				continue
			}

			// 检查是否整行为空
			isEmpty := true
			for _, cell := range row.Cells {
				if strings.TrimSpace(cell.String()) != "" {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				continue
			}

			rowData := make(map[string]any)
			for colIndex, cell := range row.Cells {
				if colIndex < len(keys) {
					rowData[keys[colIndex]] = cell.String()
				}
			}
			result = append(result, rowData)
		}
	}

	return result, nil
}

func (x *Xlsx) ReadBySheet(sheetName string) ([]map[string]any, error) {
	// 打开上传的Excel文件
	xlFile, err := xlsx.OpenFile(x.FilePath)
	if err != nil {
		return nil, ErrOpenFile
	}

	// 查找指定的sheet
	var targetSheet *xlsx.Sheet
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == sheetName {
			targetSheet = sheet
			break
		}
	}

	if targetSheet == nil {
		return nil, ErrSheetNotFound
	}

	// 处理找到的sheet
	var result []map[string]any

	// 跳过空sheet
	if len(targetSheet.Rows) == 0 {
		return result, nil
	}

	// 获取表头（第一行）
	header := targetSheet.Rows[0]

	var keys []string
	for _, cell := range header.Cells {
		keys = append(keys, cell.String())
	}

	// 遍历数据行（从第二行开始）
	for _, row := range targetSheet.Rows[1:] {
		// 跳过空行
		if row == nil || len(row.Cells) == 0 {
			continue
		}

		// 检查是否整行为空
		isEmpty := true
		for _, cell := range row.Cells {
			if strings.TrimSpace(cell.String()) != "" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			continue
		}

		rowData := make(map[string]any)
		for colIndex, cell := range row.Cells {
			if colIndex < len(keys) {
				rowData[keys[colIndex]] = cell.String()
			}
		}
		result = append(result, rowData)
	}

	return result, nil
}
