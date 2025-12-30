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

package encrypt_test

import (
	"fmt"

	"github.com/coze-dev/coze-studio/backend/pkg/encrypt"
)

// Example usage similar to mfa_service.go
func ExampleEncryptor() {
	// 创建加密器（类似mfaService结构体中的encryptor字段）
	key, _ := encrypt.GenerateKey(32) // AES-256
	encryptor, _ := encrypt.NewAESEncryptorWithKey(key)

	// 加密MFA密钥
	secret := "JBSWY3DPEHPK3PXP"
	encryptedSecret, err := encryptor.Encrypt(secret)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Encrypted: %s\n", encryptedSecret)

	// 解密MFA密钥
	decryptedSecret, err := encryptor.Decrypt(encryptedSecret)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted: %s\n", decryptedSecret)

	// Output is omitted because encryption produces different ciphertext each time
	// but decryption will always succeed and produce the original secret
	_ = encryptedSecret
	_ = decryptedSecret
}
