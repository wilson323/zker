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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 密钥管理错误
var (
	ErrKeyNotFound       = errors.New("encryption key not found")
	ErrKeyExpired        = errors.New("encryption key has expired")
	ErrInvalidKeyVersion = errors.New("invalid key version")
	ErrKeyRotationFailed = errors.New("key rotation failed")
)

// KeyMetadata 密钥元数据
type KeyMetadata struct {
	KeyID      string    `json:"key_id" gorm:"primaryKey"`
	Version    int       `json:"version" gorm:"index"`
	KeyType    string    `json:"key_type"` // AES-256-GCM, RSA-2048, etc.
	Status     string    `json:"status"`    // active, rotated, revoked
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	RotatedAt  *time.Time `json:"rotated_at"`
	RotatedTo  *string   `json:"rotated_to"` // 新版本的key_id
	CreatedBy  string    `json:"created_by"`
	TenantID   string    `json:"tenant_id" gorm:"index"`
	Algorithm  string    `json:"algorithm"`  // AES-256-GCM
	KeySize    int       `json:"key_size"`   // 256, 512
	UsageCount int64     `json:"usage_count"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

// KeyUsageLog 密钥使用日志
type KeyUsageLog struct {
	LogID      string    `json:"log_id" gorm:"primaryKey"`
	KeyID      string    `json:"key_id" gorm:"index"`
	Version    int       `json:"version"`
	Action     string    `json:"action"` // encrypt, decrypt, rotate
	UserID     string    `json:"user_id"`
	TenantID   string    `json:"tenant_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Success    bool      `json:"success"`
	ErrorMsg   string    `json:"error_msg,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	ResourceID string    `json:"resource_id,omitempty"`
}

// KeyRotationConfig 密钥轮换配置
type KeyRotationConfig struct {
	AutoRotate         bool          `json:"auto_rotate"`          // 是否自动轮换
	RotationInterval   time.Duration `json:"rotation_interval"`    // 轮换间隔
	RotationReminder   time.Duration `json:"rotation_reminder"`    // 轮换提醒
	GracePeriod        time.Duration `json:"grace_period"`         // 宽限期
	KeepOldKeyDuration time.Duration `json:"keep_old_key_duration"` // 保留旧密钥时长
}

// KeyManager 密钥管理器接口
type KeyManager interface {
	// GenerateKey 生成新密钥
	GenerateKey(ctx context.Context, keyType string, tenantID, createdBy string) (*KeyMetadata, error)

	// GetActiveKey 获取当前活跃密钥
	GetActiveKey(ctx context.Context, tenantID string) (*EncryptionKey, error)

	// GetKeyByID 根据ID获取密钥
	GetKeyByID(ctx context.Context, keyID string) (*EncryptionKey, error)

	// GetKeyByVersion 根据版本获取密钥
	GetKeyByVersion(ctx context.Context, tenantID string, version int) (*EncryptionKey, error)

	// RotateKey 轮换密钥
	RotateKey(ctx context.Context, tenantID, rotatedBy string) (*KeyMetadata, error)

	// RevokeKey 撤销密钥
	RevokeKey(ctx context.Context, keyID, revokedBy string) error

	// ListKeys 列出租户的所有密钥版本
	ListKeys(ctx context.Context, tenantID string) ([]*KeyMetadata, error)

	// LogUsage 记录密钥使用
	LogUsage(ctx context.Context, log *KeyUsageLog) error

	// GetRotationStatus 获取密钥轮换状态
	GetRotationStatus(ctx context.Context, tenantID string) (*RotationStatus, error)
}

// EncryptionKey 加密密钥
type EncryptionKey struct {
	Metadata *KeyMetadata
	KeyBytes []byte // 加密后的密钥内容
}

// RotationStatus 轮换状态
type RotationStatus struct {
	CurrentVersion      int        `json:"current_version"`
	LastRotationAt      *time.Time `json:"last_rotation_at"`
	NextRotationAt      *time.Time `json:"next_rotation_at"`
	DaysUntilRotation   int        `json:"days_until_rotation"`
	RotationRequired    bool       `json:"rotation_required"`
	RotationConfig      KeyRotationConfig `json:"rotation_config"`
}

// keyManagerImpl 密钥管理器实现
type keyManagerImpl struct {
	db              *gorm.DB
	logger          *zap.Logger
	masterKey       []byte // 主密钥(用于加密其他密钥)
	mu              sync.RWMutex
	rotationConfig  KeyRotationConfig
	cache           map[string]*EncryptionKey // 密钥缓存
}

// NewKeyManager 创建密钥管理器
func NewKeyManager(db *gorm.DB, logger *zap.Logger, masterKey string) (KeyManager, error) {
	// 从环境变量或配置读取主密钥
	if masterKey == "" {
		return nil, errors.New("master key is required")
	}

	// 主密钥必须是32字节(AES-256)
	masterKeyBytes := []byte(masterKey)
	if len(masterKeyBytes) != 32 {
		hash := sha256.Sum256(masterKeyBytes)
		masterKeyBytes = hash[:]
	}

	km := &keyManagerImpl{
		db:             db,
		logger:         logger,
		masterKey:      masterKeyBytes,
		rotationConfig: DefaultRotationConfig(),
		cache:          make(map[string]*EncryptionKey),
	}

	// 初始化数据库表
	if err := km.initTables(); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	return km, nil
}

// DefaultRotationConfig 默认轮换配置
func DefaultRotationConfig() KeyRotationConfig {
	return KeyRotationConfig{
		AutoRotate:         true,
		RotationInterval:   90 * 24 * time.Hour, // 90天
		RotationReminder:   7 * 24 * time.Hour,  // 提前7天提醒
		GracePeriod:        30 * 24 * time.Hour, // 30天宽限期
		KeepOldKeyDuration: 90 * 24 * time.Hour, // 保留旧密钥90天
	}
}

// initTables 初始化数据库表
func (km *keyManagerImpl) initTables() error {
	return km.db.AutoMigrate(&KeyMetadata{}, &KeyUsageLog{})
}

// GenerateKey 生成新密钥
func (km *keyManagerImpl) GenerateKey(ctx context.Context, keyType string, tenantID, createdBy string) (*KeyMetadata, error) {
	// 生成32字节随机密钥
	keyBytes := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// 使用主密钥加密密钥
	encryptedKey, err := km.encryptWithMasterKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	// 获取当前版本号
	var latestVersion int
	km.db.WithContext(ctx).Model(&KeyMetadata{}).
		Where("tenant_id = ?", tenantID).
		Order("version DESC").
		Select("version").
		Scan(&latestVersion)

	// 创建密钥元数据
	metadata := &KeyMetadata{
		KeyID:      km.generateKeyID(tenantID, latestVersion+1),
		Version:    latestVersion + 1,
		KeyType:    keyType,
		Status:     "active",
		CreatedAt:  time.Now(),
		CreatedBy:  createdBy,
		TenantID:   tenantID,
		Algorithm:  "AES-256-GCM",
		KeySize:    256,
		UsageCount: 0,
	}

	// 计算过期时间
	if km.rotationConfig.AutoRotate {
		expiresAt := time.Now().Add(km.rotationConfig.RotationInterval)
		metadata.ExpiresAt = &expiresAt
	}

	// 保存到数据库
	if err := km.db.WithContext(ctx).Create(metadata).Error; err != nil {
		return nil, fmt.Errorf("failed to save key metadata: %w", err)
	}

	// 缓存密钥
	km.mu.Lock()
	km.cache[metadata.KeyID] = &EncryptionKey{
		Metadata: metadata,
		KeyBytes: encryptedKey,
	}
	km.mu.Unlock()

	km.logger.Info("Generated new encryption key",
		zap.String("key_id", metadata.KeyID),
		zap.String("tenant_id", tenantID),
		zap.Int("version", metadata.Version),
	)

	return metadata, nil
}

// GetActiveKey 获取当前活跃密钥
func (km *keyManagerImpl) GetActiveKey(ctx context.Context, tenantID string) (*EncryptionKey, error) {
	var metadata KeyMetadata
	err := km.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, "active").
		Order("version DESC").
		First(&metadata).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有密钥，自动生成一个
			newKey, err := km.GenerateKey(ctx, "AES-256-GCM", tenantID, "system")
			if err != nil {
				return nil, fmt.Errorf("failed to generate key: %w", err)
			}
			// 递归获取
			return km.GetActiveKey(ctx, tenantID)
		}
		return nil, fmt.Errorf("failed to get active key: %w", err)
	}

	// 检查是否过期
	if metadata.ExpiresAt != nil && time.Now().After(*metadata.ExpiresAt) {
		// 密钥已过期，触发轮换
		km.logger.Warn("Active key has expired, triggering rotation",
			zap.String("key_id", metadata.KeyID),
			zap.Time("expired_at", *metadata.ExpiresAt),
		)
		_, _ = km.RotateKey(ctx, tenantID, "system-auto-rotation")
		return km.GetActiveKey(ctx, tenantID)
	}

	return km.loadKey(ctx, &metadata)
}

// GetKeyByID 根据ID获取密钥
func (km *keyManagerImpl) GetKeyByID(ctx context.Context, keyID string) (*EncryptionKey, error) {
	// 先查缓存
	km.mu.RLock()
	if cachedKey, exists := km.cache[keyID]; exists {
		km.mu.RUnlock()
		return cachedKey, nil
	}
	km.mu.RUnlock()

	var metadata KeyMetadata
	err := km.db.WithContext(ctx).Where("key_id = ?", keyID).First(&metadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return km.loadKey(ctx, &metadata)
}

// GetKeyByVersion 根据版本获取密钥
func (km *keyManagerImpl) GetKeyByVersion(ctx context.Context, tenantID string, version int) (*EncryptionKey, error) {
	var metadata KeyMetadata
	err := km.db.WithContext(ctx).
		Where("tenant_id = ? AND version = ?", tenantID, version).
		First(&metadata).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return km.loadKey(ctx, &metadata)
}

// loadKey 加载密钥(从数据库或缓存)
func (km *keyManagerImpl) loadKey(ctx context.Context, metadata *KeyMetadata) (*EncryptionKey, error) {
	// 检查缓存
	km.mu.RLock()
	if cachedKey, exists := km.cache[metadata.KeyID]; exists {
		km.mu.RUnlock()
		return cachedKey, nil
	}
	km.mu.RUnlock()

	// TODO: 从密钥存储服务加载加密的密钥
	// 这里简化为直接在内存中
	encryptedKey, err := km.getEncryptedKeyFromStorage(ctx, metadata.KeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to load encrypted key: %w", err)
	}

	encKey := &EncryptionKey{
		Metadata: metadata,
		KeyBytes: encryptedKey,
	}

	// 缓存密钥
	km.mu.Lock()
	km.cache[metadata.KeyID] = encKey
	km.mu.Unlock()

	return encKey, nil
}

// RotateKey 轮换密钥
func (km *keyManagerImpl) RotateKey(ctx context.Context, tenantID, rotatedBy string) (*KeyMetadata, error) {
	// 获取当前活跃密钥
	oldKey, err := km.GetActiveKey(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active key: %w", err)
	}

	// 生成新密钥
	newMetadata, err := km.GenerateKey(ctx, "AES-256-GCM", tenantID, rotatedBy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new key: %w", err)
	}

	// 更新旧密钥状态
	now := time.Now()
	oldKey.Metadata.Status = "rotated"
	oldKey.Metadata.RotatedAt = &now
	oldKey.Metadata.RotatedTo = &newMetadata.KeyID

	if err := km.db.WithContext(ctx).Save(oldKey.Metadata).Error; err != nil {
		return nil, fmt.Errorf("failed to update old key status: %w", err)
	}

	// 更新缓存
	km.mu.Lock()
	delete(km.cache, oldKey.Metadata.KeyID)
	km.mu.Unlock()

	// 记录日志
	km.logger.Info("Rotated encryption key",
		zap.String("tenant_id", tenantID),
		zap.String("old_key_id", oldKey.Metadata.KeyID),
		zap.String("new_key_id", newMetadata.KeyID),
		zap.String("rotated_by", rotatedBy),
	)

	return newMetadata, nil
}

// RevokeKey 撤销密钥
func (km *keyManagerImpl) RevokeKey(ctx context.Context, keyID, revokedBy string) error {
	var metadata KeyMetadata
	err := km.db.WithContext(ctx).Where("key_id = ?", keyID).First(&metadata).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrKeyNotFound
		}
		return fmt.Errorf("failed to get key: %w", err)
	}

	// 更新状态
	metadata.Status = "revoked"
	if err := km.db.WithContext(ctx).Save(&metadata).Error; err != nil {
		return fmt.Errorf("failed to revoke key: %w", err)
	}

	// 从缓存移除
	km.mu.Lock()
	delete(km.cache, keyID)
	km.mu.Unlock()

	km.logger.Info("Revoked encryption key",
		zap.String("key_id", keyID),
		zap.String("revoked_by", revokedBy),
	)

	return nil
}

// ListKeys 列出租户的所有密钥版本
func (km *keyManagerImpl) ListKeys(ctx context.Context, tenantID string) ([]*KeyMetadata, error) {
	var keys []*KeyMetadata
	err := km.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("version DESC").
		Find(&keys).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	return keys, nil
}

// LogUsage 记录密钥使用
func (km *keyManagerImpl) LogUsage(ctx context.Context, log *KeyUsageLog) error {
	log.LogID = km.generateLogID()
	log.CreatedAt = time.Now()

	if err := km.db.WithContext(ctx).Create(log).Error; err != nil {
		km.logger.Warn("Failed to log key usage", zap.Error(err))
		// 不返回错误，日志失败不应影响主流程
	}

	// 更新使用计数
	km.db.WithContext(ctx).Model(&KeyMetadata{}).
		Where("key_id = ?", log.KeyID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))

	// 更新最后使用时间
	now := time.Now()
	km.db.WithContext(ctx).Model(&KeyMetadata{}).
		Where("key_id = ?", log.KeyID).
		Update("last_used_at", &now)

	return nil
}

// GetRotationStatus 获取密钥轮换状态
func (km *keyManagerImpl) GetRotationStatus(ctx context.Context, tenantID string) (*RotationStatus, error) {
	var metadata KeyMetadata
	err := km.db.WithContext(ctx).
		Where("tenant_id = ? AND status = ?", tenantID, "active").
		Order("version DESC").
		First(&metadata).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active key: %w", err)
	}

	status := &RotationStatus{
		CurrentVersion: metadata.Version,
		LastRotationAt: metadata.CreatedAt,
		RotationConfig: km.rotationConfig,
	}

	// 计算下次轮换时间
	if metadata.ExpiresAt != nil {
		nextRotation := metadata.ExpiresAt
		status.NextRotationAt = nextRotation
		status.DaysUntilRotation = int(nextRotation.Sub(time.Now()).Hours() / 24)

		// 如果距离轮换时间小于提醒时间，标记为需要轮换
		if time.Until(*nextRotation) < km.rotationConfig.RotationReminder {
			status.RotationRequired = true
		}
	}

	return status, nil
}

// encryptWithMasterKey 使用主密钥加密
func (km *keyManagerImpl) encryptWithMasterKey(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decryptWithMasterKey 使用主密钥解密
func (km *keyManagerImpl) decryptWithMasterKey(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, cipherData := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// getEncryptedKeyFromStorage 从存储获取加密密钥
// TODO: 实现从安全存储(如HashiCorp Vault)获取密钥
func (km *keyManagerImpl) getEncryptedKeyFromStorage(ctx context.Context, keyID string) ([]byte, error) {
	// 简化实现:从环境变量或配置文件读取
	// 生产环境应该使用专门的密钥管理服务
	encryptedKeyStr := "" // TODO: 从存储读取
	if encryptedKeyStr == "" {
		return nil, errors.New("encrypted key not found in storage")
	}

	return base64.StdEncoding.DecodeString(encryptedKeyStr)
}

// generateKeyID 生成密钥ID
func (km *keyManagerImpl) generateKeyID(tenantID string, version int) string {
	data := fmt.Sprintf("%s-%d-%d", tenantID, version, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("key-%s-%d", base64.URLEncoding.EncodeToString(hash[:8]), version)
}

// generateLogID 生成日志ID
func (km *keyManagerImpl) generateLogID() string {
	data := fmt.Sprintf("%d", time.Now().UnixNano())
	hash := sha256.Sum256([]byte(data))
	return base64.URLEncoding.EncodeToString(hash[:8])
}

// DecryptKey 解密密钥(用于实际加密数据)
func (km *keyManagerImpl) DecryptKey(encKey *EncryptionKey) ([]byte, error) {
	// 从内存缓存或存储加载加密的密钥
	encryptedKey := encKey.KeyBytes

	// 使用主密钥解密
	keyBytes, err := km.decryptWithMasterKey(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return keyBytes, nil
}

// EncryptKey 加密密钥(用于存储)
func (km *keyManagerImpl) EncryptKey(keyBytes []byte) ([]byte, error) {
	return km.encryptWithMasterKey(keyBytes)
}
