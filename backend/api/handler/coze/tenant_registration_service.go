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

package coze

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/entity"
	"github.com/coze-dev/coze-studio/backend/domain/tenant/service"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

var (
	// 邮箱验证正则
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// 子域名验证正则（3-63字符，仅小写字母、数字、连字符）
	subdomainRegex = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{1,61}[a-z0-9])?$`)

	// 全局服务实例（TODO: 应从依赖注入容器获取）
	validationSvc        *service.TenantValidationService
	registrationSvc      *service.TenantRegistrationService
)

// VerificationCodeRequest 发送验证码请求
type VerificationCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// SendVerificationCode 发送注册验证码
// @router /api/v1/tenants/send-verification-code [POST]
func SendVerificationCode(ctx context.Context, c *app.RequestContext) {
	var req VerificationCodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error": map[string]interface{}{
				"code":        "INVALID_PARAMETER",
				"description": err.Error(),
			},
		})
		return
	}

	// 验证邮箱格式
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "邮箱格式不正确",
			"error": map[string]interface{}{
				"code": "INVALID_EMAIL",
			},
		})
		return
	}

	// TODO: 从服务容器获取验证码服务
	// codeSvc := service.GetVerificationCodeService()
	// err := codeSvc.SendCode(ctx, req.Email)

	// 临时模拟实现
	logs.Infof("Sending verification code to %s", req.Email)
	time.Sleep(100 * time.Millisecond)

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "验证码已发送",
		"data": map[string]interface{}{
			"expires_in": 300,
			"message":    "验证码有效期为5分钟",
		},
		"timestamp": time.Now().UnixMilli(),
	})
}

// CheckCompanyNameAvailabilityRequest 检查企业名称请求
type CheckCompanyNameAvailabilityRequest struct {
	CompanyName string `json:"company_name" query:"company_name" validate:"required,min=2,max=200"`
}

// CheckCompanyNameAvailability 检查企业名称是否可用
// @router /api/v1/tenants/check-company-name [GET]
func CheckCompanyNameAvailability(ctx context.Context, c *app.RequestContext) {
	companyName := c.Query("company_name")
	if companyName == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "企业名称不能为空",
		})
		return
	}

	// TODO: 从依赖注入容器获取验证服务
	// validationSvc := container.GetValidationService()

	// 临时：如果validationSvc未初始化，返回模拟数据
	if validationSvc == nil {
		logs.Warnf("TenantValidationService not initialized, returning mock data")
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"available":   true,
				"suggestions": []string{},
				"note":         "Service not initialized, returning mock data",
			},
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	// 调用服务层检查企业名称
	available, suggestions, err := validationSvc.CheckCompanyNameAvailability(ctx, companyName)
	if err != nil {
		logs.Errorf("Failed to check company name availability: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "检查企业名称失败",
			"error": map[string]interface{}{
				"code":        "CHECK_FAILED",
				"description": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"available":   available,
			"suggestions": suggestions,
		},
		"timestamp": time.Now().UnixMilli(),
	})
}

// CheckSubdomainAvailabilityRequest 检查子域名请求
type CheckSubdomainAvailabilityRequest struct {
	Subdomain string `json:"subdomain" query:"subdomain" validate:"required,min=3,max=63"`
}

// CheckSubdomainAvailability 检查子域名是否可用
// @router /api/v1/tenants/check-subdomain [GET]
func CheckSubdomainAvailability(ctx context.Context, c *app.RequestContext) {
	subdomain := c.Query("subdomain")
	if subdomain == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "子域名不能为空",
		})
		return
	}

	// 验证子域名格式
	if !subdomainRegex.MatchString(subdomain) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "子域名格式不正确，仅支持小写字母、数字、连字符",
			"error": map[string]interface{}{
				"code":        "INVALID_SUBDOMAIN",
				"description": "子域名长度为3-63个字符，仅支持小写字母、数字和连字符",
			},
		})
		return
	}

	// TODO: 从依赖注入容器获取验证服务
	// validationSvc := container.GetValidationService()

	// 临时：如果validationSvc未初始化，返回模拟数据
	if validationSvc == nil {
		logs.Warnf("TenantValidationService not initialized, returning mock data")
		fullDomain := subdomain + ".saas.coze.com"
		c.JSON(http.StatusOK, map[string]interface{}{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"available":   true,
				"suggestions": []string{},
				"full_domain":  fullDomain,
				"note":         "Service not initialized, returning mock data",
			},
			"timestamp": time.Now().UnixMilli(),
		})
		return
	}

	// 调用服务层检查子域名
	available, suggestions, err := validationSvc.CheckSubdomainAvailability(ctx, subdomain)
	if err != nil {
		logs.Errorf("Failed to check subdomain availability: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "检查子域名失败",
			"error": map[string]interface{}{
				"code":        "CHECK_FAILED",
				"description": err.Error(),
			},
		})
		return
	}

	// 构建完整域名
	fullDomain := subdomain + ".saas.coze.com"

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    0,
		"message": "success",
		"data": map[string]interface{}{
			"available":   available,
			"suggestions": suggestions,
			"full_domain":  fullDomain,
		},
		"timestamp": time.Now().UnixMilli(),
	})
}

// RegisterTenantRequest 租户注册请求
type RegisterTenantRequest struct {
	CompanyName      string `json:"company_name" validate:"required,min=2,max=200"`
	Subdomain        string `json:"subdomain" validate:"required,min=3,max=63"`
	Email            string `json:"email" validate:"required,email"`
	VerificationCode string `json:"verification_code" validate:"required,len=6"`
	AdminName        string `json:"admin_name" validate:"required,min=1,max=50"`
	AdminPassword    string `json:"admin_password" validate:"required,min=8,max=100"`
	TenantType       string `json:"tenant_type" validate:"required,oneof=individual team enterprise"`
	BillingCycle     string `json:"billing_cycle" validate:"omitempty,oneof=monthly yearly"`
}

// RegisterTenant 租户注册
// @router /api/v1/tenants/register [POST]
func RegisterTenant(ctx context.Context, c *app.RequestContext) {
	var req RegisterTenantRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误",
			"error": map[string]interface{}{
				"code":        "INVALID_PARAMETER",
				"description": err.Error(),
			},
		})
		return
	}

	// 临时：如果registrationSvc未初始化，返回错误
	if registrationSvc == nil {
		logs.Errorf("TenantRegistrationService not initialized")
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "服务未初始化，请联系管理员",
			"error": map[string]interface{}{
				"code": "SERVICE_NOT_INITIALIZED",
			},
		})
		return
	}

	// 转换请求类型
	tenantType := entity.TenantType(req.TenantType)
	serviceReq := &service.RegisterTenantRequest{
		CompanyName:      req.CompanyName,
		Subdomain:        req.Subdomain,
		Email:            req.Email,
		VerificationCode: req.VerificationCode,
		AdminName:        req.AdminName,
		AdminPassword:    req.AdminPassword,
		TenantType:       tenantType,
	}

	// 处理BillingCycle（可选参数）
	if req.BillingCycle != "" {
		serviceReq.BillingCycle = entity.BillingCycle(req.BillingCycle)
	} else {
		serviceReq.BillingCycle = entity.BillingCycleMonthly // 默认月付
	}

	// 调用服务层注册租户
	result, err := registrationSvc.RegisterTenant(ctx, serviceReq)
	if err != nil {
		logs.Errorf("Failed to register tenant: %v", err)

		// 根据错误类型返回不同的HTTP状态码
		statusCode := http.StatusInternalServerError
		errorCode := "REGISTRATION_FAILED"

		if contains(err.Error(), "验证码") {
			statusCode = http.StatusBadRequest
			errorCode = "INVALID_VERIFICATION_CODE"
		} else if contains(err.Error(), "已被占用") {
			statusCode = http.StatusConflict
			errorCode = "RESOURCE_ALREADY_EXISTS"
		} else if contains(err.Error(), "不能为空") || contains(err.Error(), "长度必须") {
			statusCode = http.StatusBadRequest
			errorCode = "INVALID_PARAMETER"
		}

		c.JSON(statusCode, map[string]interface{}{
			"code":    statusCode,
			"message": err.Error(),
			"error": map[string]interface{}{
				"code": errorCode,
			},
		})
		return
	}

	// 注册成功
	c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    0,
		"message": "租户注册成功",
		"data":    result,
		"timestamp": time.Now().UnixMilli(),
	})
}

// contains 辅助函数：检查字符串是否包含子串（不区分大小写）
func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 ||
		(len(str) > len(substr) && containsIgnoreCase(str, substr)))
}

func containsIgnoreCase(str, substr string) bool {
	// 简单实现，可以优化
	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := str[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 += 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 += 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
