/**
 * Description：
 * FileName：xlsx.go
 * Author：CJiaの用心
 * Create：2025/5/22 16:20:51
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
	ErrInvalidRow    = errors.New("无效的行数据")
)

type Xlsx struct {
	FilePath string
}

func NewXlsxFile(filePath string) *Xlsx {
	return &Xlsx{
		FilePath: filePath,
	}
}

// Read 读取整个Excel文件的所有工作表数据
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

// ReadBySheet 读取指定工作表的数据
func (x *Xlsx) ReadBySheet(sheetName string) ([]map[string]string, error) {
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

	if len(targetSheet.Rows) == 0 {
		return []map[string]string{}, nil
	}

	// 获取表头
	header := targetSheet.Rows[0]
	keys := make([]string, len(header.Cells))
	for i, cell := range header.Cells {
		keys[i] = strings.TrimSpace(cell.String())
	}

	var result []map[string]string
	for _, row := range targetSheet.Rows[1:] {
		rowData, err := x.processRow(row, keys)
		if err != nil {
			continue // 跳过无效行
		}
		result = append(result, rowData)
	}

	return result, nil
}

// processRow 处理单行数据，返回map和错误
func (x *Xlsx) processRow(row *xlsx.Row, keys []string) (map[string]string, error) {
	// 检查空行
	if row == nil || len(row.Cells) == 0 {
		return nil, ErrInvalidRow
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
		return nil, ErrInvalidRow
	}

	// 确保所有字段都有值
	rowData := make(map[string]string)
	for i, key := range keys {
		var value string
		if i < len(row.Cells) {
			value = strings.TrimSpace(row.Cells[i].String())
		}
		rowData[key] = value
	}

	return rowData, nil
}
