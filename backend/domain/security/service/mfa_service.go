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

package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/security/entity"
	secmodel "github.com/coze-dev/coze-studio/backend/domain/security/internal/dal/model"
	"github.com/coze-dev/coze-studio/backend/pkg/encrypt"
)

// MFAService MFA服务接口
type MFAService interface {
	// EnableMFA 启用MFA
	EnableMFA(ctx context.Context, tenantID, userID string, mfaType entity.MFAType) (*entity.MFAConfig, string, error)

	// DisableMFA 禁用MFA
	DisableMFA(ctx context.Context, tenantID, userID string) error

	// VerifyTOTP 验证TOTP码
	VerifyTOTP(ctx context.Context, tenantID, userID string, code string) (*entity.MFAVerification, error)

	// GenerateRecoveryCodes 生成恢复码
	GenerateRecoveryCodes(ctx context.Context, tenantID, userID string) ([]string, error)

	// UseRecoveryCode 使用恢复码
	UseRecoveryCode(ctx context.Context, tenantID, userID, code string) error

	// GetConfig 获取MFA配置
	GetConfig(ctx context.Context, tenantID, userID string) (*entity.MFAConfig, error)

	// VerifyCode 验证MFA码(通用)
	VerifyCode(ctx context.Context, tenantID, userID string, code string, mfaType entity.MFAType) error

	// LockMFA 锁定MFA
	LockMFA(ctx context.Context, tenantID, userID string, duration time.Duration) error

	// UnlockMFA 解锁MFA
	UnlockMFA(ctx context.Context, tenantID, userID string) error

	// GetVerificationHistory 获取验证历史
	GetVerificationHistory(ctx context.Context, tenantID, userID string, limit int) ([]*entity.MFAVerification, error)
}

// mfaService MFA服务实现
type mfaService struct {
	db             *gorm.DB
	encryptor      encrypt.Encryptor
	maxAttempts    int
	lockDuration   time.Duration
	codeExpiry     time.Duration
	recoveryCount  int
}

// NewMFAService 创建MFA服务实例
func NewMFAService(
	db *gorm.DB,
	encryptor encrypt.Encryptor,
) MFAService {
	return &mfaService{
		db:            db,
		encryptor:     encryptor,
		maxAttempts:   5,           // 最多5次失败尝试
		lockDuration:  30 * time.Minute, // 锁定30分钟
		codeExpiry:    5 * time.Minute,   // 验证码5分钟有效
		recoveryCount: 10,          // 生成10个恢复码
	}
}

