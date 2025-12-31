/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"github.com/coze-dev/coze-studio/backend/infra/tracing"
	pkgerrorx "github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	berrno "github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

const (
	// BucketNameTemplate Bucket命名模板：tenant-{tenant_id}
	BucketNameTemplate = "tenant-%s"

	// DefaultUploadExpiry 默认上传URL过期时间
	DefaultUploadExpiry = 24 * time.Hour

	// DefaultDownloadExpiry 默认下载URL过期时间
	DefaultDownloadExpiry = 1 * time.Hour
)

// TenantFileInfo 租户文件信息（重命名避免与storage.go中的FileInfo冲突）
type TenantFileInfo struct {
	Bucket      string            `json:"bucket"`
	Key         string            `json:"key"`
	Size        int64             `json:"size"`
	ETag        string            `json:"etag"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
}

// UploadResult 上传结果
type UploadResult struct {
	FileID   string `json:"file_id"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
	ETag     string `json:"etag"`
	Location string `json:"location"`
}

// TenantIsolatedStorage 租户隔离存储
//
// **核心功能**：
// 1. 自动为所有文件添加租户Bucket前缀
// 2. Bucket命名：tenant-{tenant_id}
// 3. 防止跨租户文件访问
// 4. 提供预签名URL
//
// **Bucket命名示例**：
// - tenant-tenant-a1b2c3d4-e5f6-7890-abcd-ef1234567890
// - tenant-tenant-b2c3d4e5-f6g7-8901-bcde-f23456789012
//
// **文件路径格式**：
// {bucket}/{category}/{year}/{month}/{day}/{file_id}
//
// **使用示例**：
//
//	storage := NewTenantIsolatedStorage(minioClient)
//	result, err := storage.Upload(ctx, tenantID, file)
type TenantIsolatedStorage struct {
	client *minio.Client

	// EnableStats 是否启用统计
	EnableStats bool

	// Stats 统计信息
	Stats *StorageStats
}

// StorageStats 存储统计
type StorageStats struct {
	TotalUploads   int64
	TotalDownloads int64
	TotalDeletes   int64
	TotalBytes     int64
}

// NewTenantIsolatedStorage 创建租户隔离存储
func NewTenantIsolatedStorage(endpoint, accessKey, secretKey string, useSSL bool) (*TenantIsolatedStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, pkgerrorx.WrapWithZap(err, berrno.ErrInitStorageFailed,
			zap.String("endpoint", endpoint),
		)
	}

	return &TenantIsolatedStorage{
		client:      client,
		EnableStats: true,
		Stats:       &StorageStats{},
	}, nil
}

// buildBucketName 构建Bucket名称：tenant-{tenant_id}
func (s *TenantIsolatedStorage) buildBucketName(tenantID string) string {
	return fmt.Sprintf(BucketNameTemplate, tenantID)
}

// Upload 上传文件
//
// **参数**：
// - ctx: 上下文
// - tenantID: 租户ID
// - category: 文件分类（如：avatar/document/knowledge等）
// - fileID: 文件ID
// - reader: 文件内容
// - size: 文件大小
// - contentType: MIME类型
// - metadata: 元数据
//
// **返回**：上传结果
//
// **示例**：
//
//	result, err := storage.Upload(ctx, tenantID, "avatar", "file-123", reader, size, "image/png", metadata)
func (s *TenantIsolatedStorage) Upload(
	ctx context.Context,
	tenantID, category, fileID string,
	reader io.Reader,
	size int64,
	contentType string,
	metadata map[string]string,
) (*UploadResult, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("category", category),
		attribute.String("file_id", fileID),
		attribute.Int64("size", size),
	)

	// 确保Bucket存在
	bucketName := s.buildBucketName(tenantID)
	if err := s.ensureBucket(ctx, bucketName); err != nil {
		return nil, err
	}

	// 构建文件路径：{category}/{year}/{month}/{day}/{fileID}
	now := time.Now()
	key := fmt.Sprintf("%s/%d/%02d/%02d/%s",
		category,
		now.Year(),
		now.Month(),
		now.Day(),
		fileID,
	)

	// 添加租户ID到元数据
	if metadata == nil {
		metadata = make(map[string]string)
	}
	metadata["tenant_id"] = tenantID
	metadata["category"] = category
	metadata["file_id"] = fileID

	// 上传文件
	uploadInfo, err := s.client.PutObject(ctx, bucketName, key, reader, size, minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: metadata,
	})

	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] upload failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return nil, pkgerrorx.WrapWithZap(err, berrno.ErrUploadFileFailed,
			zap.String("tenant_id", tenantID),
			zap.String("category", category),
			zap.String("file_id", fileID),
		)
	}

	// 记录统计
	if s.EnableStats {
		s.Stats.TotalUploads++
		s.Stats.TotalBytes += size
	}

	result := &UploadResult{
		FileID:   fileID,
		Bucket:   bucketName,
		Key:      key,
		Size:     uploadInfo.Size,
		ETag:     uploadInfo.ETag,
		Location: fmt.Sprintf("%s/%s", bucketName, key),
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] upload success: bucket=%s, key=%s, size=%d",
		bucketName, key, uploadInfo.Size)

	return result, nil
}

