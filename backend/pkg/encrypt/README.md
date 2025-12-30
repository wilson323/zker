# encrypt - AES加密解密包

## 概述

`encrypt`包提供基于AES-GCM（Galois/Counter Mode）的加密解密功能。AES-GCM是一种认证加密模式，同时提供机密性和完整性保护。

## 特性

- ✅ 支持AES-128、AES-192、AES-256
- ✅ 使用GCM模式（认证加密）
- ✅ 随机nonce生成
- ✅ Base64编码输出
- ✅ 实现`Encryptor`接口
- ✅ 完整的单元测试（86%+覆盖率）

## 安装

```bash
import "github.com/coze-dev/coze-studio/backend/pkg/encrypt"
```

## 使用方法

### 1. 基本使用（直接调用函数）

```go
package main

import (
    "fmt"
    "github.com/coze-dev/coze-studio/backend/pkg/encrypt"
)

func main() {
    // 生成密钥（32字节 = AES-256）
    key, err := encrypt.GenerateKey(32)
    if err != nil {
        panic(err)
    }

    // 加密
    plaintext := "Hello, World!"
    ciphertext, err := encrypt.AESEncrypt(key, plaintext)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Encrypted: %s\n", ciphertext)

    // 解密
    decrypted, err := encrypt.AESDecrypt(key, ciphertext)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Decrypted: %s\n", decrypted)
}
```

### 2. 使用Encryptor接口（推荐用于MFA等服务）

```go
package main

import (
    "fmt"
    "github.com/coze-dev/coze-studio/backend/pkg/encrypt"
)

func main() {
    // 创建加密器
    key, _ := encrypt.GenerateKey(32)
    encryptor, _ := encrypt.NewAESEncryptorWithKey(key)

    // 使用Encryptor接口
    secret := "JBSWY3DPEHPK3PXP" // MFA密钥

    // 加密
    encryptedSecret, err := encryptor.Encrypt(secret)
    if err != nil {
        panic(err)
    }

    // 解密
    decryptedSecret, err := encryptor.Decrypt(encryptedSecret)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Original: %s\n", secret)
    fmt.Printf("Decrypted: %s\n", decryptedSecret)
}
```

### 3. 在服务中使用（示例：MFA服务）

```go
package service

import (
    "github.com/coze-dev/coze-studio/backend/pkg/encrypt"
)

type mfaService struct {
    encryptor encrypt.Encryptor
}

func NewMFAService(secretKey []byte) (MFAService, error) {
    encryptor, err := encrypt.NewAESEncryptorWithKey(secretKey)
    if err != nil {
        return nil, err
    }

    return &mfaService{
        encryptor: encryptor,
    }, nil
}

func (s *mfaService) EnableMFA(ctx context.Context, userID string) error {
    // 生成MFA密钥
    secret := "JBSWY3DPEHPK3PXP"

    // 加密存储
    encryptedSecret, err := s.encryptor.Encrypt(secret)
    if err != nil {
        return fmt.Errorf("failed to encrypt secret: %w", err)
    }

    // 保存到数据库...
    return nil
}

func (s *mfaService) VerifyMFA(ctx context.Context, userID, code string) error {
    // 从数据库读取加密的密钥
    encryptedSecret := loadFromDatabase(userID)

    // 解密
    secret, err := s.encryptor.Decrypt(encryptedSecret)
    if err != nil {
        return fmt.Errorf("failed to decrypt secret: %w", err)
    }

    // 验证MFA码...
    return nil
}
```

## API文档

### 接口

```go
type Encryptor interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}
```

### 函数

#### GenerateKey

```go
func GenerateKey(size int) ([]byte, error)
```

生成随机密钥。

**参数**:
- `size`: 密钥大小，必须是16（AES-128）、24（AES-192）或32（AES-256）字节

**返回**:
- `[]byte`: 随机生成的密钥
- `error`: 错误信息

**示例**:
```go
key, err := encrypt.GenerateKey(32) // AES-256
```

#### NewAESEncryptor

```go
func NewAESEncryptor(keySize int) (Encryptor, error)
```

创建AES加密器（自动生成密钥）。

**参数**:
- `keySize`: 密钥大小（16/24/32字节）

**返回**:
- `Encryptor`: 加密器实例
- `error`: 错误信息

#### NewAESEncryptorWithKey

```go
func NewAESEncryptorWithKey(key []byte) (Encryptor, error)
```

使用指定密钥创建AES加密器。

**参数**:
- `key`: 密钥（16/24/32字节）

**返回**:
- `Encryptor`: 加密器实例
- `error`: 错误信息

#### AESEncrypt

```go
func AESEncrypt(key []byte, plaintext string) (string, error)
```

AES加密明文。

**参数**:
- `key`: 加密密钥
- `plaintext`: 明文

**返回**:
- `string`: Base64编码的密文
- `error`: 错误信息

#### AESDecrypt

```go
func AESDecrypt(key []byte, ciphertext string) (string, error)
```

AES解密密文。

**参数**:
- `key`: 解密密钥
- `ciphertext`: Base64编码的密文

**返回**:
- `string`: 明文
- `error`: 错误信息

## 安全说明

1. **密钥管理**: 密钥应该安全存储，推荐使用环境变量或密钥管理服务
2. **密钥长度**: 建议使用AES-256（32字节）以获得最佳安全性
3. **随机性**: 使用加密安全的随机数生成器生成密钥
4. **密钥轮换**: 定期轮换密钥以提高安全性

## 性能

```
BenchmarkAESEncrypt-8     1000000    1032 ns/op
BenchmarkAESDecrypt-8     1000000    1028 ns/op
```

## 测试

```bash
# 运行所有测试
go test ./pkg/encrypt/...

# 运行测试并查看覆盖率
go test ./pkg/encrypt/... -cover

# 运行基准测试
go test ./pkg/encrypt/... -bench=.
```

## 技术细节

### AES-GCM模式

本包使用AES-GCM（Galois/Counter Mode）模式，具有以下优势：

1. **认证加密**: 同时提供机密性和完整性保护
2. **性能优异**: 并行加密，速度快
3. **标准算法**: 广泛支持，NIST标准

### 密文格式

```
[nonce (12 bytes)] [ciphertext] [authentication tag (16 bytes)]
```

- Nonce: 每次加密随机生成，确保相同明文产生不同密文
- Authentication Tag: 验证数据完整性和真实性

## 依赖

- Go标准库 `crypto/aes`
- Go标准库 `crypto/cipher`
- Go标准库 `crypto/rand`
- Go标准库 `encoding/base64`

## 许可证

Apache License 2.0

## 作者

coze-dev
