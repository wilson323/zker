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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/coze-dev/coze-studio/backend/domain/botstore/entity"
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// MockBotStoreReviewRepository Mock评论仓库
type MockBotStoreReviewRepository struct {
	mock.Mock
}

func (m *MockBotStoreReviewRepository) Create(ctx context.Context, review *entity.BotStoreReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockBotStoreReviewRepository) Update(ctx context.Context, review *entity.BotStoreReview) error {
	args := m.Called(ctx, review)
	return args.Error(0)
}

func (m *MockBotStoreReviewRepository) Delete(ctx context.Context, reviewID string) error {
	args := m.Called(ctx, reviewID)
	return args.Error(0)
}

func (m *MockBotStoreReviewRepository) GetByID(ctx context.Context, reviewID string) (*entity.BotStoreReview, error) {
	args := m.Called(ctx, reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BotStoreReview), args.Error(1)
}

func (m *MockBotStoreReviewRepository) GetByUserAndItem(ctx context.Context, userID, itemID string) (*entity.BotStoreReview, error) {
	args := m.Called(ctx, userID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.BotStoreReview), args.Error(1)
}

func (m *MockBotStoreReviewRepository) ListByItemID(ctx context.Context, itemID string, page, pageSize int, sortBy string) ([]*entity.BotStoreReview, int, error) {
	args := m.Called(ctx, itemID, page, pageSize, sortBy)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.BotStoreReview), args.Int(1), args.Error(2)
}

func (m *MockBotStoreReviewRepository) ListByUserID(ctx context.Context, userID string, page, pageSize int) ([]*entity.BotStoreReview, int, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.BotStoreReview), args.Int(1), args.Error(2)
}

func (m *MockBotStoreReviewRepository) GetStatistics(ctx context.Context, itemID string) (*entity.ReviewStatistics, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ReviewStatistics), args.Error(1)
}

func (m *MockBotStoreReviewRepository) UpdateItemRating(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockBotStoreReviewRepository) CountByItemID(ctx context.Context, itemID string) (int, error) {
	args := m.Called(ctx, itemID)
	return args.Int(0), args.Error(1)
}

func (m *MockBotStoreReviewRepository) HasUserReviewed(ctx context.Context, userID, itemID string) (bool, error) {
	args := m.Called(ctx, userID, itemID)
	return args.Bool(0), args.Error(1)
}

// 测试创建评论
func TestCreateReview(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	userID := "user123"
	tenantID := "tenant123"
	itemID := "item123"

	req := &CreateReviewRequest{
		ItemID:  itemID,
		Rating:  5,
		Comment: "Great bot!",
	}

	// Mock商品存在且已发布
	publishedItem := &entity.BotStoreItem{
		ItemID: itemID,
		Status: entity.BotStoreItemStatusPublished,
	}
	mockStoreRepo.On("GetByID", ctx, itemID).Return(publishedItem, nil)
	mockReviewRepo.On("HasUserReviewed", ctx, userID, itemID).Return(false, nil)
	mockReviewRepo.On("Create", ctx, mock.AnythingOfType("*entity.BotStoreReview")).Return(nil)
	mockReviewRepo.On("UpdateItemRating", ctx, itemID).Return(nil)

	// 执行
	review, err := service.CreateReview(ctx, req, userID, tenantID)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, itemID, review.ItemID)
	assert.Equal(t, userID, review.UserID)
	assert.Equal(t, tenantID, review.TenantID)
	assert.Equal(t, 5, review.Rating)
	assert.Equal(t, "Great bot!", review.Comment)

	mockStoreRepo.AssertExpectations(t)
	mockReviewRepo.AssertExpectations(t)
}

// 测试创建评论 - 评分无效
func TestCreateReviewInvalidRating(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	req := &CreateReviewRequest{
		ItemID:  "item123",
		Rating:  6, // 无效评分
		Comment: "Great bot!",
	}

	// 执行
	_, err := service.CreateReview(ctx, req, "user123", "tenant123")

	// 验证
	assert.Error(t, err)
	assert.Equal(t, errno.ErrInvalidRating, err)
}