// Download 下载文件
//
// **示例**：
//
//	reader, err := storage.Download(ctx, tenantID, "avatar/file-123")
func (s *TenantIsolatedStorage) Download(
	ctx context.Context,
	tenantID, key string,
) (*minio.Object, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("key", key),
	)

	bucketName := s.buildBucketName(tenantID)

	// 获取对象
	object, err := s.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] download failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return nil, pkgerrorx.WrapWithZap(err, berrno.ErrDownloadFileFailed,
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// 记录统计
	if s.EnableStats {
		s.Stats.TotalDownloads++
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedStorage] download success: bucket=%s, key=%s", bucketName, key)
	return object, nil
}

// GetPresignedDownloadURL 获取预签名下载URL
//
// **参数**：
// - ctx: 上下文
// - tenantID: 租户ID
// - key: 文件key
// - expiry: 过期时间
//
// **返回**：预签名URL
//
// **示例**：
//
//	url, err := storage.GetPresignedDownloadURL(ctx, tenantID, "avatar/file-123", 1*time.Hour)
func (s *TenantIsolatedStorage) GetPresignedDownloadURL(
	ctx context.Context,
	tenantID, key string,
	expiry time.Duration,
) (string, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("key", key),
		attribute.String("expiry", expiry.String()),
	)

	bucketName := s.buildBucketName(tenantID)

	// 生成预签名URL
	url, err := s.client.PresignedGetObject(ctx, bucketName, key, expiry, nil)
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] generate download URL failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return "", pkgerrorx.WrapWithZap(err, berrno.ErrGenerateDownloadURLFailed,
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedStorage] download URL generated: bucket=%s, key=%s", bucketName, key)
	return url.String(), nil
}

// GetPresignedUploadURL 获取预签名上传URL
//
// **使用场景**：前端直接上传到对象存储
//
// **示例**：
//
//	url, err := storage.GetPresignedUploadURL(ctx, tenantID, "avatar", "file-123", 1*time.Hour)
func (s *TenantIsolatedStorage) GetPresignedUploadURL(
	ctx context.Context,
	tenantID, category, fileID string,
	expiry time.Duration,
) (string, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("category", category),
		attribute.String("file_id", fileID),
	)

	// 确保Bucket存在
	bucketName := s.buildBucketName(tenantID)
	if err := s.ensureBucket(ctx, bucketName); err != nil {
		return "", err
	}

	// 构建文件路径
	now := time.Now()
	key := fmt.Sprintf("%s/%d/%02d/%02d/%s",
		category,
		now.Year(),
		now.Month(),
		now.Day(),
		fileID,
	)

	// 生成预签名上传URL
	url, err := s.client.PresignedPutObject(ctx, bucketName, key, expiry)
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] generate upload URL failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return "", pkgerrorx.WrapWithZap(err, berrno.ErrGenerateUploadURLFailed,
			zap.String("tenant_id", tenantID),
			zap.String("category", category),
			zap.String("file_id", fileID),
		)
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedStorage] upload URL generated: bucket=%s, key=%s", bucketName, key)
	return url.String(), nil
}

