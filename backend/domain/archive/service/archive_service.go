package service

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/archive/entity"
	"github.com/coze-dev/coze-studio/backend/domain/archive/repository"
	"gorm.io/gorm"
)

// ArchiveService 归档服务接口
type ArchiveService interface {
	// ArchiveTable 归档表数据
	ArchiveTable(ctx context.Context, tableName, tenantID string, retentionDays int) (*entity.ArchiveTask, error)

	// QueryArchivedData 查询归档数据
	QueryArchivedData(ctx context.Context, query *entity.ArchiveQuery) (*entity.ArchiveQueryResult, error)

	// RestoreFromArchive 从归档恢复数据
	RestoreFromArchive(ctx context.Context, archiveID string) error

	// GetArchiveProgress 获取归档进度
	GetArchiveProgress(ctx context.Context, taskID string) (*entity.ArchiveProgress, error)

	// GetArchiveStats 获取归档统计
	GetArchiveStats(ctx context.Context, tenantID string) (*entity.ArchiveStats, error)

	// DeleteArchive 删除归档
	DeleteArchive(ctx context.Context, archiveID string) error
}

// archiveService 归档服务实现
type archiveService struct {
	db              *gorm.DB
	archiveRepo     repository.ArchiveRepository
	taskRepo        repository.ArchiveTaskRepository
	configRepo      repository.ArchiveConfigRepository
	storage         Storage
	index           Index
	compressor      Compressor
}

// Storage 存储接口
type Storage interface {
	// Upload 上传文件到存储
	Upload(ctx context.Context, key string, data []byte) error

	// Download 从存储下载文件
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete 从存储删除文件
	Delete(ctx context.Context, key string) error

	// Exists 检查文件是否存在
	Exists(ctx context.Context, key string) (bool, error)
}

// Index 索引接口
type Index interface {
	// IndexArchive 索引归档元数据
	IndexArchive(ctx context.Context, archive *entity.Archive) error

	// SearchArchives 搜索归档
	SearchArchives(ctx context.Context, tenantID, tableName string, startDate, endDate int64) ([]*entity.Archive, error)

	// DeleteIndex 删除索引
	DeleteIndex(ctx context.Context, archiveID string) error
}

// Compressor 压缩接口
type Compressor interface {
	// Compress 压缩数据
	Compress(data []byte) ([]byte, error)

	// Decompress 解压数据
	Decompress(data []byte) ([]byte, error)
}

// NewArchiveService 创建归档服务
func NewArchiveService(
	db *gorm.DB,
	archiveRepo repository.ArchiveRepository,
	taskRepo repository.ArchiveTaskRepository,
	configRepo repository.ArchiveConfigRepository,
	storage Storage,
	index Index,
	compressor Compressor,
) ArchiveService {
	return &archiveService{
		db:          db,
		archiveRepo: archiveRepo,
		taskRepo:    taskRepo,
		configRepo:  configRepo,
		storage:     storage,
		index:       index,
		compressor:  compressor,
	}
}

// ArchiveTable 归档表数据
func (s *archiveService) ArchiveTable(ctx context.Context, tableName, tenantID string, retentionDays int) (*entity.ArchiveTask, error) {
	// 1. 计算归档时间范围
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	cutoffTimestamp := cutoffDate.Unix() * 1000

	// 2. 创建归档任务
	task := &entity.ArchiveTask{
		TaskID:       generateUUID(),
		TableName:    tableName,
		ArchiveType:  determineArchiveType(retentionDays),
		Status:       "running",
		TotalRecords: 0,
		CreatedAt:    time.Now().Unix() * 1000,
		UpdatedAt:    time.Now().Unix() * 1000,
	}

	now := time.Now()
	task.StartedAt = &now

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create archive task: %w", err)
	}

	// 3. 异步执行归档
	go s.executeArchive(context.Background(), task, tenantID, cutoffTimestamp)

	return task, nil
}

