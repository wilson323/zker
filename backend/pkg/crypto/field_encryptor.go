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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"go.uber.org/zap"
)

// FieldEncryptor 字段加密器
type FieldEncryptor struct {
	keyManager KeyManager
	logger     *zap.Logger
	mu         sync.RWMutex
	cache      map[string]*cachedKey // 密钥缓存
}

type cachedKey struct {
	keyBytes  []byte
	expiresAt int64 // Unix timestamp
}

// FieldEncryptionConfig 字段加密配置
type FieldEncryptionConfig struct {
	// 敏感字段列表(支持通配符)
	SensitiveFields []string `json:"sensitive_fields"`

	// 加密算法
	Algorithm string `json:"algorithm"` // AES-256-GCM

	// 是否启用自动加密
	AutoEncrypt bool `json:"auto_encrypt"`

	// 是否启用批量加密
	BatchMode bool `json:"batch_mode"`

	// 批量大小
	BatchSize int `json:"batch_size"`
}

// DefaultEncryptionConfig 默认加密配置
func DefaultEncryptionConfig() *FieldEncryptionConfig {
	return &FieldEncryptionConfig{
		SensitiveFields: []string{
			"password", "passwd",
			"api_key", "apikey", "apiKey",
			"secret", "token", "access_token", "refresh_token",
			"private_key", "privatekey",
			"credit_card", "bank_account",
			"ssn", "social_security",
			"id_card", "idcard",
		},
		Algorithm:   "AES-256-GCM",
		AutoEncrypt: true,
		BatchMode:   true,
		BatchSize:   100,
	}
}

// NewFieldEncryptor 创建字段加密器
func NewFieldEncryptor(keyManager KeyManager, logger *zap.Logger) *FieldEncryptor {
	return &FieldEncryptor{
		keyManager: keyManager,
		logger:     logger,
		cache:      make(map[string]*cachedKey),
	}
}

// EncryptField 加密字段
func (fe *FieldEncryptor) EncryptField(ctx context.Context, tenantID, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// 获取加密密钥
	keyBytes, err := fe.getEncryptionKey(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to get encryption key: %w", err)
	}

	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64编码并添加版本标识
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return fmt.Sprintf("v1:%s", encoded), nil
}

// DecryptField 解密字段
func (fe *FieldEncryptor) DecryptField(ctx context.Context, tenantID, ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// 检查是否是加密格式
	if !strings.Contains(ciphertext, ":") {
		// 不是加密格式,直接返回
		return ciphertext, nil
	}

	// 分离版本和密文
	parts := strings.SplitN(ciphertext, ":", 2)
	if len(parts) != 2 {
		return ciphertext, nil // 格式不对,可能不是加密数据
	}

	version := parts[0]
	encodedCipher := parts[1]

	// 获取解密密钥
	keyBytes, err := fe.getEncryptionKey(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to get decryption key: %w", err)
	}

	// Base64解码
	data, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// 创建AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 获取nonce大小
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// 分离nonce和密文
	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		// 解密失败,可能密钥已轮换,尝试使用旧密钥
		return fe.decryptWithOldKeys(ctx, tenantID, data, err)
	}

	return string(plaintext), nil
}

// decryptWithOldKeys 使用旧密钥解密
func (fe *FieldEncryptor) decryptWithOldKeys(ctx context.Context, tenantID string, data []byte, originalErr error) (string, error) {
	// 获取所有密钥版本
	keys, err := fe.keyManager.ListKeys(ctx, tenantID)
	if err != nil {
		return "", fmt.Errorf("failed to list keys: %w (original error: %w)", err, originalErr)
	}

	// 尝试使用每个旧密钥解密
	for _, keyMeta := range keys {
		if keyMeta.Status != "active" && keyMeta.Status != "rotated" {
			continue
		}

		encKey, err := fe.keyManager.(*keyManagerImpl).GetKeyByID(ctx, keyMeta.KeyID)
		if err != nil {
			continue
		}

		keyBytes, err := fe.keyManager.(*keyManagerImpl).DecryptKey(encKey)
		if err != nil {
			continue
		}

		// 尝试解密
		block, err := aes.NewCipher(keyBytes)
		if err != nil {
			continue
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			continue
		}

		nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
			continue
		}

		nonce, cipherData := data[:nonceSize], data[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
		if err == nil {
			// 解密成功,记录日志
			fe.logger.Info("Decrypted with old key",
				zap.String("tenant_id", tenantID),
				zap.String("key_version", fmt.Sprintf("%d", keyMeta.Version)),
			)
			return string(plaintext), nil
		}
	}

	return "", fmt.Errorf("failed to decrypt with any key: %w", originalErr)
}