// EnableMFA 启用MFA
func (s *mfaService) EnableMFA(ctx context.Context, tenantID, userID string, mfaType entity.MFAType) (*entity.MFAConfig, string, error) {
	// 1. 检查是否已启用
	config, err := s.GetConfig(ctx, tenantID, userID)
	if err == nil && config != nil && config.IsEnabled() {
		return nil, "", fmt.Errorf("MFA already enabled")
	}

	// 2. 生成密钥
	secret := generateSecret()

	// 3. 加密密钥
	encryptedSecret, err := s.encryptor.Encrypt(secret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// 4. 创建配置
	config = &entity.MFAConfig{
		ConfigID:    generateConfigID(),
		TenantID:    tenantID,
		UserID:      userID,
		Type:        mfaType,
		Status:      entity.MFAStatusDisabled, // 先禁用，验证后才启用
		Secret:      encryptedSecret,
		Issuer:      "CozeStudio",
		Account:     userID,
		Algorithm:   "SHA1",
		Digits:      6,
		Period:      30,
		UsedCount:   0,
	}

	// 5. 保存配置
	if err := s.db.WithContext(ctx).Create(s.entityToDO(config)).Error; err != nil {
		return nil, "", fmt.Errorf("failed to save config: %w", err)
	}

	// 6. 生成QR码URL
	qrURL, err := s.generateTOTPURL(config.Issuer, config.Account, secret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate QR URL: %w", err)
	}

	return config, qrURL, nil
}

// DisableMFA 禁用MFA
func (s *mfaService) DisableMFA(ctx context.Context, tenantID, userID string) error {
	return s.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Delete(&secmodel.MFAConfigDO{}).Error
}

// VerifyTOTP 验证TOTP码
func (s *mfaService) VerifyTOTP(ctx context.Context, tenantID, userID string, code string) (*entity.MFAVerification, error) {
	// 1. 获取配置
	config, err := s.GetConfig(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	// 2. 检查锁定状态
	if config.IsLocked() {
		if !config.IsLockExpired() {
			return nil, fmt.Errorf("MFA locked until %v", config.LockedUntil)
		}
		// 锁定已过期，解锁
		s.UnlockMFA(ctx, tenantID, userID)
	}

	// 3. 解密密钥
	secret, err := s.encryptor.Decrypt(config.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// 4. 验证TOTP码
	valid := s.validateTOTP(secret, code, config.Digits, config.Period, time.Now())

	// 5. 创建验证记录
	verification := &entity.MFAVerification{
		VerifyID:   generateVerifyID(),
		TenantID:   tenantID,
		UserID:     userID,
		Type:       config.Type,
		Success:    valid,
	}

	if !valid {
		// 验证失败
		config.FailedAttempts++
		if config.ShouldLock(s.maxAttempts) {
			// 达到最大失败次数，锁定
			lockedUntil := time.Now().Add(s.lockDuration)
			config.LockedUntil = &lockedUntil
			config.Status = entity.MFAStatusLocked

			verification.FailureReason = "Locked due to too many failed attempts"
		} else {
			verification.FailureReason = "Invalid code"
		}

		// 更新配置
		s.db.WithContext(ctx).Model(&secmodel.MFAConfigDO{}).
			Where("config_id = ?", config.ConfigID).
			Updates(map[string]interface{}{
				"failed_attempts": config.FailedAttempts,
				"locked_until":    config.LockedUntil,
				"status":          config.Status,
			})
	} else {
		// 验证成功，启用MFA
		if !config.IsEnabled() {
			config.Status = entity.MFAStatusEnabled
			s.db.WithContext(ctx).Model(&secmodel.MFAConfigDO{}).
				Where("config_id = ?", config.ConfigID).
				Update("status", entity.MFAStatusEnabled)
		}

		// 更新使用统计
		now := time.Now()
		config.UsedCount++
		config.LastUsedAt = &now
		config.FailedAttempts = 0

		s.db.WithContext(ctx).Model(&secmodel.MFAConfigDO{}).
			Where("config_id = ?", config.ConfigID).
			Updates(map[string]interface{}{
				"used_count":   config.UsedCount,
				"last_used_at": config.LastUsedAt,
				"failed_attempts": 0,
			})
	}

	// 6. 保存验证记录
	if err := s.db.WithContext(ctx).Create(s.verificationToDO(verification)).Error; err != nil {
		return nil, err
	}

	if !valid {
		return verification, fmt.Errorf("invalid TOTP code")
	}

	return verification, nil
}

// GenerateRecoveryCodes 生成恢复码
func (s *mfaService) GenerateRecoveryCodes(ctx context.Context, tenantID, userID string) ([]string, error) {
	// 1. 获取配置
	config, err := s.GetConfig(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	// 2. 生成恢复码
	codes := make([]string, s.recoveryCount)
	for i := 0; i < s.recoveryCount; i++ {
		codes[i] = generateRecoveryCode()
	}

	// 3. 保存恢复码哈希
	for _, code := range codes {
		codeHash := hashRecoveryCode(code)
		recoveryCode := &entity.RecoveryCode{
			CodeID:    generateCodeID(),
			ConfigID:  config.ConfigID,
			UserID:    userID,
			TenantID:  tenantID,
			CodeHash:  codeHash,
		}
		if err := s.db.WithContext(ctx).Create(s.recoveryCodeToDO(recoveryCode)).Error; err != nil {
			return nil, err
		}
	}

	return codes, nil
}

// UseRecoveryCode 使用恢复码
func (s *mfaService) UseRecoveryCode(ctx context.Context, tenantID, userID, code string) error {
	// 1. 获取配置
	config, err := s.GetConfig(ctx, tenantID, userID)
	if err != nil {
		return err
	}

	// 2. 查找恢复码
	codeHash := hashRecoveryCode(code)
	var recoveryCodeDO secmodel.RecoveryCodeDO
	err = s.db.WithContext(ctx).
		Where("config_id = ? AND code_hash = ? AND used = ?", config.ConfigID, codeHash, false).
		First(&recoveryCodeDO).Error
	if err != nil {
		return fmt.Errorf("invalid recovery code")
	}

	// 3. 标记为已使用
	now := time.Now()
	err = s.db.WithContext(ctx).
		Model(&secmodel.RecoveryCodeDO{}).
		Where("code_id = ?", recoveryCodeDO.CodeID).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": now,
		}).Error
	if err != nil {
		return err
	}

	return nil
}

// GetConfig 获取MFA配置
func (s *mfaService) GetConfig(ctx context.Context, tenantID, userID string) (*entity.MFAConfig, error) {
	var configDO secmodel.MFAConfigDO
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		First(&configDO).Error
	if err != nil {
		return nil, err
	}
	return s.doToConfigEntity(&configDO), nil
}

// VerifyCode 验证MFA码(通用)
func (s *mfaService) VerifyCode(ctx context.Context, tenantID, userID, code string, mfaType entity.MFAType) error {
	switch mfaType {
	case entity.MFATypeTOTP:
		_, err := s.VerifyTOTP(ctx, tenantID, userID, code)
		return err
	// 可以添加其他类型的验证
	default:
		return fmt.Errorf("unsupported MFA type: %s", mfaType)
	}
}

// LockMFA 锁定MFA
func (s *mfaService) LockMFA(ctx context.Context, tenantID, userID string, duration time.Duration) error {
	lockedUntil := time.Now().Add(duration)
	return s.db.WithContext(ctx).
		Model(&secmodel.MFAConfigDO{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Updates(map[string]interface{}{
			"status":       entity.MFAStatusLocked,
			"locked_until": lockedUntil,
		}).Error
}

// UnlockMFA 解锁MFA
func (s *mfaService) UnlockMFA(ctx context.Context, tenantID, userID string) error {
	return s.db.WithContext(ctx).
		Model(&secmodel.MFAConfigDO{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Updates(map[string]interface{}{
			"status":         entity.MFAStatusEnabled,
			"locked_until":   nil,
			"failed_attempts": 0,
		}).Error
}

// GetVerificationHistory 获取验证历史
func (s *mfaService) GetVerificationHistory(ctx context.Context, tenantID, userID string, limit int) ([]*entity.MFAVerification, error) {
	var verificationDOs []*secmodel.MFAVerificationDO
	err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&verificationDOs).Error
	if err != nil {
		return nil, err
	}

	verifications := make([]*entity.MFAVerification, len(verificationDOs))
	for i, v := range verificationDOs {
		verifications[i] = s.doToVerificationEntity(v)
	}
	return verifications, nil
}

// generateTOTPURL 生成TOTP URL(用于QR码)
func (s *mfaService) generateTOTPURL(issuer, account, secret string) (string, error) {
	v := url.Values{}
	v.Set("secret", secret)
	v.Set("issuer", issuer)

	return fmt.Sprintf("otpauth://totp/%s:%s?%s", issuer, account, v.Encode()), nil
}

// validateTOTP 验证TOTP码
func (s *mfaService) validateTOTP(secret string, code string, digits, period int, timestamp time.Time) bool {
	// 允许时间窗口前后各1个周期
	for i := -1; i <= 1; i++ {
		expectedCode := s.generateTOTP(secret, timestamp.Add(time.Duration(i)*time.Duration(period)*time.Second), digits, period)
		if expectedCode == code {
			return true
		}
	}
	return false
}

// generateTOTP 生成TOTP码
func (s *mfaService) generateTOTP(secret string, timestamp time.Time, digits, period int) string {
	// 解码Base32密钥
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return ""
	}

	// 计算时间计数器
	counter := uint64(timestamp.Unix()) / uint64(period)

	// 将计数器转为字节
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// 使用HMAC-SHA1计算哈希
	hash := hmac.New(sha1.New, key)
	hash.Write(buf)
	h := hash.Sum(nil)

	// 动态截取
	offset := h[len(h)-1] & 0x0f
	truncated := binary.BigEndian.Uint32(h[offset : offset+4])
	truncated &= 0x7fffffff

	// 取模
	mod := uint32(1)
	for i := 0; i < digits; i++ {
		mod *= 10
	}
	code := truncated % mod

	// 格式化为指定位数
	format := fmt.Sprintf("%%0%dd", digits)
	return fmt.Sprintf(format, code)
}

// 生成密钥
func generateSecret() string {
	secret := make([]byte, 20) // 160 bits
	if _, err := rand.Read(secret); err != nil {
		panic(err) // 极少发生
	}
	return base32.StdEncoding.EncodeToString(secret)
}

// 生成恢复码
func generateRecoveryCode() string {
	// 生成8位随机码，格式: XXXX-XXXX
	code := make([]byte, 4)
	if _, err := rand.Read(code); err != nil {
		panic(err)
	}

	part1 := fmt.Sprintf("%04x", binary.BigEndian.Uint32(code[0:4]))
	part2 := make([]byte, 4)
	if _, err := rand.Read(part2); err != nil {
		panic(err)
	}
	part2Str := fmt.Sprintf("%04x", binary.BigEndian.Uint32(part2))

	return strings.ToUpper(part1 + "-" + part2Str)
}

// 哈希恢复码
func hashRecoveryCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return fmt.Sprintf("%x", hash)
}

// 生成ID
func generateConfigID() string {
	return fmt.Sprintf("mfa_cfg_%d", time.Now().UnixNano())
}

func generateCodeID() string {
	return fmt.Sprintf("mfa_rec_%d", time.Now().UnixNano())
}

func generateVerifyID() string {
	return fmt.Sprintf("mfa_ver_%d", time.Now().UnixNano())
}

// Entity <-> DO 转换方法
func (s *mfaService) entityToDO(config *entity.MFAConfig) *secmodel.MFAConfigDO {
	return &secmodel.MFAConfigDO{
		ConfigID:       config.ConfigID,
		TenantID:       config.TenantID,
		UserID:         config.UserID,
		Type:           string(config.Type),
		Status:         string(config.Status),
		Secret:         config.Secret,
		BackupCodes:    config.BackupCodes,
		Issuer:         config.Issuer,
		Account:        config.Account,
		Algorithm:      config.Algorithm,
		Digits:         config.Digits,
		Period:         config.Period,
		BackupType:     string(config.BackupType),
		BackupValue:    config.BackupValue,
		UsedCount:      config.UsedCount,
		LastUsedAt:     config.LastUsedAt,
		FailedAttempts: config.FailedAttempts,
		LockedUntil:    config.LockedUntil,
	}
}

func (s *mfaService) doToConfigEntity(do *secmodel.MFAConfigDO) *entity.MFAConfig {
	return &entity.MFAConfig{
		ConfigID:       do.ConfigID,
		TenantID:       do.TenantID,
		UserID:         do.UserID,
		Type:           entity.MFAType(do.Type),
		Status:         entity.MFAStatus(do.Status),
		Secret:         do.Secret,
		BackupCodes:    do.BackupCodes,
		Issuer:         do.Issuer,
		Account:        do.Account,
		Algorithm:      do.Algorithm,
		Digits:         do.Digits,
		Period:         do.Period,
		BackupType:     entity.MFAType(do.BackupType),
		BackupValue:    do.BackupValue,
		UsedCount:      do.UsedCount,
		LastUsedAt:     do.LastUsedAt,
		FailedAttempts: do.FailedAttempts,
		LockedUntil:    do.LockedUntil,
	}
}

func (s *mfaService) recoveryCodeToDO(code *entity.RecoveryCode) *secmodel.RecoveryCodeDO {
	return &secmodel.RecoveryCodeDO{
		CodeID:    code.CodeID,
		ConfigID:  code.ConfigID,
		UserID:    code.UserID,
		TenantID:  code.TenantID,
		CodeHash:  code.CodeHash,
		Used:      code.Used,
		UsedAt:    code.UsedAt,
		ExpiresAt: code.ExpiresAt,
	}
}

func (s *mfaService) verificationToDO(v *entity.MFAVerification) *secmodel.MFAVerificationDO {
	return &secmodel.MFAVerificationDO{
		VerifyID:      v.VerifyID,
		TenantID:      v.TenantID,
		UserID:        v.UserID,
		Type:          string(v.Type),
		Success:       v.Success,
		FailureReason: v.FailureReason,
		RequestIP:     v.RequestIP,
		UserAgent:     v.UserAgent,
	}
}

func (s *mfaService) doToVerificationEntity(do *secmodel.MFAVerificationDO) *entity.MFAVerification {
	return &entity.MFAVerification{
		VerifyID:      do.VerifyID,
		TenantID:      do.TenantID,
		UserID:        do.UserID,
		Type:          entity.MFAType(do.Type),
		Success:       do.Success,
		FailureReason: do.FailureReason,
		RequestIP:     do.RequestIP,
		UserAgent:     do.UserAgent,
	}
}
