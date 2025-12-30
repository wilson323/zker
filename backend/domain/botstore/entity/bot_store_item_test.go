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

package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestBotStoreItem_CanBePublished 测试是否可以发布
func TestBotStoreItem_CanBePublished(t *testing.T) {
	tests := []struct {
		name    string
		status  BotStoreItemStatus
		wantErr bool
	}{
		{
			name:    "草稿状态可以发布",
			status:  BotStoreItemStatusDraft,
			wantErr: false,
		},
		{
			name:    "已拒绝状态可以发布",
			status:  BotStoreItemStatusRejected,
			wantErr: false,
		},
		{
			name:    "下架状态可以发布",
			status:  BotStoreItemStatusOffline,
			wantErr: false,
		},
		{
			name:    "已发布状态不能再次发布",
			status:  BotStoreItemStatusPublished,
			wantErr: true,
		},
		{
			name:    "待审核状态不能发布",
			status:  BotStoreItemStatusPending,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &BotStoreItem{
				ItemID: "item123",
				Status: tt.status,
			}
			err := item.CanBePublished()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBotStoreItem_IsPublished 测试是否已发布
func TestBotStoreItem_IsPublished(t *testing.T) {
	tests := []struct {
		name   string
		status BotStoreItemStatus
		want   bool
	}{
		{
			name:   "已发布状态",
			status: BotStoreItemStatusPublished,
			want:   true,
		},
		{
			name:   "草稿状态",
			status: BotStoreItemStatusDraft,
			want:   false,
		},
		{
			name:   "待审核状态",
			status: BotStoreItemStatusPending,
			want:   false,
		},
		{
			name:   "已拒绝状态",
			status: BotStoreItemStatusRejected,
			want:   false,
		},
		{
			name:   "下架状态",
			status: BotStoreItemStatusOffline,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &BotStoreItem{
				ItemID: "item123",
				Status: tt.status,
			}
			got := item.IsPublished()
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestBotStoreItem_CanBeReviewed 测试是否可以审核
func TestBotStoreItem_CanBeReviewed(t *testing.T) {
	tests := []struct {
		name    string
		status  BotStoreItemStatus
		wantErr bool
	}{
		{
			name:    "待审核状态可以审核",
			status:  BotStoreItemStatusPending,
			wantErr: false,
		},
		{
			name:    "草稿状态不能审核",
			status:  BotStoreItemStatusDraft,
			wantErr: true,
		},
		{
			name:    "已发布状态不能审核",
			status:  BotStoreItemStatusPublished,
			wantErr: true,
		},
		{
			name:    "已拒绝状态不能审核",
			status:  BotStoreItemStatusRejected,
			wantErr: true,
		},
		{
			name:    "下架状态不能审核",
			status:  BotStoreItemStatusOffline,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &BotStoreItem{
				ItemID: "item123",
				Status: tt.status,
			}
			err := item.CanBeReviewed()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBotStoreItem_TableName 测试表名
func TestBotStoreItem_TableName(t *testing.T) {
	item := &BotStoreItem{}
	assert.Equal(t, "bot_store_items", item.TableName())
}

// TestBotStoreItem_Fields 测试字段
func TestBotStoreItem_Fields(t *testing.T) {
	now := time.Now()
	item := &BotStoreItem{
		ItemID:        "item123",
		BotID:         "bot123",
		TenantID:      "tenant123",
		Name:          "Test Bot",
		Description:   "A test bot",
		Category:      "cat_productivity",
		Tags:          []string{"test", "demo"},
		Price:         0.0,
		PublisherID:   "user123",
		Status:        BotStoreItemStatusPublished,
		ViewCount:     100,
		DownloadCount: 50,
		Rating:        4.5,
		RatingCount:   10,
		Screenshots:   []string{"screenshot1.jpg", "screenshot2.jpg"},
		Version:       "1.0.0",
		RejectReason:  "",
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}

	assert.Equal(t, "item123", item.ItemID)
	assert.Equal(t, "bot123", item.BotID)
	assert.Equal(t, "tenant123", item.TenantID)
	assert.Equal(t, "Test Bot", item.Name)
	assert.Equal(t, "A test bot", item.Description)
	assert.Equal(t, "cat_productivity", item.Category)
	assert.Len(t, item.Tags, 2)
	assert.Equal(t, 0.0, item.Price)
	assert.Equal(t, "user123", item.PublisherID)
	assert.Equal(t, BotStoreItemStatusPublished, item.Status)
	assert.Equal(t, 100, item.ViewCount)
	assert.Equal(t, 50, item.DownloadCount)
	assert.Equal(t, 4.5, item.Rating)
	assert.Equal(t, 10, item.RatingCount)
	assert.Len(t, item.Screenshots, 2)
	assert.Equal(t, "1.0.0", item.Version)
	assert.Empty(t, item.RejectReason)
	assert.False(t, item.DeletedAt != nil)
}
