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
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

const (
	// 验证码配置
	codeLength          = 6                // 验证码长度（6位数字）
	codeExpiry          = 5 * time.Minute  // 验证码有效期（5分钟）
	resendLimitDuration = 1 * time.Minute // 重新发送限制时长（1分钟）
	resendMaxTimes      = 5                // 同一邮箱最大发送次数

	// Redis键前缀
	redisKeyVerifyCode   = "verify_code:%s"     // 验证码存储
	redisKeyResendCount  = "verify_count:%s"    // 发送次数计数器
	redisKeyResendLimit  = "verify_limit:%s"    // 发送频率限制
)

// VerificationCodeService 验证码服务
// 提供验证码生成、发送、验证功能，支持Redis存储和频率限制
//
// 使用示例:
//
//	codeSvc := service.NewVerificationCodeService(redisClient, emailService)
//	err := codeSvc.SendCode(ctx, "user@example.com")
//	valid, err := codeSvc.VerifyCode(ctx, "user@example.com", "123456")
type VerificationCodeService struct {
	redisClient *redis.Client
	emailSvc    *EmailService
}

// NewVerificationCodeService 创建验证码服务实例
func NewVerificationCodeService(redisClient *redis.Client, emailSvc *EmailService) *VerificationCodeService {
	return &VerificationCodeService{
		redisClient: redisClient,
		emailSvc:    emailSvc,
	}
}

// SendCode 发送验证码
//
// 功能：
// 1. 检查发送频率限制（1分钟内只能发送1次）
// 2. 生成6位随机数字验证码
// 3. 存储到Redis，有效期5分钟
// 4. 发送验证码邮件
// 5. 更新发送次数计数器
//
// 参数:
//   - ctx: 上下文
//   - email: 收件人邮箱地址
//
// 返回:
//   - error: 发送失败或频率超限时返回错误
func (s *VerificationCodeService) SendCode(ctx context.Context, email string) error {
	// 1. 检查发送频率限制
	limitKey := fmt.Sprintf(redisKeyResendLimit, email)
	allowed, err := s.checkRateLimit(ctx, limitKey)
	if err != nil {
		return errorx.Wrap(err, errno.ErrVerificationCodeCheckFailed)
	}
	if !allowed {
		return errorx.NewByErrorCode(errno.ErrVerificationCodeRateLimit)
	}

	// 2. 生成6位随机数字验证码
	code := s.generateCode()

	// 3. 存储到Redis（5分钟过期）
	codeKey := fmt.Sprintf(redisKeyVerifyCode, email)
	if err := s.redisClient.Set(ctx, codeKey, code, codeExpiry).Err(); err != nil {
		return errorx.Wrap(err, errno.ErrVerificationCodeStoreFailed)
	}

	// 4. 发送验证码邮件
	if err := s.emailSvc.SendVerificationCode(ctx, email, code); err != nil {
		// 发送失败，删除已存储的验证码
		_ = s.redisClient.Del(ctx, codeKey)
		return errorx.Wrap(err, errno.ErrVerificationCodeSendFailed)
	}

	// 5. 更新发送次数计数器
	countKey := fmt.Sprintf(redisKeyResendCount, email)
	if err := s.incrementSendCount(ctx, countKey); err != nil {
		logs.Warnf("Failed to increment send count: %v", err)
	}

	logs.Infof("Verification code sent successfully to %s", email)
	return nil
}

// VerifyCode 验证验证码
//
// 功能：
// 1. 从Redis获取存储的验证码
// 2. 比对验证码是否匹配
// 3. 验证成功后删除验证码（一次性使用）
//
// 参数:
//   - ctx: 上下文
//   - email: 邮箱地址
//   - code: 用户输入的验证码
//
// 返回:
//   - bool: 验证是否成功
//   - error: Redis操作失败时返回错误
func (s *VerificationCodeService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	codeKey := fmt.Sprintf(redisKeyVerifyCode, email)

	// 1. 获取存储的验证码
	storedCode, err := s.redisClient.Get(ctx, codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			// 验证码不存在或已过期
			return false, nil
		}
		return false, errorx.Wrap(err, errno.ErrVerificationCodeGetFailed)
	}

	// 2. 验证码比对
	if storedCode != code {
		logs.Warnf("Invalid verification code for %s: expected=%s, got=%s", email, storedCode, code)
		return false, nil
	}

	// 3. 验证成功，删除验证码（一次性使用）
	if err := s.redisClient.Del(ctx, codeKey).Err(); err != nil {
		logs.Warnf("Failed to delete verified code: %v", err)
	}

	logs.Infof("Verification code verified successfully for %s", email)
	return true, nil
}

// generateCode 生成6位随机数字验证码
func (s *VerificationCodeService) generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min, max := 100000, 999999
	code := r.Intn(max-min+1) + min // 生成100000-999999之间的随机数
	return fmt.Sprintf("%06d", code)
}

// checkRateLimit 检查发送频率限制
// 返回 true 表示允许发送，false 表示超出频率限制
func (s *VerificationCodeService) checkRateLimit(ctx context.Context, limitKey string) (bool, error) {
	// 使用Redis INCR + EXPIRE 实现滑动窗口限流
	result, err := s.redisClient.Incr(ctx, limitKey).Result()
	if err != nil {
		return false, err
	}

	// 首次设置，添加过期时间
	if result == 1 {
		if err := s.redisClient.Expire(ctx, limitKey, resendLimitDuration).Err(); err != nil {
			logs.Warnf("Failed to set expiry for rate limit key: %v", err)
		}
	}

	// 检查是否超过限制
	return result <= 1, nil
}

// incrementSendCount 增加发送次数计数器
func (s *VerificationCodeService) incrementSendCount(ctx context.Context, countKey string) error {
	// 计数器每天重置
	ttl := 24 * time.Hour

	_, err := s.redisClient.Incr(ctx, countKey).Result()
	if err != nil {
		return err
	}

	// 首次设置，添加过期时间
	result, _ := s.redisClient.Get(ctx, countKey).Int64()
	if result == 1 {
		if err := s.redisClient.Expire(ctx, countKey, ttl).Err(); err != nil {
			logs.Warnf("Failed to set expiry for count key: %v", err)
		}
	}

	return nil
}

// GetSendCount 获取发送次数（用于测试和监控）
func (s *VerificationCodeService) GetSendCount(ctx context.Context, email string) (int64, error) {
	countKey := fmt.Sprintf(redisKeyResendCount, email)
	count, err := s.redisClient.Get(ctx, countKey).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