// Delete 删除文件
//
// **示例**：
//
//	err := storage.Delete(ctx, tenantID, "avatar/file-123")
func (s *TenantIsolatedStorage) Delete(ctx context.Context, tenantID, key string) error {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("key", key),
	)

	bucketName := s.buildBucketName(tenantID)

	// 删除对象
	if err := s.client.RemoveObject(ctx, bucketName, key, minio.RemoveObjectOptions{}); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] delete failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return pkgerrorx.WrapWithZap(err, berrno.ErrDeleteFileFailed,
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	// 记录统计
	if s.EnableStats {
		s.Stats.TotalDeletes++
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] delete success: bucket=%s, key=%s", bucketName, key)
	return nil
}

// DeleteMultiple 批量删除文件
//
// **示例**：
//
//	keys := []string{"avatar/file-1", "avatar/file-2", "document/file-3"}
//	err := storage.DeleteMultiple(ctx, tenantID, keys)
func (s *TenantIsolatedStorage) DeleteMultiple(
	ctx context.Context,
	tenantID string,
	keys []string,
) error {
	if len(keys) == 0 {
		return nil
	}

	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.Int("key_count", len(keys)),
	)

	bucketName := s.buildBucketName(tenantID)

	// 构建删除对象列表
	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		for _, key := range keys {
			info, _ := s.client.StatObject(ctx, bucketName, key, minio.StatObjectOptions{})
			objectsCh <- info
		}
		close(objectsCh)
	}()

	// 批量删除
	errorCh := s.client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{})

	// 收集错误
	var errors []error
	for err := range errorCh {
		if err.Err != nil {
			logs.CtxWarnf(ctx, "[TenantIsolatedStorage] delete multiple partial failure: key=%s, error=%v",
				err.ObjectName, err.Err)
			errors = append(errors, err.Err)
		}
	}

	if len(errors) > 0 {
		return pkgerrorx.NewByErrorCode(berrno.ErrDeleteFilePartialFailed)
	}

	// 记录统计
	if s.EnableStats {
		s.Stats.TotalDeletes += int64(len(keys))
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] delete multiple success: bucket=%s, count=%d", bucketName, len(keys))
	return nil
}

// GetFileInfo 获取文件信息
//
// **示例**：
//
//	info, err := storage.GetFileInfo(ctx, tenantID, "avatar/file-123")
func (s *TenantIsolatedStorage) GetFileInfo(
	ctx context.Context,
	tenantID, key string,
) (*TenantFileInfo, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("key", key),
	)

	bucketName := s.buildBucketName(tenantID)

	// 获取对象信息
	stat, err := s.client.StatObject(ctx, bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] get file info failed: bucket=%s, key=%s, error=%v",
			bucketName, key, err)

		return nil, pkgerrorx.WrapWithZap(err, berrno.ErrGetFileInfoFailed,
			zap.String("tenant_id", tenantID),
			zap.String("key", key),
		)
	}

	info := &TenantFileInfo{
		Bucket:      bucketName,
		Key:         key,
		Size:        stat.Size,
		ETag:        stat.ETag,
		ContentType: stat.ContentType,
		Metadata:    stat.UserMetadata,
		CreatedAt:   stat.LastModified,
	}

	logs.CtxDebugf(ctx, "[TenantIsolatedStorage] get file info success: bucket=%s, key=%s, size=%d",
		bucketName, key, stat.Size)

	return info, nil
}

// ListFiles 列出文件
//
// **参数**：
// - ctx: 上下文
// - tenantID: 租户ID
// - category: 文件分类
// - prefix: 前缀过滤
// - maxCount: 最大返回数量
//
// **示例**：
//
//	files, err := storage.ListFiles(ctx, tenantID, "avatar", "", 100)
func (s *TenantIsolatedStorage) ListFiles(
	ctx context.Context,
	tenantID, category, prefix string,
	maxCount int,
) ([]TenantFileInfo, error) {
	tracing.AddSpanAttributes(ctx,
		tracing.AttrTenantID.String(tenantID),
		attribute.String("category", category),
		attribute.String("prefix", prefix),
		attribute.Int("max_count", maxCount),
	)

	bucketName := s.buildBucketName(tenantID)

	// 构建完整前缀：{category}/{prefix}
	fullPrefix := category
	if prefix != "" {
		fullPrefix = fmt.Sprintf("%s/%s", category, prefix)
	}

	// 列出对象
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objectCh := s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    fullPrefix,
		Recursive: true,
		MaxKeys:   maxCount,
	})

	files := make([]TenantFileInfo, 0, maxCount)
	for object := range objectCh {
		if object.Err != nil {
			logs.CtxErrorf(ctx, "[TenantIsolatedStorage] list files error: bucket=%s, error=%v",
				bucketName, object.Err)
			continue
		}

		files = append(files, TenantFileInfo{
			Bucket:      bucketName,
			Key:         object.Key,
			Size:        object.Size,
			ETag:        object.ETag,
			ContentType: "",
			Metadata:    make(map[string]string),
			CreatedAt:   object.LastModified,
		})
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] list files success: bucket=%s, count=%d", bucketName, len(files))
	return files, nil
}

