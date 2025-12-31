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

package crypto

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TDEMiddleware 数据库透明加密中间件
type TDEMiddleware struct {
	db              *gorm.DB
	fieldEncryptor  *FieldEncryptor
	logger          *zap.Logger
	encryptionConfig *FieldEncryptionConfig
}

// TableEncryptionConfig 表加密配置
type TableEncryptionConfig struct {
	TableName      string     `json:"table_name"`
	EncryptedFields []string  `json:"encrypted_fields"` // 需要加密的字段列表
	TenantIDField  string     `json:"tenant_id_field"`  // 租户ID字段名
	Enabled        bool       `json:"enabled"`
}

// TDEConfig TDE配置
type TDEConfig struct {
	// 全局开关
	Enabled bool `json:"enabled"`

	// 表级别的加密配置
	TableConfigs []TableEncryptionConfig `json:"table_configs"`

	// 自动加密新记录
	AutoEncryptNewRecords bool `json:"auto_encrypt_new_records"`

	// 查询时自动解密
	AutoDecryptOnQuery bool `json:"auto_decrypt_on_query"`
}

// DefaultTDEConfig 默认TDE配置
func DefaultTDEConfig() *TDEConfig {
	return &TDEConfig{
		Enabled:               true,
		AutoEncryptNewRecords: true,
		AutoDecryptOnQuery:    true,
		TableConfigs: []TableEncryptionConfig{
			{
				TableName:      "users",
				EncryptedFields: []string{"password", "api_key", "secret_key"},
				TenantIDField:  "tenant_id",
				Enabled:        true,
			},
			{
				TableName:      "api_keys",
				EncryptedFields: []string{"access_token", "refresh_token"},
				TenantIDField:  "tenant_id",
				Enabled:        true,
			},
			{
				TableName:      "developer_platform",
				EncryptedFields: []string{"api_secret", "webhook_secret"},
				TenantIDField:  "tenant_id",
				Enabled:        true,
			},
		},
	}
}

// NewTDEMiddleware 创建TDE中间件
func NewTDEMiddleware(db *gorm.DB, fieldEncryptor *FieldEncryptor, logger *zap.Logger, config *TDEConfig) *TDEMiddleware {
	if config == nil {
		config = DefaultTDEConfig()
	}

	return &TDEMiddleware{
		db:              db,
		fieldEncryptor:  fieldEncryptor,
		logger:          logger,
		encryptionConfig: DefaultEncryptionConfig(),
	}
}

// RegisterHooks 注册数据库钩子
func (tde *TDEMiddleware) RegisterHooks(config *TDEConfig) error {
	if config == nil {
		config = DefaultTDEConfig()
	}

	if !config.Enabled {
		tde.logger.Info("TDE is disabled")
		return nil
	}

	// 为每个配置的表注册钩子
	for _, tableConfig := range config.TableConfigs {
		if !tableConfig.Enabled {
			continue
		}

		if err := tde.registerTableHooks(tableConfig); err != nil {
			return fmt.Errorf("failed to register hooks for table %s: %w", tableConfig.TableName, err)
		}
	}

	tde.logger.Info("Registered TDE hooks",
		zap.Int("tables", len(config.TableConfigs)),
	)

	return nil
}

// registerTableHooks 为表注册钩子
func (tde *TDEMiddleware) registerTableHooks(config TableEncryptionConfig) error {
	// 注册创建钩子(加密)
	if err := tde.db.Callback().Create().Before("gorm:create").Register("tde:before:create", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Statement.Schema.Table == config.TableName {
			if err := tde.encryptBeforeCreate(db, config); err != nil {
				tde.logger.Error("Failed to encrypt before create",
					zap.String("table", config.TableName),
					zap.Error(err),
				)
			}
		}
	}); err != nil {
		return fmt.Errorf("failed to register create hook: %w", err)
	}

	// 注册更新钩子(加密)
	if err := tde.db.Callback().Update().Before("gorm:update").Register("tde:before:update", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Statement.Schema.Table == config.TableName {
			if err := tde.encryptBeforeUpdate(db, config); err != nil {
				tde.logger.Error("Failed to encrypt before update",
					zap.String("table", config.TableName),
					zap.Error(err),
				)
			}
		}
	}); err != nil {
		return fmt.Errorf("failed to register update hook: %w", err)
	}

	// 注册查询钩子(解密)
	if err := tde.db.Callback().Query().After("gorm:query").Register("tde:after:query", func(db *gorm.DB) {
		if db.Statement.Schema != nil && db.Statement.Schema.Table == config.TableName {
			if err := tde.decryptAfterQuery(db, config); err != nil {
				tde.logger.Error("Failed to decrypt after query",
					zap.String("table", config.TableName),
					zap.Error(err),
				)
			}
		}
	}); err != nil {
		return fmt.Errorf("failed to register query hook: %w", err)
	}

	return nil
}

// encryptBeforeCreate 创建前加密
func (tde *TDEMiddleware) encryptBeforeCreate(db *gorm.DB, config TableEncryptionConfig) error {
	// 获取租户ID
	tenantID, err := tde.extractTenantID(db.Statement.Context, db, config.TenantIDField)
	if err != nil {
		return fmt.Errorf("failed to extract tenant_id: %w", err)
	}

	if tenantID == "" {
		// 没有租户ID,跳过加密
		return nil
	}

	// 加密敏感字段
	for _, field := range config.EncryptedFields {
		if fieldValue, ok := db.Statement.Schema.LookUpField(field); ok {
			if fieldValue != nil && db.Statement.Dest != nil {
				// TODO: 实现字段加密逻辑
				_ = tenantID
			}
		}
	}

	return nil
}

