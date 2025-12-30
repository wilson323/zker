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
	"strings"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "valid AES-128 key",
			size:    16,
			wantErr: false,
		},
		{
			name:    "valid AES-192 key",
			size:    24,
			wantErr: false,
		},
		{
			name:    "valid AES-256 key",
			size:    32,
			wantErr: false,
		},
		{
			name:    "invalid key size",
			size:    12,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(key) != tt.size {
				t.Errorf("GenerateKey() key length = %d, want %d", len(key), tt.size)
			}
		})
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("Failed to generate key: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "simple text",
			plaintext: "Hello, World!",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("A", 1000),
		},
		{
			name:      "special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "unicode",
			plaintext: "你好世界🌍🎉",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试加密
			ciphertext, err := AESEncrypt(key, tt.plaintext)
			if err != nil {
				t.Fatalf("Failed to encrypt: %v", err)
			}

			// 密文不应与明文相同
			if ciphertext == tt.plaintext && tt.plaintext != "" {
				t.Error("Ciphertext should not equal plaintext")
			}

			// 测试解密
			decrypted, err := AESDecrypt(key, ciphertext)
			if err != nil {
				t.Fatalf("Failed to decrypt: %v", err)
			}

			// 解密后的文本应与原文相同
			if decrypted != tt.plaintext {
				t.Errorf("Decrypted text doesn't match. Got: %s, Want: %s", decrypted, tt.plaintext)
			}
		})
	}
}

func TestAESDecryptWithWrongKey(t *testing.T) {
	key1, _ := GenerateKey(32)
	key2, _ := GenerateKey(32)

	plaintext := "Secret Message"
	ciphertext, _ := AESEncrypt(key1, plaintext)

	// 使用错误的密钥解密应该失败
	_, err := AESDecrypt(key2, ciphertext)
	if err == nil {
		t.Error("Decrypting with wrong key should fail")
	}
}

func TestAESDecryptInvalidBase64(t *testing.T) {
	key, _ := GenerateKey(32)

	// 无效的Base64字符串
	_, err := AESDecrypt(key, "invalid-base64!!!")
	if err == nil {
		t.Error("Decoding invalid base64 should fail")
	}
}

func TestAESDecryptTooShort(t *testing.T) {
	key, _ := GenerateKey(32)

	// 使用有效的Base64但数据太短
	_, err := AESDecrypt(key, "YWJj") // "abc" in base64
	if err == nil {
		t.Error("Decrypting too short ciphertext should fail")
	}
}

func TestNewAESEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		keySize int
		wantErr bool
	}{
		{
			name:    "valid AES-128",
			keySize: 16,
			wantErr: false,
		},
		{
			name:    "valid AES-256",
			keySize: 32,
			wantErr: false,
		},
		{
			name:    "invalid size",
			keySize: 12,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewAESEncryptor(tt.keySize)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encryptor == nil {
				t.Error("NewAESEncryptor() should return non-nil encryptor")
			}
		})
	}
}

func TestNewAESEncryptorWithKey(t *testing.T) {
	tests := []struct {
		name    string
		key     []byte
		wantErr bool
	}{
		{
			name:    "valid 16-byte key",
			key:     make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "valid 32-byte key",
			key:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "invalid key length",
			key:     make([]byte, 12),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewAESEncryptorWithKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptorWithKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && encryptor == nil {
				t.Error("NewAESEncryptorWithKey() should return non-nil encryptor")
			}
		})
	}
}

func TestAESEncryptorInterface(t *testing.T) {
	key, _ := GenerateKey(32)
	encryptor, _ := NewAESEncryptorWithKey(key)

	plaintext := "Test message for interface"

	// 测试Encrypt方法
	ciphertext, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encryptor.Encrypt() failed: %v", err)
	}

	// 测试Decrypt方法
	decrypted, err := encryptor.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Encryptor.Decrypt() failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("Decrypted text doesn't match. Got: %s, Want: %s", decrypted, plaintext)
	}
}

// 基准测试
func BenchmarkAESEncrypt(b *testing.B) {
	key, _ := GenerateKey(32)
	plaintext := "Hello, World!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AESEncrypt(key, plaintext)
	}
}

func BenchmarkAESDecrypt(b *testing.B) {
	key, _ := GenerateKey(32)
	plaintext := "Hello, World!"
	ciphertext, _ := AESEncrypt(key, plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AESDecrypt(key, ciphertext)
	}
}
