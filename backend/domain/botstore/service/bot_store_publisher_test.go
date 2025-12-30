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

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockBotStoreRepository Mock仓库
type MockBotStoreRepository struct {
	mock.Mock
}

func (m *MockBotStoreRepository) Create(ctx context.Context, item *entity.BotStoreItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockBotStoreRepository) Update(ctx context.Context, item *entity.BotStoreItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockBotStoreRepository) GetByID(ctx context.Context, itemID string) (*entity.BotStoreItem, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BotStoreItem), args.Error(1)
}

func (m *MockBotStoreRepository) GetByBotID(ctx context.Context, botID string) (*entity.BotStoreItem, error) {
	args := m.Called(ctx, botID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BotStoreItem), args.Error(1)
}

func (m *MockBotStoreRepository) Delete(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockBotStoreRepository) List(ctx context.Context, req *repository.ListRequest) ([]*entity.BotStoreItem, int, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.BotStoreItem), args.Int(1), args.Error(2)
}

func (m *MockBotStoreRepository) Search(ctx context.Context, req *repository.SearchRequest) ([]*entity.BotStoreItem, int, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.BotStoreItem), args.Int(1), args.Error(2)
}

func (m *MockBotStoreRepository) GetPendingReviews(ctx context.Context, page, pageSize int) ([]*entity.BotStoreItem, int, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.BotStoreItem), args.Int(1), args.Error(2)
}

func (m *MockBotStoreRepository) IncrementViewCount(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockBotStoreRepository) IncrementDownloadCount(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockBotStoreRepository) UpdateRating(ctx context.Context, itemID string, rating float64) error {
	args := m.Called(ctx, itemID, rating)
	return args.Error(0)
}

// MockBotStoreCategoryRepository Mock分类仓库
type MockBotStoreCategoryRepository struct {
	mock.Mock
}

func (m *MockBotStoreCategoryRepository) Create(ctx context.Context, category *entity.BotStoreCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockBotStoreCategoryRepository) Update(ctx context.Context, category *entity.BotStoreCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockBotStoreCategoryRepository) GetByID(ctx context.Context, categoryID string) (*entity.BotStoreCategory, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BotStoreCategory), args.Error(1)
}

func (m *MockBotStoreCategoryRepository) List(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BotStoreCategory), args.Error(1)
}

func (m *MockBotStoreCategoryRepository) ListActive(ctx context.Context) ([]*entity.BotStoreCategory, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.BotStoreCategory), args.Error(1)
}

