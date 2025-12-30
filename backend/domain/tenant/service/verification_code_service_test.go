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
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

// setupTestRedis 设置测试用Redis
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return mr, client
}

func TestVerificationCodeService_SendCode(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	emailSvc := NewEmailService(false, "noreply@coze.com", "Coze Studio")
	codeSvc := NewVerificationCodeService(redisClient, emailSvc)

	ctx := context.Background()
	testEmail := "test@example.com"

	t.Run("成功发送验证码", func(t *testing.T) {
		err := codeSvc.SendCode(ctx, testEmail)
		assert.NoError(t, err)

		// 验证Redis中存储了验证码
		codeKey := "verify_code:" + testEmail
		storedCode, err := redisClient.Get(ctx, codeKey).Result()
		assert.NoError(t, err)
		assert.Len(t, storedCode, 6)
		assert.Regexp(t, `^\d{6}$`, storedCode)
	})

	t.Run("发送频率限制", func(t *testing.T) {
		// 第一次发送成功
		err := codeSvc.SendCode(ctx, testEmail)
		assert.NoError(t, err)

		// 1分钟内第二次发送应被限制
		err = codeSvc.SendCode(ctx, testEmail)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "发送过于频繁")
	})

	t.Run("等待1分钟后可重新发送", func(t *testing.T) {
		// 快进时间61秒
		mr.FastForward(61 * time.Second)

		err := codeSvc.SendCode(ctx, testEmail)
		assert.NoError(t, err)
	})
}

func TestVerificationCodeService_VerifyCode(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	emailSvc := NewEmailService(false, "noreply@coze.com", "Coze Studio")
	codeSvc := NewVerificationCodeService(redisClient, emailSvc)

	ctx := context.Background()
	testEmail := "test@example.com"

	t.Run("验证成功", func(t *testing.T) {
		// 1. 发送验证码
		err := codeSvc.SendCode(ctx, testEmail)
		require.NoError(t, err)

		// 2. 获取发送的验证码
		codeKey := "verify_code:" + testEmail
		storedCode, err := redisClient.Get(ctx, codeKey).Result()
		require.NoError(t, err)

		// 3. 验证码验证
		valid, err := codeSvc.VerifyCode(ctx, testEmail, storedCode)
		assert.NoError(t, err)
		assert.True(t, valid)

		// 4. 验证码应被删除（一次性使用）
		_, err = redisClient.Get(ctx, codeKey).Result()
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("验证码错误", func(t *testing.T) {
		// 1. 发送验证码
		err := codeSvc.SendCode(ctx, testEmail)
		require.NoError(t, err)

		// 2. 使用错误的验证码
		valid, err := codeSvc.VerifyCode(ctx, testEmail, "000000")
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("验证码不存在", func(t *testing.T) {
		valid, err := codeSvc.VerifyCode(ctx, "nonexistent@example.com", "123456")
		assert.NoError(t, err)
		assert.False(t, valid)
	})

	t.Run("验证码过期", func(t *testing.T) {
		// 1. 发送验证码
		err := codeSvc.SendCode(ctx, testEmail)
		require.NoError(t, err)

		// 2. 快进时间6分钟（超过5分钟有效期）
		mr.FastForward(6 * time.Minute)

		// 3. 验证码应已过期
		valid, err := codeSvc.VerifyCode(ctx, testEmail, "123456")
		assert.NoError(t, err)
		assert.False(t, valid)
	})
}

func TestVerificationCodeService_generateCode(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	emailSvc := NewEmailService(false, "noreply@coze.com", "Coze Studio")
	codeSvc := NewVerificationCodeService(redisClient, emailSvc)

	t.Run("生成6位数字验证码", func(t *testing.T) {
		code := codeSvc.generateCode()
		assert.Len(t, code, 6)
		assert.Regexp(t, `^\d{6}$`, code)
	})

	t.Run("验证码唯一性测试", func(t *testing.T) {
		codes := make(map[string]bool)
		for i := 0; i < 100; i++ {
			code := codeSvc.generateCode()
			codes[code] = true
		}
		// 100次生成应该有重复（随机性）
		assert.True(t, len(codes) < 100)
	})
}

func TestVerificationCodeService_GetSendCount(t *testing.T) {
	mr, redisClient := setupTestRedis(t)
	defer mr.Close()

	emailSvc := NewEmailService(false, "noreply@coze.com", "Coze Studio")
	codeSvc := NewVerificationCodeService(redisClient, emailSvc)

	ctx := context.Background()
	testEmail := "test@example.com"

	t.Run("初始发送次数为0", func(t *testing.T) {
		count, err := codeSvc.GetSendCount(ctx, testEmail)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("发送后计数增加", func(t *testing.T) {
		err := codeSvc.SendCode(ctx, testEmail)
		require.NoError(t, err)

		count, err := codeSvc.GetSendCount(ctx, testEmail)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 第二次发送
		mr.FastForward(61 * time.Second)
		err = codeSvc.SendCode(ctx, testEmail)
		require.NoError(t, err)

		count, err = codeSvc.GetSendCount(ctx, testEmail)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}