// ClearTenantBucket 清空租户Bucket
//
// **使用场景**：租户删除、数据清理
//
// **示例**：
//
//	count, size, err := storage.ClearTenantBucket(ctx, tenantID)
func (s *TenantIsolatedStorage) ClearTenantBucket(
	ctx context.Context,
	tenantID string,
) (int64, int64, error) {
	logs.CtxInfof(ctx, "[TenantIsolatedStorage] clearing bucket for tenant: %s", tenantID)

	bucketName := s.buildBucketName(tenantID)

	// 检查Bucket是否存在
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] check bucket exists failed: bucket=%s, error=%v",
			bucketName, err)
		return 0, 0, pkgerrorx.WrapWithZap(err, berrno.ErrCheckBucketFailed,
			zap.String("tenant_id", tenantID),
		)
	}

	if !exists {
		logs.CtxInfof(ctx, "[TenantIsolatedStorage] bucket not exists: %s", bucketName)
		return 0, 0, nil
	}

	// 列出所有对象
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objectCh := s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	var count int64
	var size int64
	for object := range objectCh {
		if object.Err != nil {
			logs.CtxWarnf(ctx, "[TenantIsolatedStorage] list object error: key=%s, error=%v",
				object.Key, object.Err)
			continue
		}

		count++
		size += object.Size
	}

	// 删除所有对象
	if count > 0 {
		objectsCh := make(chan minio.ObjectInfo)
		go func() {
			for object := range objectCh {
				if object.Err == nil {
					objectsCh <- object
				}
			}
			close(objectsCh)
		}()

		errorCh := s.client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{})
		for err := range errorCh {
			if err.Err != nil {
				logs.CtxWarnf(ctx, "[TenantIsolatedStorage] delete object error: key=%s, error=%v",
					err.ObjectName, err.Err)
			}
		}
	}

	// 删除Bucket
	if err := s.client.RemoveBucket(ctx, bucketName); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] remove bucket failed: bucket=%s, error=%v",
			bucketName, err)
		return 0, 0, pkgerrorx.WrapWithZap(err, berrno.ErrDeleteBucketFailed,
			zap.String("tenant_id", tenantID),
		)
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] bucket cleared: bucket=%s, count=%d, size=%d",
		bucketName, count, size)

	return count, size, nil
}

// GetStats 获取存储统计
func (s *TenantIsolatedStorage) GetStats() StorageStats {
	return *s.Stats
}

// ResetStats 重置统计
func (s *TenantIsolatedStorage) ResetStats() {
	s.Stats = &StorageStats{}
}

// ensureBucket 确保Bucket存在
func (s *TenantIsolatedStorage) ensureBucket(ctx context.Context, bucketName string) error {
	// 检查Bucket是否存在
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] check bucket exists failed: bucket=%s, error=%v",
			bucketName, err)
		return pkgerrorx.WrapWithZap(err, berrno.ErrCheckBucketFailed,
			zap.String("bucket", bucketName),
		)
	}

	if exists {
		return nil
	}

	// 创建Bucket
	if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
		logs.CtxErrorf(ctx, "[TenantIsolatedStorage] create bucket failed: bucket=%s, error=%v",
			bucketName, err)
		return pkgerrorx.WrapWithZap(err, berrno.ErrCreateBucketFailed,
			zap.String("bucket", bucketName),
		)
	}

	logs.CtxInfof(ctx, "[TenantIsolatedStorage] bucket created: %s", bucketName)
	return nil
}
