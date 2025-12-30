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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/coze-dev/coze-studio/backend/domain/org/entity"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockWorkExperienceRepository 工作经历仓储Mock
type MockWorkExperienceRepository struct {
	mock.Mock
}

func (m *MockWorkExperienceRepository) Create(ctx context.Context, exp *entity.WorkExperience) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockWorkExperienceRepository) Update(ctx context.Context, exp *entity.WorkExperience) error {
	args := m.Called(ctx, exp)
	return args.Error(0)
}

func (m *MockWorkExperienceRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWorkExperienceRepository) GetByID(ctx context.Context, id int64) (*entity.WorkExperience, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WorkExperience), args.Error(1)
}

func (m *MockWorkExperienceRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.WorkExperience, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.WorkExperience), args.Error(1)
}

// MockEducationRepository 教育经历仓储Mock
type MockEducationRepository struct {
	mock.Mock
}

func (m *MockEducationRepository) Create(ctx context.Context, edu *entity.Education) error {
	args := m.Called(ctx, edu)
	return args.Error(0)
}

func (m *MockEducationRepository) Update(ctx context.Context, edu *entity.Education) error {
	args := m.Called(ctx, edu)
	return args.Error(0)
}

func (m *MockEducationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEducationRepository) GetByID(ctx context.Context, id int64) (*entity.Education, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Education), args.Error(1)
}

func (m *MockEducationRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.Education, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Education), args.Error(1)
}

// MockMemberSkillRepository 技能仓储Mock
type MockMemberSkillRepository struct {
	mock.Mock
}

func (m *MockMemberSkillRepository) Create(ctx context.Context, skill *entity.MemberSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

func (m *MockMemberSkillRepository) Update(ctx context.Context, skill *entity.MemberSkill) error {
	args := m.Called(ctx, skill)
	return args.Error(0)
}

func (m *MockMemberSkillRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMemberSkillRepository) GetByID(ctx context.Context, id int64) (*entity.MemberSkill, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MemberSkill), args.Error(1)
}

func (m *MockMemberSkillRepository) ListByUserID(ctx context.Context, userID string) ([]*entity.MemberSkill, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MemberSkill), args.Error(1)
}

func (m *MockMemberSkillRepository) ListByTenantID(ctx context.Context, tenantID string) ([]*entity.MemberSkill, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MemberSkill), args.Error(1)
}

func (m *MockMemberSkillRepository) SearchByName(ctx context.Context, tenantID, skillName string, limit int) ([]*entity.MemberSkill, error) {
	args := m.Called(ctx, tenantID, skillName, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.MemberSkill), args.Error(1)
}

// MockEmployeeRepository 员工仓储Mock
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, emp *entity.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, emp *entity.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, empID string) error {
	args := m.Called(ctx, empID)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, empID string) (*entity.Employee, error) {
	args := m.Called(ctx, empID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByUserID(ctx context.Context, userID string) (*entity.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) List(ctx context.Context, filter interface{}) ([]*entity.Employee, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.Employee), args.Get(1).(int64), args.Error(2)
}

// TestMemberProfileService_AddWorkExperience 测试添加工作经历
func TestMemberProfileService_AddWorkExperience(t *testing.T) {
	// 准备测试数据
	ctx := context.Background()
	userID := "test-user-001"
	tenantID := "test-tenant-001"

	mockEmpRepo := new(MockEmployeeRepository)
	mockWorkExpRepo := new(MockWorkExperienceRepository)
	mockEduRepo := new(MockEducationRepository)
	mockSkillRepo := new(MockMemberSkillRepository)

	service := NewMemberProfileService(
		mockWorkExpRepo,
		mockEduRepo,
		mockSkillRepo,
		mockEmpRepo,
	)

	// 设置Mock期望
	emp := &entity.Employee{
		EmpID:    "emp-001",
		UserID:   userID,
		TenantID: tenantID,
	}
	mockEmpRepo.On("GetByUserID", ctx, userID).Return(emp, nil)
	mockWorkExpRepo.On("Create", ctx, mock.Anything).Return(nil)

	// 测试用例1: 正常添加
	t.Run("正常添加工作经历", func(t *testing.T) {
		exp := &entity.WorkExperience{
			UserID:      userID,
			TenantID:    tenantID,
			CompanyName: "Tech Corp",
			Position:    "Software Engineer",
			StartDate:   time.Now().UnixMilli(),
		}

		err := service.AddWorkExperience(ctx, exp)
		assert.NoError(t, err)
		assert.NotZero(t, exp.CreatedAt)
		assert.NotZero(t, exp.UpdatedAt)
	})

	// 测试用例2: 参数验证失败
	t.Run("参数验证失败", func(t *testing.T) {
		exp := &entity.WorkExperience{
			UserID:   "", // 空用户ID
			TenantID: tenantID,
		}

		err := service.AddWorkExperience(ctx, exp)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrInvalidParam, err)
	})

	// 测试用例3: 员工不存在
	t.Run("员工不存在", func(t *testing.T) {
		mockEmpRepo2 := new(MockEmployeeRepository)
		service2 := NewMemberProfileService(
			new(MockWorkExperienceRepository),
			mockEduRepo,
			mockSkillRepo,
			mockEmpRepo2,
		)

		mockEmpRepo2.On("GetByUserID", ctx, "nonexistent-user").Return(nil, nil)

		exp := &entity.WorkExperience{
			UserID:      "nonexistent-user",
			TenantID:    tenantID,
			CompanyName: "Tech Corp",
			Position:    "Software Engineer",
		}

		err := service2.AddWorkExperience(ctx, exp)
		assert.Error(t, err)
		assert.Equal(t, errno.ErrEmployeeNotFound, err)
	})

	// 验证Mock调用
	mockEmpRepo.AssertExpectations(t)
	mockWorkExpRepo.AssertExpectations(t)
}

