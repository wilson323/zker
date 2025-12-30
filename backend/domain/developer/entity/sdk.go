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

package entity

// SDKLanguage SDK语言
type SDKLanguage string

const (
	SDKLanguagePython     SDKLanguage = "python"     // Python
	SDKLanguageJavaScript SDKLanguage = "javascript" // JavaScript/TypeScript
	SDKLanguageGo         SDKLanguage = "go"         // Go
	SDKLanguageJava       SDKLanguage = "java"       // Java
)

// SDKStatus SDK状态
type SDKStatus string

const (
	SDKStatusDraft     SDKStatus = "draft"     // 草稿
	SDKStatusPublished SDKStatus = "published" // 已发布
	SDKStatusArchived  SDKStatus = "archived"  // 已归档
)

// SDK SDK实体
type SDK struct {
	SDKID        string      `json:"sdk_id" gorm:"primaryKey;type:varchar(36);comment:SDK ID"`
	TenantID     string      `json:"tenant_id" gorm:"type:varchar(36);not null;index:idx_tenant_id;comment:租户ID"`
	ProjectID    string      `json:"project_id" gorm:"type:varchar(36);not null;index:idx_project_id;comment:项目ID"`
	SDKName      string      `json:"sdk_name" gorm:"type:varchar(100);not null;comment:SDK名称"`
	Language     SDKLanguage `json:"language" gorm:"type:enum('python','javascript','go','java');not null;index:idx_language;comment:编程语言"`
	Version      string      `json:"version" gorm:"type:varchar(20);not null;comment:版本号"`
	Description  string      `json:"description" gorm:"type:text;comment:描述"`
	Code         string      `json:"code,omitempty" gorm:"type:longtext;comment:生成的代码"`
	PackageURL   string      `json:"package_url,omitempty" gorm:"type:varchar(500);comment:包下载地址"`
	Readme       string      `json:"readme,omitempty" gorm:"type:text;comment:使用文档"`
	Status       SDKStatus   `json:"status" gorm:"type:enum('draft','published','archived');default:'draft';not null;index:idx_status;comment:状态"`
	DownloadCount int64      `json:"download_count" gorm:"default:0;comment:下载次数"`
	CreatedAt    int64       `json:"created_at" gorm:"not null;default:0;comment:创建时间"`
	UpdatedAt    int64       `json:"updated_at" gorm:"not null;default:0;comment:更新时间"`
	DeletedAt    *int64      `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`

	// 关联
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID;references:ProjectID"`
}

// TableName 指定表名
func (SDK) TableName() string {
	return "developer_sdks"
}

// IsPublished 是否已发布
func (s *SDK) IsPublished() bool {
	return s.Status == SDKStatusPublished
}

// GetVersionedName 获取带版本号的名称
func (s *SDK) GetVersionedName() string {
	return s.SDKName + "-v" + s.Version
}

// SDKPackage SDK包信息
type SDKPackage struct {
	SDKID       string     `json:"sdk_id"`
	FileName    string     `json:"file_name"`
	FileSize    int64      `json:"file_size"`
	FileMD5     string     `json:"file_md5"`
	DownloadURL string     `json:"download_url"`
	ExpiresAt   *int64     `json:"expires_at,omitempty"`
}

// SDKTemplate SDK模板配置
type SDKTemplate struct {
	Language     SDKLanguage `json:"language"`
	TemplateName string      `json:"template_name"`
	Content      string      `json:"content"`
	Dependencies []string    `json:"dependencies,omitempty"`
	Examples     []string    `json:"examples,omitempty"`
}
