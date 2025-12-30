// Package saga 提供租户注册相关的Saga实现
package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TenantRegistrationCommand 租户注册命令
type TenantRegistrationCommand struct {
	EnterpriseName string
	Industry       string
	Scale          string
	Subdomain      string
	ContactName    string
	ContactEmail   string
	ContactPhone   string
}

// Tenant 租户实体(简化版)
type Tenant struct {
	ID            string
	Name          string
	Industry      string
	Scale         string
	Subdomain     string
	Status        string
	TrialEndsAt   time.Time
	ContactName   string
	ContactEmail  string
	ContactPhone  string
}

// Organization 组织实体(简化版)
type Organization struct {
	ID       string
	TenantID string
	Name     string
	Type     string
	ParentID *string
	Path     string
	Level    int
	Status   string
}

// User 用户实体(简化版)
type User struct {
	ID             string
	TenantID       string
	Name           string
	Email          string
	Phone          string
	OrganizationID string
	Role           string
	Status         string
}

// SagaDefinition Saga定义
type SagaDefinition struct {
	SagaID   string      // Saga唯一标识
	SagaType string      // Saga类型
	Steps    []*SagaStep // 执行步骤
	Payload  interface{} // 负载数据
}

// SagaStep Saga步骤
type SagaStep struct {
	Name       string        // 步骤名称
	Execute    StepFunc      // 执行函数
	Compensate CompensateFunc // 补偿函数
}

// StepFunc 执行函数类型
type StepFunc func(ctx context.Context, data interface{}) (interface{}, error)

// CompensateFunc 补偿函数类型
type CompensateFunc func(ctx context.Context, data interface{}) error

// generateID 生成唯一ID
func generateID() string {
	return uuid.New().String()
}

// NewTenantRegistrationSaga 创建租户注册Saga
//
// 业务流程:
// 1. 创建租户记录
// 2. 创建默认组织
// 3. 创建管理员账户
// 4. 发送欢迎邮件
func NewTenantRegistrationSaga() *SagaDefinition {
	return &SagaDefinition{
		SagaID:    "saga-tenant-registration",
		SagaType:  "tenant_registration",
		Steps: []*SagaStep{
			{
				Name: "create_tenant",
				Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
					cmd := data.(*TenantRegistrationCommand)
					tenant := &Tenant{
						ID:           generateID(),
						Name:         cmd.EnterpriseName,
						Industry:     cmd.Industry,
						Scale:        cmd.Scale,
						Subdomain:    cmd.Subdomain,
						Status:       "trial",
						TrialEndsAt:  time.Now().Add(14 * 24 * time.Hour),
						ContactName:  cmd.ContactName,
						ContactEmail: cmd.ContactEmail,
						ContactPhone: cmd.ContactPhone,
					}
					// 实际场景: tenantRepo.Save(ctx, tenant)
					return tenant, nil
				},
				Compensate: func(ctx context.Context, data interface{}) error {
					tenant := data.(*Tenant)
					// 实际场景: tenantRepo.Delete(ctx, tenant.ID)
					fmt.Printf("[补偿] 删除租户: %s\n", tenant.ID)
					return nil
				},
			},
			{
				Name: "create_organization",
				Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
					tenant := data.(*Tenant)
					org := &Organization{
						ID:       generateID(),
						TenantID: tenant.ID,
						Name:     tenant.Name + "-默认组织",
						Type:     "company",
						ParentID: nil,
						Path:     "/" + generateID(),
						Level:    1,
						Status:   "active",
					}
					// 实际场景: orgRepo.Save(ctx, org)
					return org, nil
				},
				Compensate: func(ctx context.Context, data interface{}) error {
					org := data.(*Organization)
					// 实际场景: orgRepo.Delete(ctx, org.ID)
					fmt.Printf("[补偿] 删除组织: %s\n", org.ID)
					return nil
				},
			},
			{
				Name: "create_admin_user",
				Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
					org := data.(*Organization)
					user := &User{
						ID:             generateID(),
						TenantID:       org.TenantID,
						Name:           "Admin",
						Email:          fmt.Sprintf("admin@%s.com", org.TenantID),
						Phone:          "13800138000",
						OrganizationID: org.ID,
						Role:           "admin",
						Status:         "active",
					}
					// 实际场景: userRepo.Save(ctx, user)
					return user, nil
				},
				Compensate: func(ctx context.Context, data interface{}) error {
					user := data.(*User)
					// 实际场景: userRepo.Delete(ctx, user.ID)
					fmt.Printf("[补偿] 删除用户: %s\n", user.ID)
					return nil
				},
			},
			{
				Name: "send_welcome_email",
				Execute: func(ctx context.Context, data interface{}) (interface{}, error) {
					user := data.(*User)
					// 实际场景: emailService.Send(ctx, user.Email, "welcome")
					fmt.Printf("[执行] 发送欢迎邮件给: %s\n", user.Email)
					return nil, nil
				},
				Compensate: func(ctx context.Context, data interface{}) error {
					user := data.(*User)
					// 发送失败通知邮件,不阻塞补偿流程
					fmt.Printf("[补偿] 发送注册失败通知给: %s\n", user.Email)
					return nil
				},
			},
		},
		Payload: nil,
	}
}