// encryptBeforeUpdate 更新前加密
func (tde *TDEMiddleware) encryptBeforeUpdate(db *gorm.DB, config TableEncryptionConfig) error {
	return tde.encryptBeforeCreate(db, config)
}

// decryptAfterQuery 查询后解密
func (tde *TDEMiddleware) decryptAfterQuery(db *gorm.DB, config TableEncryptionConfig) error {
	// 获取租户ID
	tenantID, err := tde.extractTenantID(db.Statement.Context, db, config.TenantIDField)
	if err != nil {
		return fmt.Errorf("failed to extract tenant_id: %w", err)
	}

	if tenantID == "" {
		// 没有租户ID,跳过解密
		return nil
	}

	// 解密敏感字段
	dest := db.Statement.Dest
	if dest == nil {
		return nil
	}

	// TODO: 实现字段解密逻辑
	_ = tenantID

	return nil
}

// extractTenantID 提取租户ID
func (tde *TDEMiddleware) extractTenantID(ctx context.Context, db *gorm.DB, tenantIDField string) (string, error) {
	// 1. 从上下文获取
	if tenantID := ctx.Value("tenant_id"); tenantID != nil {
		if str, ok := tenantID.(string); ok {
			return str, nil
		}
	}

	// 2. 从数据库字段获取
	if db.Statement.Schema != nil {
		if field := db.Statement.Schema.LookUpField(tenantIDField); field != nil {
			if fieldValue, ok := field.ValueOf(db.Statement.Context).(string); ok {
				return fieldValue, nil
			}
		}
	}

	return "", errors.New("tenant_id not found")
}

// EncryptRow 手动加密一行数据
func (tde *TDEMiddleware) EncryptRow(ctx context.Context, tableName string, data interface{}, tenantID string) error {
	// 查找表配置
	var config *TableEncryptionConfig
	for i := range tde.encryptionConfig.SensitiveFields {
		// TODO: 实现配置查找
		_ = i
	}

	if config == nil || !config.Enabled {
		return nil // 没有配置或未启用,不加密
	}

	// 加密数据
	return tde.fieldEncryptor.EncryptStruct(ctx, tenantID, data, tde.encryptionConfig)
}

// DecryptRow 手动解密一行数据
func (tde *TDEMiddleware) DecryptRow(ctx context.Context, tableName string, data interface{}, tenantID string) error {
	// 查找表配置
	var config *TableEncryptionConfig
	for i := range tde.encryptionConfig.SensitiveFields {
		// TODO: 实现配置查找
		_ = i
	}

	if config == nil || !config.Enabled {
		return nil // 没有配置或未启用,不解密
	}

	// 解密数据
	return tde.fieldEncryptor.DecryptStruct(ctx, tenantID, data, tde.encryptionConfig)
}

// EncryptBatch 批量加密
func (tde *TDEMiddleware) EncryptBatch(ctx context.Context, tableName string, items []interface{}, tenantID string) error {
	for i, item := range items {
		if err := tde.EncryptRow(ctx, tableName, item, tenantID); err != nil {
			return fmt.Errorf("failed to encrypt item %d: %w", i, err)
		}
	}
	return nil
}

// DecryptBatch 批量解密
func (tde *TDEMiddleware) DecryptBatch(ctx context.Context, tableName string, items []interface{}, tenantID string) error {
	for i, item := range items {
		if err := tde.DecryptRow(ctx, tableName, item, tenantID); err != nil {
			return fmt.Errorf("failed to decrypt item %d: %w", i, err)
		}
	}
	return nil
}

// ReencryptWithNewKey 使用新密钥重新加密(密钥轮换时使用)
func (tde *TDEMiddleware) ReencryptWithNewKey(ctx context.Context, tableName string, tenantID string) error {
	// TODO: 实现密钥轮换时的数据重新加密逻辑
	// 1. 读取所有加密数据
	// 2. 使用旧密钥解密
	// 3. 使用新密钥加密
	// 4. 更新数据库
	return errors.New("not implemented")
}

// GetTableEncryptionStatus 获取表加密状态
func (tde *TDEMiddleware) GetTableEncryptionStatus(tableName string) (*TableEncryptionStatus, error) {
	// TODO: 实现状态查询
	return &TableEncryptionStatus{
		TableName:        tableName,
		EncryptionEnabled: false,
		EncryptedFields:  []string{},
		TotalRows:        0,
		EncryptedRows:    0,
	}, nil
}

// TableEncryptionStatus 表加密状态
type TableEncryptionStatus struct {
	TableName         string   `json:"table_name"`
	EncryptionEnabled bool     `json:"encryption_enabled"`
	EncryptedFields   []string `json:"encrypted_fields"`
	TotalRows         int64    `json:"total_rows"`
	EncryptedRows     int64    `json:"encrypted_rows"`
	UnencryptedRows   int64    `json:"unencrypted_rows"`
	EncryptionRatio   float64  `json:"encryption_ratio"`
}

// VerifyEncryption 验证数据是否已加密
func (tde *TDEMiddleware) VerifyEncryption(data interface{}, field string) bool {
	// TODO: 实现加密验证逻辑
	return false
}
