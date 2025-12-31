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

import "time"

// ChannelType 渠道类型枚举
type ChannelType string

const (
	// ChannelTypeWechatOfficialAccount 微信公众号
	ChannelTypeWechatOfficialAccount ChannelType = "wechat_official_account"
	// ChannelTypeFeishuApp 飞书应用
	ChannelTypeFeishuApp ChannelType = "feishu_app"
	// ChannelTypeDiscordBot Discord Bot
	ChannelTypeDiscordBot ChannelType = "discord_bot"
	// ChannelTypeDingTalk 钉钉
	ChannelTypeDingTalk ChannelType = "ding_talk"
	// ChannelTypeSlack Slack
	ChannelTypeSlack ChannelType = "slack"
)

// ChannelStatus 渠道状态
type ChannelStatus string

const (
	// ChannelStatusDraft 草稿
	ChannelStatusDraft ChannelStatus = "draft"
	// ChannelStatusPublished 已发布
	ChannelStatusPublished ChannelStatus = "published"
	// ChannelStatusUnpublished 已取消发布
	ChannelStatusUnpublished ChannelStatus = "unpublished"
	// ChannelStatusError 错误
	ChannelStatusError ChannelStatus = "error"
)

// ChannelConfig 渠道配置实体
type ChannelConfig struct {
	ID               string                 `json:"id" gorm:"primaryKey;size:36"`
	TenantID         string                 `json:"tenant_id" gorm:"not null;index:idx_tenant_channel;size:36"`
	BotID            string                 `json:"bot_id" gorm:"not null;index:idx_bot_channel;size:36"`
	ChannelType      ChannelType           `json:"channel_type" gorm:"not null;type:varchar(50)"`
	Enabled          bool                  `json:"enabled" gorm:"not null;default:true"`
	Config           string                 `json:"config" gorm:"type:json"` // JSON配置
	WebhookURL       string                 `json:"webhook_url" gorm:"type:text"`
	WebhookSecret    string                 `json:"webhook_secret" gorm:"size:255"` // Webhook验证密钥
	Status           ChannelStatus         `json:"status" gorm:"type:varchar(20);default:'draft'"`
	LastErrorMessage string                 `json:"last_error_message" gorm:"type:text"`
	PublishedAt       *time.Time             `json:"published_at" gorm:"index"`
	UnpublishedAt    *time.Time             `json:"unpublished_at"`
	CreatedAt        time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt        *time.Time             `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (ChannelConfig) TableName() string {
	return "channel_configs"
}

// ChannelMessage 渠道消息 (统一格式)
type ChannelMessage struct {
	MessageID   string                 `json:"message_id"`
	ChannelID   string                 `json:"channel_id"`
	UserID      string                 `json:"user_id"`       // 渠道用户ID
	MessageType string                 `json:"message_type"`  // text, image, audio, video, card, file
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// ChannelResponse 渠道响应 (统一格式)
type ChannelResponse struct {
	Success  bool                   `json:"success"`
	MessageID string                 `json:"message_id,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PublishResult 发布结果
type PublishResult struct {
	Success    bool    `json:"success"`
	ChannelID  string  `json:"channel_id"`
	WebhookURL string  `json:"webhook_url,omitempty"`
	QRCodeURL  string  `json:"qr_code_url,omitempty"`
	Error      string  `json:"error,omitempty"`
}
