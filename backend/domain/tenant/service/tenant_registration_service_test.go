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

	"github.com/stretchr/testify/assert"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
)

// TestValidateRequest 测试请求验证逻辑
func TestValidateRequest(t *testing.T) {
	// 创建服务实例（DB可为nil，因为validateRequest不使用DB）
	service := NewTenantRegistrationService(nil, nil, nil, nil, nil, nil)
	ctx := context.Background()

	t.Run("参数验证失败-请求为空", func(t *testing.T) {
		err := service.validateRequest(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "请求不能为空")
	})

	t.Run("参数验证失败-企业名称为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "企业名称不能为空")
	})

	t.Run("参数验证失败-企业名称太短", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "a",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "企业名称长度必须在2-200个字符之间")
	})

	t.Run("参数验证失败-企业名称太长", func(t *testing.T) {
		longName := string(make([]byte, 201))
		req := &RegisterTenantRequest{
			CompanyName:      longName,
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "企业名称长度必须在2-200个字符之间")
	})

	t.Run("参数验证失败-子域名为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "子域名不能为空")
	})

	t.Run("参数验证失败-子域名太短", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "ab",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "子域名长度必须在3-63个字符之间")
	})

	t.Run("参数验证失败-子域名太长", func(t *testing.T) {
		longSubdomain := string(make([]byte, 64))
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        longSubdomain,
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "子域名长度必须在3-63个字符之间")
	})

	t.Run("参数验证失败-邮箱为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "邮箱不能为空")
	})

	t.Run("参数验证失败-验证码为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "验证码不能为空")
	})

	t.Run("参数验证失败-验证码长度不正确", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "12345",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "验证码格式不正确")
	})

	t.Run("参数验证失败-管理员姓名为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "管理员姓名不能为空")
	})

	t.Run("参数验证失败-管理员密码为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "管理员密码不能为空")
	})

	t.Run("参数验证失败-管理员密码太短", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "pass",
			TenantType:       entity.TenantTypeIndividual,
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "管理员密码长度不能少于8位")
	})

	t.Run("参数验证失败-租户类型为空", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       "",
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "租户类型不能为空")
	})

	t.Run("参数验证失败-无效的租户类型", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantType("invalid"),
		}
		err := service.validateRequest(ctx, req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "无效的租户类型")
	})

	t.Run("参数验证成功-个人租户", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试公司",
			Subdomain:        "testcompany",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeIndividual,
			BillingCycle:     entity.BillingCycleMonthly,
		}
		err := service.validateRequest(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("参数验证成功-团队租户", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试团队",
			Subdomain:        "testteam",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeTeam,
			BillingCycle:     entity.BillingCycleYearly,
		}
		err := service.validateRequest(ctx, req)
		assert.NoError(t, err)
	})

	t.Run("参数验证成功-企业租户", func(t *testing.T) {
		req := &RegisterTenantRequest{
			CompanyName:      "测试企业",
			Subdomain:        "testenterprise",
			Email:            "user@example.com",
			VerificationCode: "123456",
			AdminName:        "管理员",
			AdminPassword:    "password123",
			TenantType:       entity.TenantTypeEnterprise,
			BillingCycle:     entity.BillingCycleYearly,
		}
		err := service.validateRequest(ctx, req)
		assert.NoError(t, err)
	})
}