// 测试创建评论 - 商品不存在
func TestCreateReviewItemNotFound(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	req := &CreateReviewRequest{
		ItemID:  "nonexistent",
		Rating:  5,
		Comment: "Great bot!",
	}

	mockStoreRepo.On("GetByID", ctx, "nonexistent").Return(nil, errors.New("not found"))

	// 执行
	_, err := service.CreateReview(ctx, req, "user123", "tenant123")

	// 验证
	assert.Error(t, err)
	assert.Equal(t, errno.ErrBotStoreItemNotFound, err)
}

// 测试创建评论 - 用户已评论
func TestCreateReviewAlreadyReviewed(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	userID := "user123"
	itemID := "item123"

	req := &CreateReviewRequest{
		ItemID:  itemID,
		Rating:  5,
		Comment: "Great bot!",
	}

	publishedItem := &entity.BotStoreItem{
		ItemID: itemID,
		Status: entity.BotStoreItemStatusPublished,
	}

	mockStoreRepo.On("GetByID", ctx, itemID).Return(publishedItem, nil)
	mockReviewRepo.On("HasUserReviewed", ctx, userID, itemID).Return(true, nil)

	// 执行
	_, err := service.CreateReview(ctx, req, userID, "tenant123")

	// 验证
	assert.Error(t, err)
	assert.Equal(t, errno.ErrAlreadyReviewed, err)
}

// 测试更新评论
func TestUpdateReview(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	userID := "user123"
	reviewID := "review123"

	existingReview := &entity.BotStoreReview{
		ReviewID: reviewID,
		UserID:   userID,
		ItemID:   "item123",
		Rating:   4,
		Comment:  "Good",
	}

	req := &UpdateReviewRequest{
		Rating:  5,
		Comment: "Excellent!",
	}

	mockReviewRepo.On("GetByID", ctx, reviewID).Return(existingReview, nil)
	mockReviewRepo.On("Update", ctx, mock.AnythingOfType("*entity.BotStoreReview")).Return(nil)
	mockReviewRepo.On("UpdateItemRating", ctx, "item123").Return(nil)

	// 执行
	err := service.UpdateReview(ctx, reviewID, req, userID)

	// 验证
	assert.NoError(t, err)
	mockReviewRepo.AssertExpectations(t)
}

// 测试更新评论 - 权限错误
func TestUpdateReviewPermissionDenied(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	reviewID := "review123"

	existingReview := &entity.BotStoreReview{
		ReviewID: reviewID,
		UserID:   "user456", // 不同的用户
		ItemID:   "item123",
		Rating:   4,
		Comment:  "Good",
	}

	req := &UpdateReviewRequest{
		Rating:  5,
		Comment: "Excellent!",
	}

	mockReviewRepo.On("GetByID", ctx, reviewID).Return(existingReview, nil)

	// 执行
	err := service.UpdateReview(ctx, reviewID, req, "user123")

	// 验证
	assert.Error(t, err)
	assert.Equal(t, errno.ErrPermissionDenied, err)
}

// 测试删除评论
func TestDeleteReview(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	userID := "user123"
	reviewID := "review123"

	existingReview := &entity.BotStoreReview{
		ReviewID: reviewID,
		UserID:   userID,
		ItemID:   "item123",
		Rating:   4,
		Comment:  "Good",
	}

	mockReviewRepo.On("GetByID", ctx, reviewID).Return(existingReview, nil)
	mockReviewRepo.On("Delete", ctx, reviewID).Return(nil)
	mockReviewRepo.On("UpdateItemRating", ctx, "item123").Return(nil)

	// 执行
	err := service.DeleteReview(ctx, reviewID, userID)

	// 验证
	assert.NoError(t, err)
	mockReviewRepo.AssertExpectations(t)
}

