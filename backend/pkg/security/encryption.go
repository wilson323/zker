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

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrInvalidKey 密钥无效
	ErrInvalidKey = errors.New("encryption key must be 32 bytes for AES-256")

	// ErrInvalidCiphertext 密文无效
	ErrInvalidCiphertext = errors.New("invalid ciphertext")

	// ErrCiphertextTooShort 密文过短
	ErrCiphertextTooShort = errors.New("ciphertext too short")
)

// EncryptionService 加密服务接口
type EncryptionService interface {
	// Encrypt 加密明文
	Encrypt(plaintext string) (string, error)

	// Decrypt 解密密文
	Decrypt(ciphertext string) (string, error)
}

// aesGCMEncryptor AES-256-GCM加密器
type aesGCMEncryptor struct {
	key []byte
}

// NewAESGCMEncryptor 创建AES-256-GCM加密器
// 密钥必须是32字节（AES-256）
func NewAESGCMEncryptor(key string) (*aesGCMEncryptor, error) {
	// 如果密钥不是32字节，使用SHA256哈希转换为32字节
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		hash := sha256.Sum256(keyBytes)
		keyBytes = hash[:]
	}

	return &aesGCMEncryptor{
		key: keyBytes,
	}, nil
}

// Encrypt 使用AES-256-GCM加密
// 返回Base64编码的密文（包含nonce）
func (e *aesGCMEncryptor) Encrypt(plaintext string) (string, error) {
	// 创建AES cipher
	block, err := aes.NewCipher(e.key)
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

	// 加密数据（GCM模式会自动添加认证标签）
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用AES-256-GCM解密
// ciphertext是Base64编码的密文
func (e *aesGCMEncryptor) Decrypt(ciphertext string) (string, error) {
	// Base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// 创建AES cipher
	block, err := aes.NewCipher(e.key)
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
		return "", ErrCiphertextTooShort
	}

	// 分离nonce和密文
	nonce, cipherData := data[:nonceSize], data[nonceSize:]

	// 解密并验证认证标签
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GenerateEncryptionKey 生成32字节的随机加密密钥
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate encryption key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// HashKeyTo32Bytes 将任意长度的密钥哈希为32字节
func HashKeyTo32Bytes(key string) string {
	hash := sha256.Sum256([]byte(key))
	return base64.StdEncoding.EncodeToString(hash[:])
}
