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
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/digital_employee/entity"
)

// MockProfileRepository 模拟ProfileRepository
type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) Create(ctx context.Context, profile *entity.EmployeeProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) Update(ctx context.Context, profile *entity.EmployeeProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockProfileRepository) GetByID(ctx context.Context, employeeID string) (*entity.EmployeeProfile, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeProfile), args.Error(1)
}

func (m *MockProfileRepository) GetByTenantID(ctx context.Context, tenantID string, role *entity.EmployeeRole, status *entity.EmployeeStatus, page, pageSize int) ([]*entity.EmployeeProfile, int64, error) {
	args := m.Called(ctx, tenantID, role, status, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*entity.EmployeeProfile), args.Get(1).(int64), args.Error(2)
}

func (m *MockProfileRepository) Delete(ctx context.Context, employeeID string) error {
	args := m.Called(ctx, employeeID)
	return args.Error(0)
}

func (m *MockProfileRepository) GetByBotID(ctx context.Context, botID string) (*entity.EmployeeProfile, error) {
	args := m.Called(ctx, botID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmployeeProfile), args.Error(1)
}

func (m *MockProfileRepository) ExistsByName(ctx context.Context, tenantID, name string) (bool, error) {
	args := m.Called(ctx, tenantID, name)
	return args.Bool(0), args.Error(1)
}

// TestCreateProfile 测试创建员工画像
func TestCreateProfile(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	req := &entity.CreateProfileRequest{
		TenantID:    "tenant-123",
		Name:        "测试员工",
		Role:        entity.EmployeeRoleCustomerService,
		Skills:      []string{"客户咨询", "问题解答"},
		Personality: "热情友好",
	}

	// Mock: 名称不存在
	mockRepo.On("ExistsByName", ctx, req.TenantID, req.Name).Return(false, nil)
	// Mock: 创建成功
	mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.EmployeeProfile")).Return(nil)

	profile, err := service.CreateProfile(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, req.Name, profile.Name)
	assert.Equal(t, req.Role, profile.Role)
	assert.Equal(t, entity.EmployeeStatusActive, profile.Status)

	mockRepo.AssertExpectations(t)
}

// TestCreateProfile_DuplicateName 测试创建员工画像-名称重复
func TestCreateProfile_DuplicateName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	req := &entity.CreateProfileRequest{
		TenantID: "tenant-123",
		Name:     "重复名称",
		Role:     entity.EmployeeRoleSales,
		Skills:   []string{"销售"},
	}

	// Mock: 名称已存在
	mockRepo.On("ExistsByName", ctx, req.TenantID, req.Name).Return(true, nil)

	profile, err := service.CreateProfile(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "already exists")

	mockRepo.AssertExpectations(t)
}

// TestGetProfile 测试获取员工画像
func TestGetProfile(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	employeeID := "employee-123"
	expectedProfile := &entity.EmployeeProfile{
		EmployeeID: employeeID,
		Name:       "测试员工",
		Role:       entity.EmployeeRoleTechSupport,
		Status:     entity.EmployeeStatusActive,
	}

	// Mock: 获取成功
	mockRepo.On("GetByID", ctx, employeeID).Return(expectedProfile, nil)

	profile, err := service.GetProfile(ctx, employeeID)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, employeeID, profile.EmployeeID)
	assert.Equal(t, "测试员工", profile.Name)

	mockRepo.AssertExpectations(t)
}

// TestGetProfile_NotFound 测试获取员工画像-不存在
func TestGetProfile_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	employeeID := "non-existent"

	// Mock: 不存在
	mockRepo.On("GetByID", ctx, employeeID).Return(nil, assert.AnError)

	profile, err := service.GetProfile(ctx, employeeID)

	assert.Error(t, err)
	assert.Nil(t, profile)

	mockRepo.AssertExpectations(t)
}

// TestListProfiles 测试列出员工画像
func TestListProfiles(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	req := &entity.ListProfilesRequest{
		TenantID: "tenant-123",
		Page:     1,
		PageSize: 20,
	}

	expectedProfiles := []*entity.EmployeeProfile{
		{
			EmployeeID: "employee-1",
			Name:       "员工1",
			Role:       entity.EmployeeRoleCustomerService,
		},
		{
			EmployeeID: "employee-2",
			Name:       "员工2",
			Role:       entity.EmployeeRoleSales,
		},
	}

	// Mock: 查询成功
	mockRepo.On("GetByTenantID", ctx, req.TenantID, (*entity.EmployeeRole)(nil), (*entity.EmployeeStatus)(nil), req.Page, req.PageSize).
		Return(expectedProfiles, int64(2), nil)

	resp, err := service.ListProfiles(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Profiles))
	assert.Equal(t, int64(2), resp.Total)
	assert.Equal(t, req.Page, resp.Page)
	assert.Equal(t, req.PageSize, resp.PageSize)

	mockRepo.AssertExpectations(t)
}

// TestDeleteProfile 测试删除员工画像
func TestDeleteProfile(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockProfileRepository)
	service := NewProfileService(mockRepo)

	employeeID := "employee-123"

	// Mock: 删除成功
	mockRepo.On("Delete", ctx, employeeID).Return(nil)

	err := service.DeleteProfile(ctx, employeeID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