// EncryptStruct 加密结构体中的敏感字段
func (fe *FieldEncryptor) EncryptStruct(ctx context.Context, tenantID string, data interface{}, config *FieldEncryptionConfig) error {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("input must be a struct or pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 检查是否应该加密
		if !fe.shouldEncryptField(field.Name, config) {
			continue
		}

		// 只处理字符串字段
		if fieldValue.Kind() != reflect.String {
			continue
		}

		plaintext := fieldValue.String()
		if plaintext == "" {
			continue
		}

		// 加密
		ciphertext, err := fe.EncryptField(ctx, tenantID, plaintext)
		if err != nil {
			return fmt.Errorf("failed to encrypt field %s: %w", field.Name, err)
		}

		fieldValue.SetString(ciphertext)
	}

	return nil
}

// DecryptStruct 解密结构体中的敏感字段
func (fe *FieldEncryptor) DecryptStruct(ctx context.Context, tenantID string, data interface{}, config *FieldEncryptionConfig) error {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("input must be a struct or pointer to struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// 检查是否应该解密
		if !fe.shouldEncryptField(field.Name, config) {
			continue
		}

		// 只处理字符串字段
		if fieldValue.Kind() != reflect.String {
			continue
		}

		ciphertext := fieldValue.String()
		if ciphertext == "" {
			continue
		}

		// 解密
		plaintext, err := fe.DecryptField(ctx, tenantID, ciphertext)
		if err != nil {
			fe.logger.Warn("Failed to decrypt field",
				zap.String("field", field.Name),
				zap.Error(err),
			)
			// 解密失败不解密,保留原值
			continue
		}

		fieldValue.SetString(plaintext)
	}

	return nil
}

// EncryptBatch 批量加密
func (fe *FieldEncryptor) EncryptBatch(ctx context.Context, tenantID string, items []interface{}, config *FieldEncryptionConfig) error {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	for i, item := range items {
		if err := fe.EncryptStruct(ctx, tenantID, item, config); err != nil {
			return fmt.Errorf("failed to encrypt item %d: %w", i, err)
		}
	}

	return nil
}

// DecryptBatch 批量解密
func (fe *FieldEncryptor) DecryptBatch(ctx context.Context, tenantID string, items []interface{}, config *FieldEncryptionConfig) error {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	for i, item := range items {
		if err := fe.DecryptStruct(ctx, tenantID, item, config); err != nil {
			return fmt.Errorf("failed to decrypt item %d: %w", i, err)
		}
	}

	return nil
}