// TestCreateInitialQuotas 测试初始配额创建逻辑
func TestCreateInitialQuotas(t *testing.T) {
	service := NewTenantRegistrationService(nil, nil, nil, nil, nil, nil)

	t.Run("创建初始配额", func(t *testing.T) {
		tenantID := "test_tenant_id"
		quotas := service.createInitialQuotas(tenantID)

		assert.Equal(t, 4, len(quotas), "应该创建4个配额")

		// 验证每个配额
		resourceTypes := make(map[entity.ResourceType]*entity.Quota)
		for i := range quotas {
			quota := &quotas[i]
			assert.Equal(t, tenantID, quota.TenantID)
			assert.NotEmpty(t, quota.QuotaID)
			assert.Equal(t, 0, quota.UsedCount, "初始使用量应该为0")
			assert.Greater(t, quota.CreatedAt, int64(0))
			assert.Greater(t, quota.UpdatedAt, int64(0))
			resourceTypes[quota.ResourceType] = quota
		}

		// 验证Bot配额
		botQuota := resourceTypes[entity.ResourceTypeBots]
		assert.NotNil(t, botQuota)
		assert.Equal(t, freeTierBotsLimit, botQuota.MaxLimit)
		assert.Equal(t, entity.ResetCycleNever, botQuota.ResetCycle)

		// 验证消息配额
		msgQuota := resourceTypes[entity.ResourceTypeMessages]
		assert.NotNil(t, msgQuota)
		assert.Equal(t, freeTierMessagesLimit, msgQuota.MaxLimit)
		assert.Equal(t, entity.ResetCycleMonthly, msgQuota.ResetCycle)

		// 验证存储配额
		storageQuota := resourceTypes[entity.ResourceTypeStorage]
		assert.NotNil(t, storageQuota)
		assert.Equal(t, freeTierStorageLimit, storageQuota.MaxLimit)
		assert.Equal(t, entity.ResetCycleNever, storageQuota.ResetCycle)

		// 验证团队成员配额
		memberQuota := resourceTypes[entity.ResourceTypeTeamMembers]
		assert.NotNil(t, memberQuota)
		assert.Equal(t, freeTierTeamMembersLimit, memberQuota.MaxLimit)
		assert.Equal(t, entity.ResetCycleNever, memberQuota.ResetCycle)
	})

	t.Run("验证配额ID唯一性", func(t *testing.T) {
		tenantID := "test_tenant_id"
		quotas := service.createInitialQuotas(tenantID)

		// 验证所有配额ID不同
		ids := make(map[string]bool)
		for _, quota := range quotas {
			assert.False(t, ids[quota.QuotaID], "配额ID应该唯一: %s", quota.QuotaID)
			ids[quota.QuotaID] = true
		}
		assert.Equal(t, len(quotas), len(ids), "所有配额ID应该唯一")
	})

	t.Run("验证配额时间戳", func(t *testing.T) {
		tenantID := "test_tenant_id"
		quotas := service.createInitialQuotas(tenantID)

		for _, quota := range quotas {
			assert.Greater(t, quota.CreatedAt, int64(0), "创建时间应该大于0")
			assert.Greater(t, quota.UpdatedAt, int64(0), "更新时间应该大于0")
			assert.Equal(t, quota.CreatedAt, quota.UpdatedAt, "创建和更新时间应该相同")
			assert.Equal(t, quota.CreatedAt, quota.LastResetAt, "重置时间应该与创建时间相同")
		}
	})
}

// TestGenerateID 测试ID生成逻辑
func TestGenerateID(t *testing.T) {
	t.Run("生成ID格式验证", func(t *testing.T) {
		id1 := generateID("test")
		id2 := generateID("test")

		assert.NotEmpty(t, id1)
		assert.NotEmpty(t, id2)
		assert.NotEqual(t, id1, id2, "相同前缀的ID应该不同")
		assert.Contains(t, id1, "test_")
		assert.Contains(t, id2, "test_")
	})

	t.Run("生成不同前缀的ID", func(t *testing.T) {
		tenantID := generateTenantID()
		subID := generateSubscriptionID()
		quotaID := generateQuotaID()
		userID := generateUserID()

		assert.Contains(t, tenantID, "tenant_")
		assert.Contains(t, subID, "sub_")
		assert.Contains(t, quotaID, "quota_")
		assert.Contains(t, userID, "user_")
	})

	t.Run("ID唯一性测试", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := generateID("test")
			assert.False(t, ids[id], "生成的ID应该唯一: %s", id)
			ids[id] = true
		}
		assert.Equal(t, 100, len(ids), "所有生成的ID应该唯一")
	})
}
