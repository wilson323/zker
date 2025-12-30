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

package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encryptor 加密器接口
type Encryptor interface {
	// Encrypt 加密明文
	Encrypt(plaintext string) (string, error)
	// Decrypt 解密密文
	Decrypt(ciphertext string) (string, error)
}

// aesEncryptor AES加密器实现
type aesEncryptor struct {
	key []byte
}

// NewAESEncryptor 创建AES加密器
// keySize: 密钥大小，16(AES-128)、24(AES-192)或32(AES-256)字节
func NewAESEncryptor(keySize int) (Encryptor, error) {
	key, err := GenerateKey(keySize)
	if err != nil {
		return nil, err
	}
	return &aesEncryptor{key: key}, nil
}

// NewAESEncryptorWithKey 使用指定密钥创建AES加密器
func NewAESEncryptorWithKey(key []byte) (Encryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("key size must be 16, 24, or 32 bytes")
	}
	return &aesEncryptor{key: key}, nil
}

// Encrypt 加密明文
func (e *aesEncryptor) Encrypt(plaintext string) (string, error) {
	return AESEncrypt(e.key, plaintext)
}

// Decrypt 解密密文
func (e *aesEncryptor) Decrypt(ciphertext string) (string, error) {
	return AESDecrypt(e.key, ciphertext)
}

// AESEncrypt AES加密
// key: 加密密钥（16/24/32字节对应AES-128/192/256）
// plaintext: 明文
// 返回：Base64编码的密文
func AESEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 使用GCM模式（Galois/Counter Mode）
	// GCM提供认证加密，同时保证机密性和完整性
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密并认证
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecrypt AES解密
// key: 解密密钥
// ciphertext: Base64编码的密文
// 返回：明文
func AESDecrypt(key []byte, ciphertext string) (string, error) {
	// Base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 检查长度
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// 提取nonce和密文
	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// 解密并验证
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GenerateKey 生成随机密钥
// size: 密钥大小（16/24/32字节）
func GenerateKey(size int) ([]byte, error) {
	if size != 16 && size != 24 && size != 32 {
		return nil, errors.New("key size must be 16, 24, or 32 bytes")
	}

	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	return key, nil
}
