package entity

import (
	"time"
)

// Archive 归档记录实体
type Archive struct {
	ArchiveID      string    `gorm:"column:archive_id;primaryKey;size:36" json:"archive_id"`
	TableName      string    `gorm:"column:table_name;size:64;not null" json:"table_name"`
	TenantID       string    `gorm:"column:tenant_id;size:36;not null;index:idx_tenant_id" json:"tenant_id"`
	DateStart      int64     `gorm:"column:date_start;not null;index:idx_date_range" json:"date_start"`
	DateEnd        int64     `gorm:"column:date_end;not null;index:idx_date_range" json:"date_end"`
	RecordCount    int64     `gorm:"column:record_count;not null;default:0" json:"record_count"`
	FilePath       string    `gorm:"column:file_path;size:512;not null" json:"file_path"`
	FileSize       int64     `gorm:"column:file_size;not null;default:0" json:"file_size"`
	Checksum       string    `gorm:"column:checksum;size:64;not null" json:"checksum"`
	Compression    string    `gorm:"column:compression;size:16;not null;default:'zstd'" json:"compression"`
	StorageClass   string    `gorm:"column:storage_class;size:32;not null;default:'glacier'" json:"storage_class"`
	Status         string    `gorm:"column:status;size:16;not null;default:'pending';index:idx_status" json:"status"`
	ErrorMessage   string    `gorm:"column:error_message;size:512" json:"error_message,omitempty"`
	CreatedAt      int64     `gorm:"column:created_at;not null;index:idx_created_at" json:"created_at"`
	UpdatedAt      int64     `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName 指定表名
func (Archive) TableName() string {
	return "archives"
}

// ArchiveTask 归档任务实体
type ArchiveTask struct {
	TaskID        string    `gorm:"column:task_id;primaryKey;size:36" json:"task_id"`
	TableName     string    `gorm:"column:table_name;size:64;not null" json:"table_name"`
	ArchiveType   string    `gorm:"column:archive_type;size:16;not null" json:"archive_type"` // hot, warm, cold
	PartitionKey  string    `gorm:"column:partition_key;size:64" json:"partition_key,omitempty"`
	Status        string    `gorm:"column:status;size:16;not null;default:'pending';index:idx_status" json:"status"`
	TotalRecords  int64     `gorm:"column:total_records;default:0" json:"total_records"`
	ProcessedRecords int64   `gorm:"column:processed_records;default:0" json:"processed_records"`
	FailedRecords int64     `gorm:"column:failed_records;default:0" json:"failed_records"`
	StartedAt     *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	CompletedAt   *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	ErrorMessage  string    `gorm:"column:error_message;size:512" json:"error_message,omitempty"`
	CreatedAt     int64     `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt     int64     `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName 指定表名
func (ArchiveTask) TableName() string {
	return "archive_tasks"
}

// ArchiveConfig 归档配置实体
type ArchiveConfig struct {
	ConfigID      string `gorm:"column:config_id;primaryKey;size:36" json:"config_id"`
	TableName     string `gorm:"column:table_name;size:64;not null;uniqueIndex:uk_table" json:"table_name"`
	ArchiveRule   string `gorm:"column:archive_rule;size:16;not null" json:"archive_rule"` // monthly, quarterly, yearly
	RetentionDays int    `gorm:"column:retention_days;not null;default:180" json:"retention_days"`
	Compression   string `gorm:"column:compression;size:16;not null;default:'zstd'" json:"compression"`
	StorageClass  string `gorm:"column:storage_class;size:32;not null;default:'standard_ia'" json:"storage_class"`
	Enabled       bool   `gorm:"column:enabled;not null;default:true" json:"enabled"`
	Priority      int    `gorm:"column:priority;not null;default:0" json:"priority"`
	CreatedAt     int64  `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName 指定表名
func (ArchiveConfig) TableName() string {
	return "archive_configs"
}

// ArchiveQuery 归档查询请求
type ArchiveQuery struct {
	TenantID    string `json:"tenant_id"`
	TableName   string `json:"table_name"`
	StartDate   int64  `json:"start_date"`
	EndDate     int64  `json:"end_date"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
}

// ArchiveQueryResult 归档查询结果
type ArchiveQueryResult struct {
	TotalRecords int64            `json:"total_records"`
	Records      []ArchiveRecord  `json:"records"`
	QueryTime    int64            `json:"query_time_ms"`
	ArchivesUsed []string         `json:"archives_used"`
}

// ArchiveRecord 归档数据记录
type ArchiveRecord struct {
	TenantID    string                 `json:"tenant_id"`
	Data        map[string]interface{} `json:"data"`
	CreatedAt   int64                  `json:"created_at"`
}

// ArchiveStats 归档统计信息
type ArchiveStats struct {
	TotalArchives    int64   `json:"total_archives"`
	TotalRecords     int64   `json:"total_records"`
	TotalSize        int64   `json:"total_size_gb"`
	CompressionRatio float64 `json:"compression_ratio"`
	StorageCost      float64 `json:"storage_cost_monthly"`
}

// ArchiveProgress 归档进度
type ArchiveProgress struct {
	TaskID           string  `json:"task_id"`
	Status           string  `json:"status"`
	TotalRecords     int64   `json:"total_records"`
	ProcessedRecords int64   `json:"processed_records"`
	Progress         float64 `json:"progress_percent"`
	ETA              int64   `json:"eta_seconds"`
}