// TestMemberProfileService_GenerateResume 测试生成简历
func TestMemberProfileService_GenerateResume(t *testing.T) {
	ctx := context.Background()
	userID := "test-user-002"
	tenantID := "test-tenant-002"

	mockEmpRepo := new(MockEmployeeRepository)
	mockWorkExpRepo := new(MockWorkExperienceRepository)
	mockEduRepo := new(MockEducationRepository)
	mockSkillRepo := new(MockMemberSkillRepository)

	service := NewMemberProfileService(
		mockWorkExpRepo,
		mockEduRepo,
		mockSkillRepo,
		mockEmpRepo,
	)

	// 准备测试数据
	emp := &entity.Employee{
		EmpID:    "emp-002",
		UserID:   userID,
		TenantID: tenantID,
		EmpName:  "张三",
	}

	workExps := []*entity.WorkExperience{
		{
			ID:          1,
			UserID:      userID,
			TenantID:    tenantID,
			CompanyName: "公司A",
			Position:    "工程师",
			StartDate:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).UnixMilli(),
		},
	}

	edus := []*entity.Education{
		{
			ID:         1,
			UserID:     userID,
			TenantID:   tenantID,
			SchoolName: "清华大学",
			Major:      "计算机科学",
			Degree:     entity.DegreeBachelor,
		},
	}

	skills := []*entity.MemberSkill{
		{
			ID:          1,
			UserID:      userID,
			TenantID:    tenantID,
			SkillName:   "Go",
			Proficiency: entity.SkillLevelExpert,
		},
	}

	// 设置Mock期望
	mockEmpRepo.On("GetByUserID", ctx, userID).Return(emp, nil)
	mockWorkExpRepo.On("ListByUserID", ctx, userID).Return(workExps, nil)
	mockEduRepo.On("ListByUserID", ctx, userID).Return(edus, nil)
	mockSkillRepo.On("ListByUserID", ctx, userID).Return(skills, nil)

	// 执行测试
	resume, err := service.GenerateResume(ctx, userID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, resume)
	assert.Equal(t, userID, resume.UserID)
	assert.Equal(t, tenantID, resume.TenantID)
	assert.NotNil(t, resume.Employee)
	assert.Len(t, resume.WorkExperiences, 1)
	assert.Len(t, resume.Educations, 1)
	assert.Len(t, resume.Skills, 1)
	assert.NotZero(t, resume.GeneratedAt)

	// 验证简历方法
	assert.Equal(t, entity.DegreeBachelor, resume.GetHighestDegree())
	assert.Contains(t, resume.GetExpertSkills(), "Go")

	// 验证Mock调用
	mockEmpRepo.AssertExpectations(t)
	mockWorkExpRepo.AssertExpectations(t)
	mockEduRepo.AssertExpectations(t)
	mockSkillRepo.AssertExpectations(t)
}