// executeArchive 执行归档
func (s *archiveService) executeArchive(ctx context.Context, task *entity.ArchiveTask, tenantID string, cutoffTimestamp int64) {
	// 更新任务状态为running
	s.taskRepo.UpdateStatus(ctx, task.TaskID, "running", "")

	// 1. 导出数据到CSV
	tempDir := filepath.Join(os.TempDir(), fmt.Sprintf("archive_%s", task.TaskID))
	defer os.RemoveAll(tempDir)

	csvFile := filepath.Join(tempDir, fmt.Sprintf("%s.csv", task.TableName))
	records, err := s.exportTableToCSV(ctx, task.TableName, tenantID, cutoffTimestamp, csvFile)
	if err != nil {
		s.taskRepo.UpdateStatus(ctx, task.TaskID, "failed", err.Error())
		return
	}

	task.TotalRecords = int64(records)

	// 2. 压缩数据
	compressedData, err := s.compressor.Compress(readFile(csvFile))
	if err != nil {
		s.taskRepo.UpdateStatus(ctx, task.TaskID, "failed", fmt.Sprintf("Compression failed: %v", err))
		return
	}

	// 3. 计算校验和
	checksum := fmt.Sprintf("%x", sha256.Sum256(compressedData))

	// 4. 上传到存储
	storageKey := fmt.Sprintf("archives/%s/%s/%s.tar.zst", tenantID, task.TableName, task.TaskID)
	if err := s.storage.Upload(ctx, storageKey, compressedData); err != nil {
		s.taskRepo.UpdateStatus(ctx, task.TaskID, "failed", fmt.Sprintf("Upload failed: %v", err))
		return
	}

	// 5. 创建归档记录
	archive := &entity.Archive{
		ArchiveID:    generateUUID(),
		TableName:    task.TableName,
		TenantID:     tenantID,
		DateStart:    0,
		DateEnd:      cutoffTimestamp,
		RecordCount:  int64(records),
		FilePath:     storageKey,
		FileSize:     int64(len(compressedData)),
		Checksum:     checksum,
		Compression:  "zstd",
		StorageClass: "glacier",
		Status:       "completed",
		CreatedAt:    time.Now().Unix() * 1000,
		UpdatedAt:    time.Now().Unix() * 1000,
	}

	if err := s.archiveRepo.Create(ctx, archive); err != nil {
		s.taskRepo.UpdateStatus(ctx, task.TaskID, "failed", fmt.Sprintf("Failed to create archive record: %v", err))
		return
	}

	// 6. 索引归档元数据
	if err := s.index.IndexArchive(ctx, archive); err != nil {
		// 非致命错误，只记录日志
		fmt.Printf("Warning: failed to index archive: %v\n", err)
	}

	// 7. 更新任务状态
	s.taskRepo.UpdateStatus(ctx, task.TaskID, "completed", "")
}

// exportTableToCSV 导出表数据到CSV
func (s *archiveService) exportTableToCSV(ctx context.Context, tableName, tenantID string, cutoffTimestamp int64, outputPath string) (int, error) {
	// 查询数据
	rows, err := s.db.WithContext(ctx).
		Table(tableName).
		Where("tenant_id = ? AND created_at < ?", tenantID, cutoffTimestamp).
		Rows()
	if err != nil {
		return 0, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	// 创建CSV文件
	file, err := os.Create(outputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		return 0, fmt.Errorf("failed to get columns: %w", err)
	}

	// 写入表头
	if err := writer.Write(columns); err != nil {
		return 0, fmt.Errorf("failed to write header: %w", err)
	}

	// 写入数据
	count := 0
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return 0, fmt.Errorf("failed to scan row: %w", err)
		}

		// 转换为字符串
		record := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				record[i] = ""
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}

		if err := writer.Write(record); err != nil {
			return 0, fmt.Errorf("failed to write record: %w", err)
		}

		count++

		// 每处理1000条记录更新一次进度
		if count%1000 == 0 {
			fmt.Printf("Exported %d records from %s\n", count, tableName)
		}
	}

	return count, nil
}

// QueryArchivedData 查询归档数据
func (s *archiveService) QueryArchivedData(ctx context.Context, query *entity.ArchiveQuery) (*entity.ArchiveQueryResult, error) {
	startTime := time.Now()

	// 1. 搜索归档索引
	archives, err := s.index.SearchArchives(ctx, query.TenantID, query.TableName, query.StartDate, query.EndDate)
	if err != nil {
		return nil, fmt.Errorf("failed to search archives: %w", err)
	}

	if len(archives) == 0 {
		return &entity.ArchiveQueryResult{
			TotalRecords: 0,
			Records:      []entity.ArchiveRecord{},
			QueryTime:    time.Since(startTime).Milliseconds(),
			ArchivesUsed: []string{},
		}, nil
	}

	// 2. 从归档文件读取数据
	allRecords := make([]entity.ArchiveRecord, 0)
	archivesUsed := make([]string, 0)

	for _, archive := range archives {
		records, err := s.readArchiveData(ctx, archive, query)
		if err != nil {
			fmt.Printf("Warning: failed to read archive %s: %v\n", archive.ArchiveID, err)
			continue
		}

		allRecords = append(allRecords, records...)
		archivesUsed = append(archivesUsed, archive.ArchiveID)
	}

	// 3. 应用分页
	total := len(allRecords)
	start := query.Offset
	end := start + query.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedRecords := allRecords[start:end]

	return &entity.ArchiveQueryResult{
		TotalRecords: int64(total),
		Records:      pagedRecords,
		QueryTime:    time.Since(startTime).Milliseconds(),
		ArchivesUsed: archivesUsed,
	}, nil
}