func (m *MockBotStoreCategoryRepository) Delete(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

func (m *MockBotStoreCategoryRepository) IncrementBotCount(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

func (m *MockBotStoreCategoryRepository) DecrementBotCount(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

// TestBotStorePublisher_PublishBot 测试发布Bot
func TestBotStorePublisher_PublishBot(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	publisher := NewBotStorePublisher(mockStoreRepo, mockCategoryRepo)

	req := &PublishBotRequest{
		BotID:       "bot123",
		Name:        "Test Bot",
		Description: "A test bot",
		Category:    "cat_productivity",
		Tags:        []string{"test", "demo"},
		Price:       0.0,
		PublisherID: "user123",
		TenantID:    "tenant123",
	}

	// Mock分类存在
	mockCategoryRepo.On("GetByID", ctx, "cat_productivity").Return(&entity.BotStoreCategory{
		CategoryID: "cat_productivity",
		Name:       "效率工具",
		IsActive:   true,
	}, nil)

	// Mock Bot不存在
	mockStoreRepo.On("GetByBotID", ctx, "bot123").Return(nil, errno.ErrBotStoreItemNotFound)

	// Mock创建成功
	mockStoreRepo.On("Create", ctx, mock.AnythingOfType("*entity.BotStoreItem")).Return(nil)

	// Mock增加分类计数
	mockCategoryRepo.On("IncrementBotCount", ctx, "cat_productivity").Return(nil)

	item, err := publisher.PublishBot(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "bot123", item.BotID)
	assert.Equal(t, "Test Bot", item.Name)
	assert.Equal(t, entity.BotStoreItemStatusPending, item.Status)

	mockStoreRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

// TestBotStorePublisher_PublishBot_AlreadyPublished 测试发布已发布的Bot
func TestBotStorePublisher_PublishBot_AlreadyPublished(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	publisher := NewBotStorePublisher(mockStoreRepo, mockCategoryRepo)

	req := &PublishBotRequest{
		BotID:       "bot123",
		Name:        "Test Bot",
		Description: "A test bot",
		Category:    "cat_productivity",
		Price:       0.0,
		PublisherID: "user123",
		TenantID:    "tenant123",
	}

	// Mock Bot已发布
	existingItem := &entity.BotStoreItem{
		ItemID: "item123",
		BotID:  "bot123",
		Status: entity.BotStoreItemStatusPublished,
	}
	mockStoreRepo.On("GetByBotID", ctx, "bot123").Return(existingItem, nil)

	item, err := publisher.PublishBot(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, item)
	assert.Equal(t, errno.ErrAlreadyPublished, err)

	mockStoreRepo.AssertExpectations(t)
}

// TestBotStorePublisher_UnpublishBot 测试下架Bot
func TestBotStorePublisher_UnpublishBot(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	publisher := NewBotStorePublisher(mockStoreRepo, mockCategoryRepo)

	itemID := "item123"
	userID := "user123"

	// Mock获取项目
	existingItem := &entity.BotStoreItem{
		ItemID:      itemID,
		BotID:       "bot123",
		PublisherID: userID,
		Status:      entity.BotStoreItemStatusPublished,
		Category:    "cat_productivity",
	}
	mockStoreRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	// Mock更新
	mockStoreRepo.On("Update", ctx, mock.AnythingOfType("*entity.BotStoreItem")).Return(nil)

	// Mock减少分类计数
	mockCategoryRepo.On("DecrementBotCount", ctx, "cat_productivity").Return(nil)

	err := publisher.UnpublishBot(ctx, itemID, userID)

	assert.NoError(t, err)
	mockStoreRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

// TestBotStorePublisher_UnpublishBot_PermissionDenied 测试无权限下架
func TestBotStorePublisher_UnpublishBot_PermissionDenied(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	publisher := NewBotStorePublisher(mockStoreRepo, mockCategoryRepo)

	itemID := "item123"
	userID := "user123"

	// Mock获取项目（发布者不是当前用户）
	existingItem := &entity.BotStoreItem{
		ItemID:      itemID,
		BotID:       "bot123",
		PublisherID: "user456", // 不同的用户
		Status:      entity.BotStoreItemStatusPublished,
	}
	mockStoreRepo.On("GetByID", ctx, itemID).Return(existingItem, nil)

	err := publisher.UnpublishBot(ctx, itemID, userID)

	assert.Error(t, err)
	assert.Equal(t, errno.ErrPermissionDenied, err)

	mockStoreRepo.AssertExpectations(t)
}

// TestBotStoreBrowser_ListBots 测试列出Bot
func TestBotStoreBrowser_ListBots(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	browser := NewBotStoreBrowser(mockStoreRepo, mockCategoryRepo)

	req := &ListBotsRequest{
		Category: "cat_productivity",
		SortBy:   "popular",
		Page:     1,
		PageSize: 20,
	}

	// Mock返回数据
	mockItems := []*entity.BotStoreItem{
		{
			ItemID:   "item1",
			Name:     "Bot 1",
			Status:   entity.BotStoreItemStatusPublished,
		},
		{
			ItemID:   "item2",
			Name:     "Bot 2",
			Status:   entity.BotStoreItemStatusPublished,
		},
	}
	mockStoreRepo.On("List", ctx, mock.AnythingOfType("*repository.ListRequest")).Return(mockItems, 2, nil)

	resp, err := browser.ListBots(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Items, 2)
	assert.Equal(t, 2, resp.Total)
	assert.Equal(t, 1, resp.TotalPages)

	mockStoreRepo.AssertExpectations(t)
}

// TestBotStoreReviewer_ReviewBot 测试审核Bot
func TestBotStoreReviewer_ReviewBot(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	reviewer := NewBotStoreReviewer(mockStoreRepo)

	req := &ReviewBotRequest{
		ItemID:     "item123",
		Approved:   true,
		Reason:     "",
		ReviewerID: "admin123",
	}

	// Mock获取待审核项目
	pendingItem := &entity.BotStoreItem{
		ItemID: "item123",
		BotID:  "bot123",
		Status: entity.BotStoreItemStatusPending,
	}
	mockStoreRepo.On("GetByID", ctx, "item123").Return(pendingItem, nil)

	// Mock更新
	mockStoreRepo.On("Update", ctx, mock.AnythingOfType("*entity.BotStoreItem")).Return(nil)

	err := reviewer.ReviewBot(ctx, req)

	assert.NoError(t, err)
	mockStoreRepo.AssertExpectations(t)
}

// TestBotStoreReviewer_ReviewBot_Reject 测试拒绝Bot
func TestBotStoreReviewer_ReviewBot_Reject(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	reviewer := NewBotStoreReviewer(mockStoreRepo)

	req := &ReviewBotRequest{
		ItemID:     "item123",
		Approved:   false,
		Reason:     "不符合规范",
		ReviewerID: "admin123",
	}

	// Mock获取待审核项目
	pendingItem := &entity.BotStoreItem{
		ItemID: "item123",
		BotID:  "bot123",
		Status: entity.BotStoreItemStatusPending,
	}
	mockStoreRepo.On("GetByID", ctx, "item123").Return(pendingItem, nil)

	// Mock更新
	mockStoreRepo.On("Update", ctx, mock.AnythingOfType("*entity.BotStoreItem")).Return(nil)

	err := reviewer.ReviewBot(ctx, req)

	assert.NoError(t, err)
	mockStoreRepo.AssertExpectations(t)
}

// TestBotStoreReviewer_ReviewBot_MissingReason 测试拒绝但未提供原因
func TestBotStoreReviewer_ReviewBot_MissingReason(t *testing.T) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	reviewer := NewBotStoreReviewer(mockStoreRepo)

	req := &ReviewBotRequest{
		ItemID:     "item123",
		Approved:   false,
		Reason:     "", // 缺少原因
		ReviewerID: "admin123",
	}

	// Mock获取待审核项目
	pendingItem := &entity.BotStoreItem{
		ItemID: "item123",
		BotID:  "bot123",
		Status: entity.BotStoreItemStatusPending,
	}
	mockStoreRepo.On("GetByID", ctx, "item123").Return(pendingItem, nil)

	err := reviewer.ReviewBot(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, errno.ErrRejectReasonRequired, err)

	mockStoreRepo.AssertExpectations(t)
}

// BenchmarkBotStorePublisher_PublishBot 性能测试
func BenchmarkBotStorePublisher_PublishBot(b *testing.B) {
	ctx := context.Background()

	mockStoreRepo := new(MockBotStoreRepository)
	mockCategoryRepo := new(MockBotStoreCategoryRepository)
	publisher := NewBotStorePublisher(mockStoreRepo, mockCategoryRepo)

	req := &PublishBotRequest{
		BotID:       "bot123",
		Name:        "Test Bot",
		Description: "A test bot",
		Category:    "cat_productivity",
		Price:       0.0,
		PublisherID: "user123",
		TenantID:    "tenant123",
	}

	// 设置Mock
	mockCategoryRepo.On("GetByID", ctx, "cat_productivity").Return(&entity.BotStoreCategory{
		CategoryID: "cat_productivity",
		Name:       "效率工具",
		IsActive:   true,
	}, nil)
	mockStoreRepo.On("GetByBotID", ctx, "bot123").Return(nil, errno.ErrBotStoreItemNotFound)
	mockStoreRepo.On("Create", ctx, mock.AnythingOfType("*entity.BotStoreItem")).Return(nil)
	mockCategoryRepo.On("IncrementBotCount", ctx, "cat_productivity").Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = publisher.PublishBot(ctx, req)
	}
}