// 测试获取评论列表
func TestListReviews(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	itemID := "item123"

	req := &ListReviewsRequest{
		ItemID:   itemID,
		Page:     1,
		PageSize: 20,
		SortBy:   "latest",
	}

	reviews := []*entity.BotStoreReview{
		{
			ReviewID: "review1",
			ItemID:   itemID,
			UserID:   "user1",
			Rating:   5,
			Comment:  "Excellent!",
		},
		{
			ReviewID: "review2",
			ItemID:   itemID,
			UserID:   "user2",
			Rating:   4,
			Comment:  "Good",
		},
	}

	stats := &entity.ReviewStatistics{
		AverageRating: 4.5,
		RatingCount:   2,
	}

	publishedItem := &entity.BotStoreItem{
		ItemID: itemID,
		Status: entity.BotStoreItemStatusPublished,
	}

	mockStoreRepo.On("GetByID", ctx, itemID).Return(publishedItem, nil)
	mockReviewRepo.On("ListByItemID", ctx, itemID, 1, 20, "latest").Return(reviews, 2, nil)
	mockReviewRepo.On("GetStatistics", ctx, itemID).Return(stats, nil)

	// 执行
	resp, err := service.ListReviews(ctx, req)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Reviews, 2)
	assert.Equal(t, 2, resp.Total)
	assert.Equal(t, 4.5, resp.AverageRating)
	assert.Equal(t, 2, resp.RatingCount)

	mockStoreRepo.AssertExpectations(t)
	mockReviewRepo.AssertExpectations(t)
}

// 测试获取评论统计
func TestGetReviewStatistics(t *testing.T) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	itemID := "item123"

	stats := &entity.ReviewStatistics{
		AverageRating: 4.5,
		RatingCount:   10,
		Rating1Count:  1,
		Rating2Count:  0,
		Rating3Count:  1,
		Rating4Count:  2,
		Rating5Count:  6,
	}

	publishedItem := &entity.BotStoreItem{
		ItemID: itemID,
		Status: entity.BotStoreItemStatusPublished,
	}

	mockStoreRepo.On("GetByID", ctx, itemID).Return(publishedItem, nil)
	mockReviewRepo.On("GetStatistics", ctx, itemID).Return(stats, nil)

	// 执行
	result, err := service.GetReviewStatistics(ctx, itemID)

	// 验证
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 4.5, result.AverageRating)
	assert.Equal(t, 10, result.RatingCount)
	assert.Equal(t, 1, result.Rating1Count)
	assert.Equal(t, 6, result.Rating5Count)

	mockStoreRepo.AssertExpectations(t)
	mockReviewRepo.AssertExpectations(t)
}

// 测试评论验证
func TestReviewValidation(t *testing.T) {
	tests := []struct {
		name    string
		request interface{}
		wantErr error
	}{
		{
			name: "有效创建评论请求",
			request: &CreateReviewRequest{
				ItemID:  "item123",
				Rating:  5,
				Comment: "Great!",
			},
			wantErr: nil,
		},
		{
			name: "评分过低",
			request: &CreateReviewRequest{
				ItemID:  "item123",
				Rating:  0,
				Comment: "Bad",
			},
			wantErr: errno.ErrInvalidRating,
		},
		{
			name: "评分过高",
			request: &CreateReviewRequest{
				ItemID:  "item123",
				Rating:  6,
				Comment: "Too high",
			},
			wantErr: errno.ErrInvalidRating,
		},
		{
			name: "评论过长",
			request: &CreateReviewRequest{
				ItemID:  "item123",
				Rating:  5,
				Comment: string(make([]byte, 1001)), // 超过1000字符
			},
			wantErr: errno.ErrCommentTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch req := tt.request.(type) {
			case *CreateReviewRequest:
				err = req.Validate()
			case *UpdateReviewRequest:
				err = req.Validate()
			}

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// 基准测试
func BenchmarkCreateReview(b *testing.B) {
	ctx := context.Background()
	mockStoreRepo := new(MockBotStoreRepository)
	mockReviewRepo := new(MockBotStoreReviewRepository)
	service := NewBotStoreReviewService(mockStoreRepo, mockReviewRepo)

	userID := "user123"
	tenantID := "tenant123"
	itemID := "item123"

	req := &CreateReviewRequest{
		ItemID:  itemID,
		Rating:  5,
		Comment: "Great bot!",
	}

	publishedItem := &entity.BotStoreItem{
		ItemID: itemID,
		Status: entity.BotStoreItemStatusPublished,
	}

	mockStoreRepo.On("GetByID", ctx, itemID).Return(publishedItem, nil).Maybe()
	mockReviewRepo.On("HasUserReviewed", ctx, userID, itemID).Return(false, nil).Maybe()
	mockReviewRepo.On("Create", ctx, mock.AnythingOfType("*entity.BotStoreReview")).Return(nil).Maybe()
	mockReviewRepo.On("UpdateItemRating", ctx, itemID).Return(nil).Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreateReview(ctx, req, userID, tenantID)
	}
}