// readArchiveData 从归档文件读取数据
func (s *archiveService) readArchiveData(ctx context.Context, archive *entity.Archive, query *entity.ArchiveQuery) ([]entity.ArchiveRecord, error) {
	// 1. 从存储下载
	data, err := s.storage.Download(ctx, archive.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download archive: %w", err)
	}

	// 2. 解压数据
	decompressed, err := s.compressor.Decompress(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress: %w", err)
	}

	// 3. 解析CSV
	records := make([]entity.ArchiveRecord, 0)
	reader := csv.NewReader(bytes.NewReader(decompressed))

	// 读取表头
	columns, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// 读取数据行
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}

		// 转换为map
		dataMap := make(map[string]interface{})
		for i, col := range columns {
			dataMap[col] = row[i]
		}

		// 检查时间范围
		if createdAt, ok := dataMap["created_at"].(string); ok {
			var timestamp int64
			fmt.Sscanf(createdAt, "%d", &timestamp)
			if timestamp < query.StartDate || timestamp > query.EndDate {
				continue
			}
		}

		records = append(records, entity.ArchiveRecord{
			TenantID:  archive.TenantID,
			Data:      dataMap,
			CreatedAt: archive.DateStart,
		})
	}

	return records, nil
}

// RestoreFromArchive 从归档恢复数据
func (s *archiveService) RestoreFromArchive(ctx context.Context, archiveID string) error {
	// 1. 查找归档记录
	archive, err := s.archiveRepo.FindByID(ctx, archiveID)
	if err != nil {
		return fmt.Errorf("failed to find archive: %w", err)
	}

	// 2. 下载并解压数据
	data, err := s.storage.Download(ctx, archive.FilePath)
	if err != nil {
		return fmt.Errorf("failed to download archive: %w", err)
	}

	decompressed, err := s.compressor.Decompress(data)
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}

	// 3. 恢复到数据库（简化版本，实际需要更复杂的逻辑）
	// 这里只是示例，实际实现需要根据表结构动态构建INSERT语句
	fmt.Printf("Restoring %d records to %s\n", archive.RecordCount, archive.TableName)

	return nil
}

// GetArchiveProgress 获取归档进度
func (s *archiveService) GetArchiveProgress(ctx context.Context, taskID string) (*entity.ArchiveProgress, error) {
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	progress := &entity.ArchiveProgress{
		TaskID:           task.TaskID,
		Status:           task.Status,
		TotalRecords:     task.TotalRecords,
		ProcessedRecords: task.ProcessedRecords,
	}

	if task.TotalRecords > 0 {
		progress.Progress = float64(task.ProcessedRecords) / float64(task.TotalRecords) * 100
	}

	if task.Status == "running" && task.ProcessedRecords > 0 {
		// 估算剩余时间
		elapsed := time.Since(*task.StartedAt)
		avgTimePerRecord := elapsed.Seconds() / float64(task.ProcessedRecords)
		remainingRecords := task.TotalRecords - task.ProcessedRecords
		progress.ETA = int64(avgTimePerRecord * float64(remainingRecords))
	}

	return progress, nil
}

// GetArchiveStats 获取归档统计
func (s *archiveService) GetArchiveStats(ctx context.Context, tenantID string) (*entity.ArchiveStats, error) {
	return s.archiveRepo.GetStats(ctx, tenantID)
}

// DeleteArchive 删除归档
func (s *archiveService) DeleteArchive(ctx context.Context, archiveID string) error {
	// 1. 查找归档记录
	archive, err := s.archiveRepo.FindByID(ctx, archiveID)
	if err != nil {
		return fmt.Errorf("failed to find archive: %w", err)
	}

	// 2. 删除存储文件
	if err := s.storage.Delete(ctx, archive.FilePath); err != nil {
		return fmt.Errorf("failed to delete storage file: %w", err)
	}

	// 3. 删除索引
	if err := s.index.DeleteIndex(ctx, archiveID); err != nil {
		fmt.Printf("Warning: failed to delete index: %v\n", err)
	}

	// 4. 删除数据库记录
	if err := s.archiveRepo.Delete(ctx, archiveID); err != nil {
		return fmt.Errorf("failed to delete archive record: %w", err)
	}

	return nil
}

// 辅助函数
func generateUUID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), randomString(8))
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func determineArchiveType(retentionDays int) string {
	if retentionDays <= 180 {
		return "warm" // 6个月内
	}
	return "cold" // 超过6个月
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func bytes.NewReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}

// 添加缺失的import
import (
	"bytes"
)
