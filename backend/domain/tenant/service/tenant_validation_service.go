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
	"regexp"
	"strings"
	"unicode"

	"github.com/coze-dev/coze-studio/backend/domain/tenant/repository"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// TenantValidationService 租户验证服务
// 提供企业名称、子域名唯一性检查和智能建议功能
type TenantValidationService struct {
	tenantRepo repository.TenantRepository
}

// NewTenantValidationService 创建租户验证服务实例
func NewTenantValidationService(tenantRepo repository.TenantRepository) *TenantValidationService {
	return &TenantValidationService{
		tenantRepo: tenantRepo,
	}
}

// CheckCompanyNameAvailability 检查企业名称是否可用
//
// 功能：
// 1. 去除企业名称首尾空格和特殊字符
// 2. 检查名称长度（2-200字符）
// 3. 查询数据库判断是否已存在
// 4. 如果已存在，生成智能建议名称
//
// 参数:
//   - ctx: 上下文
//   - companyName: 企业名称
//
// 返回:
//   - available: 是否可用
//   - suggestions: 不可用时的建议名称列表
//   - error: 查询失败时返回错误
func (s *TenantValidationService) CheckCompanyNameAvailability(ctx context.Context, companyName string) (available bool, suggestions []string, err error) {
	// 1. 清理企业名称
	cleanName := s.cleanCompanyName(companyName)

	// 2. 验证名称长度
	if len(cleanName) < 2 {
		return false, []string{}, errorx.NewByErrorCode(errno.ErrTenantNameTooShort)
	}
	if len(cleanName) > 200 {
		return false, []string{}, errorx.NewByErrorCode(errno.ErrTenantNameTooLong)
	}

	// 3. 检查名称是否已被使用
	existingTenant, err := s.tenantRepo.GetByName(ctx, cleanName)
	if err != nil && err != repository.ErrRecordNotFound {
		logs.Errorf("Failed to check company name availability: %v", err)
		return false, nil, errorx.Wrap(err, errno.ErrTenantCheckNameFailed)
	}

	// 4. 名称可用（记录不存在）
	if existingTenant == nil {
		return true, []string{}, nil
	}

	// 5. 名称已被占用，生成建议
	suggestions = s.generateCompanyNameSuggestions(cleanName)
	logs.Infof("Company name '%s' already taken, generated %d suggestions", cleanName, len(suggestions))

	return false, suggestions, nil
}

// CheckSubdomainAvailability 检查子域名是否可用
//
// 功能：
// 1. 验证子域名格式（3-63字符，小写字母数字连字符）
// 2. 转换为小写
// 3. 查询数据库判断是否已存在
// 4. 如果已存在，生成智能建议子域名
//
// 参数:
//   - ctx: 上下文
//   - subdomain: 子域名
//
// 返回:
//   - available: 是否可用
//   - suggestions: 不可用时的建议子域名列表
//   - error: 验证或查询失败时返回错误
func (s *TenantValidationService) CheckSubdomainAvailability(ctx context.Context, subdomain string) (available bool, suggestions []string, err error) {
	// 1. 转换为小写并清理
	subdomain = strings.ToLower(strings.TrimSpace(subdomain))

	// 2. 移除特殊字符，只保留字母、数字、连字符
	cleanSubdomain := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, subdomain)

	// 3. 验证格式
	if len(cleanSubdomain) < 3 {
		return false, []string{}, errorx.NewByErrorCode(errno.ErrTenantSubdomainTooShort)
	}
	if len(cleanSubdomain) > 63 {
		return false, []string{}, errorx.NewByErrorCode(errno.ErrTenantSubdomainTooLong)
	}

	// 不能以连字符开头或结尾
	if strings.HasPrefix(cleanSubdomain, "-") || strings.HasSuffix(cleanSubdomain, "-") {
		return false, []string{}, errorx.NewByErrorCode(errno.ErrTenantSubdomainInvalidFormat)
	}

	// 4. 检查是否已被使用
	existingTenant, err := s.tenantRepo.GetBySubdomain(ctx, cleanSubdomain)
	if err != nil && err != repository.ErrRecordNotFound {
		logs.Errorf("Failed to check subdomain availability: %v", err)
		return false, nil, errorx.Wrap(err, errno.ErrTenantCheckSubdomainFailed)
	}

	// 5. 子域名可用（记录不存在）
	if existingTenant == nil {
		return true, []string{}, nil
	}

	// 6. 子域名已被占用，生成建议
	suggestions = s.generateSubdomainSuggestions(cleanSubdomain)
	logs.Infof("Subdomain '%s' already taken, generated %d suggestions", cleanSubdomain, len(suggestions))

	return false, suggestions, nil
}

// cleanCompanyName 清理企业名称
// 去除首尾空格、特殊字符，保留中英文、数字、常用符号
func (s *TenantValidationService) cleanCompanyName(name string) string {
	// 去除首尾空格
	name = strings.TrimSpace(name)

	// 去除多余空格
	reg := regexp.MustCompile(`\s+`)
	name = reg.ReplaceAllString(name, " ")

	// 去除不可见字符（除了空格）
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) || r == ' ' {
			return r
		}
		return -1
	}, name)

	return cleaned
}

// generateCompanyNameSuggestions 生成企业名称建议
// 基于原名称生成3-5个可用建议
func (s *TenantValidationService) generateCompanyNameSuggestions(companyName string) []string {
	suggestions := []string{}

	// 策略1: 添加后缀
	suffixes := []string{"（科技）", "（信息科技）", "（智能）", "（数字）", "（云）", "-科技", "-智能", "-2025"}
	for _, suffix := range suffixes {
		suggestions = append(suggestions, companyName+suffix)
		if len(suggestions) >= 5 {
			break
		}
	}

	// 策略2: 如果名称中有括号，尝试修改
	if strings.Contains(companyName, "(") && strings.Contains(companyName, ")") {
		// 将括号内容改为其他描述
		baseName := companyName[:strings.Index(companyName, "(")]
		descriptions := []string{"（有限合伙）", "（集团）", "（控股）"}
		for _, desc := range descriptions {
			suggestions = append(suggestions, baseName+desc)
			if len(suggestions) >= 8 {
				break
			}
		}
	}

	// 限制返回数量
	if len(suggestions) > 8 {
		suggestions = suggestions[:8]
	}

	return suggestions
}

// generateSubdomainSuggestions 生成子域名建议
// 基于原子域名生成5-8个可用建议
func (s *TenantValidationService) generateSubdomainSuggestions(subdomain string) []string {
	suggestions := []string{}
	base := strings.TrimSuffix(subdomain, "-")

	// 策略1: 添加数字后缀
	for i := 1; i <= 3; i++ {
		suggestions = append(suggestions, fmt.Sprintf("%s%d", base, i))
	}

	// 策略2: 添加年份后缀
	suggestions = append(suggestions, base+"-2025", base+"-2024", base+"-new")

	// 策略3: 添加行业后缀
	industrySuffixes := []string{"-tech", "-cloud", "-ai", "-data", "-sys"}
	for _, suffix := range industrySuffixes {
		suggestions = append(suggestions, base+suffix)
		if len(suggestions) >= 8 {
			break
		}
	}

	// 限制返回数量
	if len(suggestions) > 8 {
		suggestions = suggestions[:8]
	}

	return suggestions
}