// EncryptJSONMap 加密JSON Map中的敏感字段
func (fe *FieldEncryptor) EncryptJSONMap(ctx context.Context, tenantID string, data map[string]interface{}, config *FieldEncryptionConfig) (map[string]interface{}, error) {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	result := make(map[string]interface{})
	for key, value := range data {
		if !fe.shouldEncryptField(key, config) {
			result[key] = value
			continue
		}

		// 处理字符串值
		if strValue, ok := value.(string); ok && strValue != "" {
			encrypted, err := fe.EncryptField(ctx, tenantID, strValue)
			if err != nil {
				return nil, fmt.Errorf("failed to encrypt field %s: %w", key, err)
			}
			result[key] = encrypted
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// DecryptJSONMap 解密JSON Map中的敏感字段
func (fe *FieldEncryptor) DecryptJSONMap(ctx context.Context, tenantID string, data map[string]interface{}, config *FieldEncryptionConfig) (map[string]interface{}, error) {
	if config == nil {
		config = DefaultEncryptionConfig()
	}

	result := make(map[string]interface{})
	for key, value := range data {
		if !fe.shouldEncryptField(key, config) {
			result[key] = value
			continue
		}

		// 处理字符串值
		if strValue, ok := value.(string); ok && strValue != "" {
			decrypted, err := fe.DecryptField(ctx, tenantID, strValue)
			if err != nil {
				fe.logger.Warn("Failed to decrypt field",
					zap.String("field", key),
					zap.Error(err),
				)
				result[key] = value // 保留原值
			} else {
				result[key] = decrypted
			}
		} else {
			result[key] = value
		}
	}

	return result, nil
}

// EncryptJSONBytes 加密JSON字节数组
func (fe *FieldEncryptor) EncryptJSONBytes(ctx context.Context, tenantID string, data []byte, config *FieldEncryptionConfig) ([]byte, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	encrypted, err := fe.EncryptJSONMap(ctx, tenantID, jsonMap, config)
	if err != nil {
		return nil, err
	}

	return json.Marshal(encrypted)
}

// DecryptJSONBytes 解密JSON字节数组
func (fe *FieldEncryptor) DecryptJSONBytes(ctx context.Context, tenantID string, data []byte, config *FieldEncryptionConfig) ([]byte, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	decrypted, err := fe.DecryptJSONMap(ctx, tenantID, jsonMap, config)
	if err != nil {
		return nil, err
	}

	return json.Marshal(decrypted)
}

// shouldEncryptField 判断字段是否应该加密
func (fe *FieldEncryptor) shouldEncryptField(fieldName string, config *FieldEncryptionConfig) bool {
	fieldName = strings.ToLower(fieldName)

	for _, pattern := range config.SensitiveFields {
		pattern = strings.ToLower(pattern)
		// 精确匹配
		if fieldName == pattern {
			return true
		}
		// 包含匹配
		if strings.Contains(fieldName, pattern) {
			return true
		}
	}

	return false
}

// getEncryptionKey 获取加密密钥(带缓存)
func (fe *FieldEncryptor) getEncryptionKey(ctx context.Context, tenantID string) ([]byte, error) {
	// 检查缓存
	fe.mu.RLock()
	cached, exists := fe.cache[tenantID]
	fe.mu.RUnlock()

	if exists {
		// 检查是否过期
		if cached.expiresAt > 0 && cached.expiresAt < getCurrentTimestamp() {
			// 缓存过期,删除
			fe.mu.Lock()
			delete(fe.cache, tenantID)
			fe.mu.Unlock()
		} else {
			// 缓存有效
			return cached.keyBytes, nil
		}
	}

	// 获取活跃密钥
	encKey, err := fe.keyManager.GetActiveKey(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 解密密钥
	keyBytes, err := fe.keyManager.(*keyManagerImpl).DecryptKey(encKey)
	if err != nil {
		return nil, err
	}

	// 缓存密钥(5分钟)
	fe.mu.Lock()
	fe.cache[tenantID] = &cachedKey{
		keyBytes:  keyBytes,
		expiresAt: getCurrentTimestamp() + 300,
	}
	fe.mu.Unlock()

	return keyBytes, nil
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() int64 {
	return 0 // TODO: 实现获取当前时间戳
}

// IsEncrypted 检查字符串是否已加密
func IsEncrypted(data string) bool {
	return strings.HasPrefix(data, "v1:") || strings.HasPrefix(data, "v2:")
}

// MaskField 快速脱敏字段(用于日志)
func MaskField(value string, showFirst, showLast int) string {
	if value == "" {
		return ""
	}

	runes := []rune(value)
	length := len(runes)

	if length <= showFirst+showLast {
		return strings.Repeat("*", length)
	}

	maskedLength := length - showFirst - showLast
	return string(runes[:showFirst]) + strings.Repeat("*", maskedLength) + string(runes[length-showLast:])
}
