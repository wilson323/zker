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
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/channel/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// WechatOfficialAccountAdapter 微信公众号适配器
type WechatOfficialAccountAdapter struct {
	BaseAdapter
	httpClient *http.Client
}

// NewWechatOfficialAccountAdapter 创建微信公众号适配器
func NewWechatOfficialAccountAdapter() *WechatOfficialAccountAdapter {
	return &WechatOfficialAccountAdapter{
		BaseAdapter: BaseAdapter{
			ChannelType: entity.ChannelTypeWechatOfficialAccount,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ValidateConfig 验证微信公众号配置
func (a *WechatOfficialAccountAdapter) ValidateConfig(ctx context.Context, config map[string]string) error {
	// 必需字段检查
	requiredFields := []string{"appId", "appSecret"}
	for _, field := range requiredFields {
		if config[field] == "" {
			return errorx.New(errno.ErrChannelInvalidParamCode, errorx.KV("field", field))
		}
	}

	// 可选字段
	// token, encodingAESKey 可以为空，但建议配置

	return nil
}

// PublishBot 发布Bot到微信公众号
func (a *WechatOfficialAccountAdapter) PublishBot(ctx context.Context, botID string, botName string, config map[string]string) (*entity.PublishResult, error) {
	// 1. 创建自定义菜单
	if err := a.createMenu(ctx, botID, config); err != nil {
		logs.CtxWarnf(ctx, "[WechatAdapter] 创建菜单失败: %v", err)
		// 菜单创建失败不影响发布
	}

	// 2. 生成Webhook URL
	webhookURL := fmt.Sprintf("https://%s.com/api/webhooks/wechat/%s", getDomain(ctx), botID)

	// 3. 生成二维码URL
	qrCodeURL := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", "TICKET_PLACEHOLDER")

	return &entity.PublishResult{
		Success:    true,
		ChannelID:  botID,
		WebhookURL: webhookURL,
		QRCodeURL:  qrCodeURL,
	}, nil
}

// UnpublishBot 取消发布Bot
func (a *WechatOfficialAccountAdapter) UnpublishBot(ctx context.Context, botID string, config map[string]string) error {
	// 1. 删除自定义菜单
	if err := a.deleteMenu(ctx, config); err != nil {
		logs.CtxWarnf(ctx, "[WechatAdapter] 删除菜单失败: %v", err)
	}

	return nil
}

// SendMessage 发送消息到微信公众号
func (a *WechatOfficialAccountAdapter) SendMessage(ctx context.Context, msg *entity.ChannelMessage) (*entity.ChannelResponse, error) {
	// 微信公众号API: 客服消息接口
	// https://api.weixin.qq.com/cgi-bin/message/custom/send

	// 根据消息类型转换格式
	_, err := a.TransformMessage(ctx, msg)
	if err != nil {
		return nil, err
	}

	// 发送HTTP请求
	// TODO: 实际调用微信API

	return &entity.ChannelResponse{
		Success:   true,
		MessageID: msg.MessageID,
	}, nil
}

// HandleWebhook 处理微信Webhook回调
func (a *WechatOfficialAccountAdapter) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*entity.ChannelMessage, error) {
	// 1. 验证签名
	if err := a.verifySignature(ctx, payload, headers); err != nil {
		return nil, err
	}

	// 2. 解析XML消息
	var wxMsg WechatMessage
	if err := xml.Unmarshal(payload, &wxMsg); err != nil {
		return nil, errorx.WrapByCode(err, errno.ErrChannelMessageParseFailed)
	}

	// 3. 转换为内部格式
	channelMsg := &entity.ChannelMessage{
		MessageID:   wxMsg.MsgID,
		ChannelID:   wxMsg.ToUserName,
		UserID:      wxMsg.FromUserName,
		MessageType: a.convertMessageType(wxMsg.MsgType),
		Content:     wxMsg.Content,
		Metadata: map[string]interface{}{
			"create_time": wxMsg.CreateTime,
			"msg_type":    wxMsg.MsgType,
			"event":       wxMsg.Event,
		},
		Timestamp: time.Now(),
	}

	return channelMsg, nil
}

// TransformMessage 转换消息格式 (内部格式 -> 微信格式)
func (a *WechatOfficialAccountAdapter) TransformMessage(ctx context.Context, msg *entity.ChannelMessage) (interface{}, error) {
	wechatMsg := map[string]interface{}{
		"touser":  msg.UserID,
		"msgtype": msg.MessageType,
	}

	switch msg.MessageType {
	case "text":
		wechatMsg["text"] = map[string]string{
			"content": msg.Content,
		}
	case "image":
		wechatMsg["image"] = map[string]string{
			"media_id": msg.Content,
		}
	case "news":
		// 图文消息
		wechatMsg["news"] = map[string]interface{}{
			"articles": []map[string]string{
				{
					"title":       msg.Content,
					"description": msg.Metadata["description"].(string),
					"url":          msg.Metadata["url"].(string),
					"picurl":       msg.Metadata["picurl"].(string),
				},
			},
		}
	default:
		return nil, errorx.New(errno.ErrChannelMessageTypeNotSupported)
	}

	return wechatMsg, nil
}

// TransformResponse 转换响应格式 (微信格式 -> 内部格式)
func (a *WechatOfficialAccountAdapter) TransformResponse(ctx context.Context, resp interface{}) (*entity.ChannelResponse, error) {
	// 微信API响应格式解析
	// {"errcode":0,"errmsg":"ok"}
	return &entity.ChannelResponse{
		Success: true,
	}, nil
}

// ==================== 私有方法 ====================

// verifySignature 验证微信签名
func (a *WechatOfficialAccountAdapter) verifySignature(ctx context.Context, payload []byte, headers map[string]string) error {
	signature := headers["signature"]
	timestamp := headers["timestamp"]
	nonce := headers["nonce"]

	// TODO: 从配置中获取token
	token := "YOUR_TOKEN"

	// 1. 将token、timestamp、nonce三个参数进行字典序排序
	params := []string{token, timestamp, nonce}
	sort.Strings(params)

	// 2. 将三个参数字符串拼接成一个字符串进行sha1加密
	str := strings.Join(params, "")
	h := sha1.New()
	h.Write([]byte(str))
	sha1Sum := hex.EncodeToString(h.Sum(nil))

	// 3. 加密后的字符串与signature对比
	if sha1Sum != signature {
		return errorx.New(errno.ErrChannelSignatureInvalid)
	}

	return nil
}

// createMenu 创建自定义菜单
func (a *WechatOfficialAccountAdapter) createMenu(ctx context.Context, botID string, config map[string]string) error {
	// 微信自定义菜单API
	// https://api.weixin.qq.com/cgi-bin/menu/create?access_token=TOKEN

	// TODO: 实际调用微信API创建菜单
	return nil
}

// deleteMenu 删除自定义菜单
func (a *WechatOfficialAccountAdapter) deleteMenu(ctx context.Context, config map[string]string) error {
	// 微信删除菜单API
	// https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=TOKEN

	// TODO: 实际调用微信API删除菜单
	return nil
}

// convertMessageType 转换消息类型
func (a *WechatOfficialAccountAdapter) convertMessageType(wxType string) string {
	typeMap := map[string]string{
		"text":       "text",
		"image":      "image",
		"voice":      "audio",
		"video":      "video",
		"shortvideo": "video",
		"location":   "text",
		"link":       "text",
		"event":      "card",
	}

	if t, ok := typeMap[wxType]; ok {
		return t
	}
	return "text"
}

// getDomain 获取域名
func getDomain(ctx context.Context) string {
	// TODO: 从配置中获取域名
	return "api.example.com"
}

// ==================== 微信消息结构 ====================

// WechatMessage 微信消息结构
type WechatMessage struct {
	XMLName      struct{} `xml:"xml"`
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	MsgID        string `xml:"MsgId"`
	Content      string `xml:"Content"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
}

// WechatApiResponse 微信API响应
type WechatApiResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// ==================== 辅助方法 ====================

// doHTTPRequest 执行HTTP请求
func (a *WechatOfficialAccountAdapter) doHTTPRequest(ctx context.Context, url string, body interface{}) (*WechatApiResponse, error) {
	// TODO: 实现HTTP请求
	return &WechatApiResponse{
		ErrCode: 0,
		ErrMsg:  "ok",
	}, nil
}

// getAccessToken 获取Access Token
func (a *WechatOfficialAccountAdapter) getAccessToken(ctx context.Context, appID, appSecret string) (string, error) {
	// 微信获取Access Token API
	// https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET

	// TODO: 实际调用微信API
	return "ACCESS_TOKEN_PLACEHOLDER", nil
}

// uploadMedia 上传媒体文件
func (a *WechatOfficialAccountAdapter) uploadMedia(ctx context.Context, accessToken, mediaType string, file io.Reader) (string, error) {
	// 微信上传媒体文件API
	// https://api.weixin.qq.com/cgi-bin/media/upload?access_token=TOKEN&type=TYPE

	// TODO: 实际调用微信API
	return "MEDIA_ID_PLACEHOLDER", nil
}
