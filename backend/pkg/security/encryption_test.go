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
	"strings"
	"testing"
)

func TestNewAESGCMEncryptor_ValidKey(t *testing.T) {
	// 测试32字节密钥
	key32 := strings.Repeat("a", 32)
	encryptor, err := NewAESGCMEncryptor(key32)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}
	if encryptor == nil {
		t.Fatal("encryptor should not be nil")
	}

	// 测试非32字节密钥（应该被哈希为32字节）
	keyShort := "short_key"
	encryptor, err = NewAESGCMEncryptor(keyShort)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor with short key failed: %v", err)
	}
	if encryptor == nil {
		t.Fatal("encryptor should not be nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := strings.Repeat("b", 32)
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	testCases := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "Simple text",
			plaintext: "Hello, World!",
		},
		{
			name:      "API Key",
			plaintext: "sk_test_1234567890abcdef",
		},
		{
			name:      "Empty string",
			plaintext: "",
		},
		{
			name:      "Special characters",
			plaintext: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name:      "Long text",
			plaintext: strings.Repeat("a", 1000),
		},
		{
			name:      "Unicode",
			plaintext: "中文测试🔒",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 加密
			ciphertext, err := encryptor.Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// 密文不应该为空
			if ciphertext == "" {
				t.Fatal("ciphertext should not be empty")
			}

			// 密文应该不同于明文
			if ciphertext == tc.plaintext {
				t.Error("ciphertext should be different from plaintext")
			}

			// 解密
			decrypted, err := encryptor.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// 解密结果应该等于原文
			if decrypted != tc.plaintext {
				t.Errorf("decrypted text mismatch: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestEncrypt_DifferentEachTime(t *testing.T) {
	key := strings.Repeat("c", 32)
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	plaintext := "same plaintext"

	// 加密两次
	ciphertext1, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First Encrypt failed: %v", err)
	}

	ciphertext2, err := encryptor.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second Encrypt failed: %v", err)
	}

	// 由于使用随机nonce，两次加密结果应该不同
	if ciphertext1 == ciphertext2 {
		t.Error("encrypting the same plaintext twice should produce different ciphertexts")
	}

	// 但解密后应该得到相同的结果
	decrypted1, err := encryptor.Decrypt(ciphertext1)
	if err != nil {
		t.Fatalf("Decrypt ciphertext1 failed: %v", err)
	}

	decrypted2, err := encryptor.Decrypt(ciphertext2)
	if err != nil {
		t.Fatalf("Decrypt ciphertext2 failed: %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("decrypted text should match original plaintext")
	}
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := strings.Repeat("d", 32)
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	testCases := []struct {
		name       string
		ciphertext string
		wantErr    error
	}{
		{
			name:       "Empty string",
			ciphertext: "",
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "Invalid base64",
			ciphertext: "not-base64!!!",
			wantErr:    ErrInvalidCiphertext,
		},
		{
			name:       "Too short",
			ciphertext: "YWJj", // base64 of "abc"
			wantErr:    ErrCiphertextTooShort,
		},
		{
			name:       "Valid base64 but invalid ciphertext",
			ciphertext: "dGVzdHRlc3R0ZXN0dGVzdA==", // base64 of "testtesttesttest"
			wantErr:    ErrInvalidCiphertext,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := encryptor.Decrypt(tc.ciphertext)
			if err == nil {
				t.Fatal("Decrypt should return error")
			}
			if !strings.Contains(err.Error(), tc.wantErr.Error()) &&
				!strings.Contains(err.Error(), "failed to") {
				t.Errorf("error should contain %q, got %q", tc.wantErr, err)
			}
		})
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	key1 := strings.Repeat("e", 32)
	key2 := strings.Repeat("f", 32)

	encryptor1, err := NewAESGCMEncryptor(key1)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	encryptor2, err := NewAESGCMEncryptor(key2)
	if err != nil {
		t.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	plaintext := "secret message"

	// 使用key1加密
	ciphertext, err := encryptor1.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 使用key2解密（应该失败）
	_, err = encryptor2.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("Decrypt with wrong key should fail")
	}
}

func TestGenerateEncryptionKey(t *testing.T) {
	// 生成两个密钥
	key1, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey failed: %v", err)
	}

	key2, err := GenerateEncryptionKey()
	if err != nil {
		t.Fatalf("GenerateEncryptionKey failed: %v", err)
	}

	// 两个密钥应该不同
	if key1 == key2 {
		t.Error("generated keys should be different")
	}

	// 密钥应该是44字符（base64编码的32字节）
	if len(key1) != 44 {
		t.Errorf("key length should be 44, got %d", len(key1))
	}

	// 使用生成的密钥创建加密器
	_, err = NewAESGCMEncryptor(key1)
	if err != nil {
		t.Errorf("NewAESGCMEncryptor with generated key failed: %v", err)
	}
}

func TestHashKeyTo32Bytes(t *testing.T) {
	testCases := []struct {
		name string
		key  string
	}{
		{name: "Short key", key: "short"},
		{name: "Medium key", key: "medium_length_key"},
		{name: "Long key", key: strings.Repeat("x", 100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashed := HashKeyTo32Bytes(tc.key)

			// 哈希后的密钥应该是44字符（base64编码的32字节）
			if len(hashed) != 44 {
				t.Errorf("hashed key length should be 44, got %d", len(hashed))
			}

			// 相同的输入应该产生相同的哈希
			hashed2 := HashKeyTo32Bytes(tc.key)
			if hashed != hashed2 {
				t.Error("same input should produce same hash")
			}

			// 不同的输入应该产生不同的哈希
			differentKey := tc.key + "_different"
			hashed3 := HashKeyTo32Bytes(differentKey)
			if hashed == hashed3 {
				t.Error("different inputs should produce different hashes")
			}
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	key := strings.Repeat("g", 32)
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		b.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	plaintext := "benchmark_test_plaintext"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encryptor.Encrypt(plaintext)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key := strings.Repeat("h", 32)
	encryptor, err := NewAESGCMEncryptor(key)
	if err != nil {
		b.Fatalf("NewAESGCMEncryptor failed: %v", err)
	}

	plaintext := "benchmark_test_plaintext"
	ciphertext, _ := encryptor.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = encryptor.Decrypt(ciphertext)
	}
}
