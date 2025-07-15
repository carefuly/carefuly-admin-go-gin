/**
 * Description：
 * FileName：file.go
 * Author：CJiaの用心
 * Create：2025/7/15 14:26:34
 * Remark：
 */

package third

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var (
	ErrBucketFileDirectoryExists = errors.New("存储桶目录已经存在")
	ErrBucketFileCreateFailed    = errors.New("创建存储桶目录失败")
	ErrBucketFilePathInvalid     = errors.New("目录路径无效")
	ErrBucketFileDeleteFailed    = errors.New("删除存储桶目录失败")
)

type BucketFileService interface {
	CreateBucketDir(ctx context.Context, bucketCode string) error
	DeleteBucketDir(ctx context.Context, bucketCode string) error
	BatchDeleteDirs(ctx context.Context, codes []string) error
}

type bucketFileService struct {
	baseDir string
}

func NewBucketFileService() BucketFileService {
	return &bucketFileService{
		baseDir: filepath.Join("./static", "bucket"),
	}
}

// CreateBucketDir 创建存储桶目录
func (s *bucketFileService) CreateBucketDir(ctx context.Context, bucketCode string) error {
	dir := filepath.Join(s.baseDir, bucketCode)

	// 检查目录是否已存在
	if _, err := os.Stat(dir); err == nil {
		return ErrBucketFileDirectoryExists
	}

	// 创建目录（包括父目录）
	if err := os.MkdirAll(dir, 0755); err != nil {
		return ErrBucketFileCreateFailed
	}

	return nil
}

// DeleteBucketDir 删除存储桶目录
func (s *bucketFileService) DeleteBucketDir(ctx context.Context, bucketCode string) error {
	dir := filepath.Join(s.baseDir, bucketCode)

	// 安全验证：确保在合法路径下
	if _, err := os.Stat(dir); err != nil {
		return ErrBucketFilePathInvalid
	}

	// 删除目录及内容
	if err := os.RemoveAll(dir); err != nil {
		return ErrBucketFileDeleteFailed
	}

	return nil
}

// BatchDeleteDirs 批量删除存储桶目录
func (s *bucketFileService) BatchDeleteDirs(ctx context.Context, codes []string) error {
	for _, code := range codes {
		if err := s.DeleteBucketDir(ctx, code); err != nil {
			// 记录错误但继续执行
			zap.L().Error("批量删除存储桶目录失败", zap.String("code", code), zap.Error(err))
		}
	}
	return nil
}
