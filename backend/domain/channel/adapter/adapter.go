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

package adapter

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/channel/entity"
)

// ChannelAdapter 渠道适配器接口
// 所有渠道适配器必须实现此接口
type ChannelAdapter interface {
	// ValidateConfig 验证渠道配置是否有效
	ValidateConfig(ctx context.Context, config map[string]string) error

	// PublishBot 发布Bot到渠道
	// 返回: webhook地址、二维码URL等发布信息
	PublishBot(ctx context.Context, botID string, botName string, config map[string]string) (*entity.PublishResult, error)

	// UnpublishBot 取消发布Bot
	UnpublishBot(ctx context.Context, botID string, config map[string]string) error

	// SendMessage 发送消息到渠道
	SendMessage(ctx context.Context, msg *entity.ChannelMessage) (*entity.ChannelResponse, error)

	// HandleWebhook 处理渠道Webhook回调
	HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*entity.ChannelMessage, error)

	// GetChannelType 获取渠道类型
	GetChannelType() entity.ChannelType

	// TransformMessage 转换消息格式 (从内部格式到渠道格式)
	TransformMessage(ctx context.Context, msg *entity.ChannelMessage) (interface{}, error)

	// TransformResponse 转换响应格式 (从渠道格式到内部格式)
	TransformResponse(ctx context.Context, resp interface{}) (*entity.ChannelResponse, error)
}

// AdapterFactory 适配器工厂接口
type AdapterFactory interface {
	// CreateAdapter 创建渠道适配器
	CreateAdapter(channelType entity.ChannelType) (ChannelAdapter, error)

	// GetSupportedChannels 获取支持的渠道列表
	GetSupportedChannels() []entity.ChannelType
}

// BaseAdapter 基础适配器，提供通用功能
type BaseAdapter struct {
	ChannelType entity.ChannelType
}

// GetChannelType 获取渠道类型
func (b *BaseAdapter) GetChannelType() entity.ChannelType {
	return b.ChannelType
}
