# API接口文档:多渠道发布框架模块

**文档编号**: DE-API-2025-DOC-028
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考设计**: 《28-企业应用中心_多渠道发布框架.md》

---

## 📋 目录

1. [渠道管理API](#1-渠道管理api)
2. [Bot发布API](#2-bot发布api)
3. [Webhook接收API](#3-webhook接收api)
4. [消息管理API](#4-消息管理api)
5. [渠道统计API](#5-渠道统计api)
6. [渠道配置API](#6-渠道配置api)

---

## 1. 渠道管理API

### 1.1 获取支持的渠道类型列表

**接口描述**: 获取所有支持的渠道类型及配置要求

**请求方式**: `GET /api/v1/channels/types`

**权限要求**: 无需认证

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "type": "wechat_official_account",
      "name": "微信公众号",
      "icon": "https://cdn.coze.com/channels/wechat.png",
      "description": "将Bot发布为微信公众号,支持消息推送、菜单配置",
      "category": "social",
      "features": [
        "text_message",
        "image_message",
        "menu_config",
        "auto_reply"
      ],
      "config_fields": [
        {
          "name": "app_id",
          "label": "AppID",
          "type": "text",
          "required": true,
          "description": "微信公众平台AppID"
        },
        {
          "name": "app_secret",
          "label": "AppSecret",
          "type": "password",
          "required": true,
          "description": "微信公众平台AppSecret"
        },
        {
          "name": "token",
          "label": "Token",
          "type": "text",
          "required": true,
          "description": "消息验证Token"
        },
        {
          "name": "encoding_aes_key",
          "label": "EncodingAESKey",
          "type": "text",
          "required": false,
          "description": "消息加密密钥"
        }
      ],
      "setup_guide_url": "https://docs.coze.com/channels/wechat",
      "pricing": {
        "free": true,
        "limitations": []
      }
    },
    {
      "type": "feishu_app",
      "name": "飞书应用",
      "icon": "https://cdn.coze.com/channels/feishu.png",
      "description": "将Bot发布为飞书应用,支持卡片消息、事件订阅",
      "category": "enterprise",
      "features": [
        "text_message",
        "card_message",
        "event_subscription"
      ],
      "config_fields": [
        {
          "name": "app_id",
          "label": "AppID",
          "type": "text",
          "required": true
        },
        {
          "name": "app_secret",
          "label": "AppSecret",
          "type": "password",
          "required": true
        },
        {
          "name": "encrypt_key",
          "label": "EncryptKey",
          "type": "text",
          "required": false
        },
        {
          "name": "verification_token",
          "label": "VerificationToken",
          "type": "text",
          "required": true
        }
      ],
      "setup_guide_url": "https://docs.coze.com/channels/feishu"
    },
    {
      "type": "discord_bot",
      "name": "Discord Bot",
      "icon": "https://cdn.coze.com/channels/discord.png",
      "description": "将Bot发布为Discord Bot,支持命令处理、Slash Command",
      "category": "community",
      "features": [
        "text_message",
        "slash_command",
        "embed_message"
      ],
      "config_fields": [
        {
          "name": "bot_token",
          "label": "Bot Token",
          "type": "password",
          "required": true
        },
        {
          "name": "client_id",
          "label": "Client ID",
          "type": "text",
          "required": true
        },
        {
          "name": "client_secret",
          "label": "Client Secret",
          "type": "password",
          "required": true
        }
      ],
      "setup_guide_url": "https://docs.coze.com/channels/discord"
    }
  ]
}
```

### 1.2 创建渠道配置

**接口描述**: 为Bot创建渠道配置

**请求方式**: `POST /api/v1/bots/{bot_id}/channels`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**请求体**:

```json
{
  "data": {
    "channel_type": "wechat_official_account",
    "name": "我的微信公众号",
    "enabled": true,
    "config": {
      "app_id": "wx1234567890abcdef",
      "app_secret": "xxxxxxxxxxxxxxxxxx",
      "token": "cozebot_token",
      "encoding_aes_key": "xxxxxxxxxxxxxxxxxx"
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| channel_type | Enum | ✅ | 渠道类型 |
| name | String | ✅ | 配置名称 |
| enabled | Boolean | 否 | 是否启用(默认false) |
| config | Object | ✅ | 渠道配置(根据channel_type不同而不同) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Channel config created",
  "data": {
    "channel_id": 123,
    "bot_id": 789,
    "channel_type": "wechat_official_account",
    "name": "我的微信公众号",
    "enabled": false,
    "status": "DRAFT",
    "config": {
      "app_id": "wx1234567890abcdef",
      "app_secret": "********",
      "token": "cozebot_token",
      "encoding_aes_key": "********"
    },
    "webhook_url": null,
    "published_at": null,
    "created_at": "2025-12-30T10:00:00Z",
    "updated_at": "2025-12-30T10:00:00Z"
  }
}
```

### 1.3 获取Bot的渠道列表

**接口描述**: 获取Bot的所有渠道配置

**请求方式**: `GET /api/v1/bots/{bot_id}/channels`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "channel_id": 123,
      "bot_id": 789,
      "channel_type": "wechat_official_account",
      "name": "我的微信公众号",
      "enabled": true,
      "status": "PUBLISHED",
      "config": {
        "app_id": "wx1234567890abcdef"
      },
      "webhook_url": "https://api.coze.com/api/webhooks/wechat/123",
      "published_at": "2025-12-25T10:00:00Z",
      "last_message_at": "2025-12-30T09:30:00Z",
      "message_count_today": 152,
      "created_at": "2025-12-20T10:00:00Z"
    },
    {
      "channel_id": 124,
      "bot_id": 789,
      "channel_type": "feishu_app",
      "name": "飞书应用",
      "enabled": false,
      "status": "DRAFT",
      "config": {},
      "webhook_url": null,
      "published_at": null,
      "created_at": "2025-12-28T15:00:00Z"
    }
  ]
}
```

### 1.4 更新渠道配置

**接口描述**: 更新渠道配置信息

**请求方式**: `PATCH /api/v1/bots/{bot_id}/channels/{channel_id}`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**请求体**:

```json
{
  "data": {
    "name": "我的微信公众号(更新)",
    "enabled": true,
    "config": {
      "token": "new_token_value"
    }
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Channel config updated",
  "data": {
    "channel_id": 123,
    "name": "我的微信公众号(更新)",
    "updated_at": "2025-12-30T11:00:00Z"
  }
}
```

### 1.5 删除渠道配置

**接口描述**: 删除渠道配置

**请求方式**: `DELETE /api/v1/bots/{bot_id}/channels/{channel_id}`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Channel config deleted",
  "data": {
    "channel_id": 123,
    "deleted_at": "2025-12-30T12:00:00Z"
  }
}
```

---

## 2. Bot发布API

### 2.1 发布Bot到渠道

**接口描述**: 将Bot发布到指定渠道

**请求方式**: `POST /api/v1/bots/{bot_id}/channels/{channel_id}/publish`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**请求体**:

```json
{
  "data": {
    "auto_create_menu": true,
    "menu_config": {
      "button": [
        {
          "type": "click",
          "name": "开始对话",
          "key": "START_CONVERSATION"
        },
        {
          "type": "view",
          "name": "查看官网",
          "url": "https://example.com"
        }
      ]
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| auto_create_menu | Boolean | 否 | 是否自动创建菜单 |
| menu_config | Object | 否 | 菜单配置 |

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot published to channel",
  "data": {
    "channel_id": 123,
    "bot_id": 789,
    "channel_type": "wechat_official_account",
    "status": "PUBLISHED",
    "webhook_url": "https://api.coze.com/api/webhooks/wechat/123",
    "qr_code_url": "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=xxxxx",
    "publish_result": {
      "server_url_configured": true,
      "menu_created": true,
      "access_token": "access_token_value",
      "expires_in": 7200
    },
    "next_steps": [
      {
        "title": "配置服务器地址",
        "description": "在微信公众平台配置服务器URL",
        "url": "https://mp.weixin.qq.com",
        "config_value": "https://api.coze.com/api/webhooks/wechat/123"
      },
      {
        "title": "扫描二维码",
        "description": "使用微信扫描二维码关注公众号"
      }
    ],
    "published_at": "2025-12-30T13:00:00Z"
  }
}
```

### 2.2 取消发布Bot

**接口描述**: 取消Bot的渠道发布

**请求方式**: `POST /api/v1/bots/{bot_id}/channels/{channel_id}/unpublish`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Bot unpublished from channel",
  "data": {
    "channel_id": 123,
    "status": "UNPUBLISHED",
    "unpublished_at": "2025-12-30T14:00:00Z"
  }
}
```

### 2.3 获取发布状态

**接口描述**: 获取Bot在渠道的发布状态

**请求方式**: `GET /api/v1/bots/{bot_id}/channels/{channel_id}/status`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "channel_id": 123,
    "bot_id": 789,
    "channel_type": "wechat_official_account",
    "status": "PUBLISHED",
    "enabled": true,
    "health_check": {
      "status": "HEALTHY",
      "last_check_at": "2025-12-30T14:30:00Z",
      "response_time_ms": 125,
      "error_message": null
    },
    "webhook_configured": true,
    "webhook_url": "https://api.coze.com/api/webhooks/wechat/123",
    "published_at": "2025-12-25T10:00:00Z",
    "last_message_at": "2025-12-30T14:25:00Z",
    "stats": {
      "total_messages": 5230,
      "messages_today": 152,
      "errors_today": 0,
      "avg_response_time_ms": 350
    }
  }
}
```

### 2.4 测试渠道连接

**接口描述**: 测试渠道配置是否正确

**请求方式**: `POST /api/v1/bots/{bot_id}/channels/{channel_id}/test`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Channel connection test passed",
  "data": {
    "channel_id": 123,
    "test_result": "SUCCESS",
    "tests": [
      {
        "name": "config_validation",
        "status": "PASSED",
        "message": "配置参数验证通过"
      },
      {
        "name": "api_access",
        "status": "PASSED",
        "message": "API访问正常",
        "details": {
          "access_token_obtained": true
        }
      },
      {
        "name": "webhook_reachable",
        "status": "PASSED",
        "message": "Webhook地址可访问"
      }
    ],
    "tested_at": "2025-12-30T15:00:00Z"
  }
}
```

---

## 3. Webhook接收API

### 3.1 接收微信Webhook

**接口描述**: 接收微信公众平台的Webhook消息

**请求方式**: `POST /api/webhooks/wechat/{channel_id}`

**权限要求**: 无需认证(签名验证)

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| channel_id | Long | 渠道配置ID |

**请求头**:

```http
Content-Type: application/xml or application/json
X-Wechat-Signature: {signature}
X-Wechat-Timestamp: {timestamp}
X-Wechat-Nonce: {nonce}
```

**请求体**(XML格式):

```xml
<xml>
  <ToUserName><![CDATA[gh_xxxxx]]></ToUserName>
  <FromUserName><![CDATA[oXXXXXXX]]></FromUserName>
  <CreateTime>1234567890</CreateTime>
  <MsgType><![CDATA[text]]></MsgType>
  <Content><![CDATA[你好]]></Content>
  <MsgId>1234567890123456</MsgId>
</xml>
```

**响应**:

```xml
<xml>
  <ReturnCode><![CDATA[0]]></ReturnCode>
  <ReturnMsg><![CDATA[ok]]></ReturnMsg>
</xml>
```

### 3.2 接收飞书Webhook

**接口描述**: 接收飞书应用的事件推送

**请求方式**: `POST /api/webhooks/feishu/{channel_id}`

**权限要求**: 无需认证(签名验证)

**请求体**:

```json
{
  "token": "verify_token",
  "challenge": "challenge_string",
  "type": "url_verification"
}
```

**响应**(URL验证):

```json
{
  "challenge": "challenge_string"
}
```

**请求体**(事件推送):

```json
{
  "schema": "2.0",
  "header": {
    "event_id": "xxxxx",
    "event_type": "im.message.receive_v1",
    "create_time": "1234567890",
    "tenant_key": "xxxxx",
    "app_id": "xxxxx"
  },
  "event": {
    "sender": {
      "sender_id": {
        "open_id": "ou_xxxxx"
      }
    },
    "message": {
      "message_id": "om_xxxxx",
      "content": "{\"text\":\"你好\"}"
    }
  }
}
```

**响应**:

```json
{
  "code": 0,
  "msg": "success"
}
```

### 3.3 接收Discord Webhook

**接口描述**: 接收Discord的交互请求

**请求方式**: `POST /api/webhooks/discord/{channel_id}`

**权限要求**: 无需认证(签名验证)

**请求头**:

```http
Content-Type: application/json
X-Signature-Ed25519: {signature}
X-Signature-Timestamp: {timestamp}
```

**请求体**:

```json
{
  "type": 1,
  "version": 1,
  "application_id": "xxxxx",
  "interaction_token": "xxxxx",
  "token": "xxxxx"
}
```

**响应**(PING):

```json
{
  "type": 1
}
```

---

## 4. 消息管理API

### 4.1 获取渠道消息列表

**接口描述**: 获取渠道的消息记录

**请求方式**: `GET /api/v1/bots/{bot_id}/channels/{channel_id}/messages`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| start_date | String | 否 | 开始日期(ISO8601) |
| end_date | String | 否 | 结束日期(ISO8601) |
| status | String | 否 | 状态筛选(processed/failed/pending) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "message_id": "msg_123456",
      "external_message_id": "wechat_msg_xxxxx",
      "channel_type": "wechat_official_account",
      "direction": "INCOMING",
      "user": {
        "external_user_id": "oXXXXXXX",
        "nickname": "张三"
      },
      "message_type": "TEXT",
      "content": "你好",
      "processed": true,
      "bot_response": {
        "response_id": "resp_123456",
        "content": "你好!我是智能助手,有什么可以帮您?",
        "processed_at": "2025-12-30T14:00:00Z"
      },
      "error": null,
      "created_at": "2025-12-30T14:00:00Z"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_count": 152,
      "total_pages": 8
    }
  }
}
```

### 4.2 获取消息详情

**接口描述**: 获取消息的完整详情

**请求方式**: `GET /api/v1/bots/{bot_id}/channels/{channel_id}/messages/{message_id}`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |
| message_id | String | 消息ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message_id": "msg_123456",
    "external_message_id": "wechat_msg_xxxxx",
    "channel_type": "wechat_official_account",
    "direction": "INCOMING",
    "user": {
      "external_user_id": "oXXXXXXX",
      "nickname": "张三",
      "avatar_url": "http://thirdwx.qlogo.cn/xxx"
    },
    "message_type": "TEXT",
    "content": "你好",
    "raw_message": "<xml>...</xml>",
    "processed": true,
    "processing_time_ms": 350,
    "bot_response": {
      "response_id": "resp_123456",
      "content": "你好!我是智能助手,有什么可以帮您?",
      "raw_response": "{\"msgtype\":\"text\",\"text\":{\"content\":\"你好!...\"}}",
      "sent_at": "2025-12-30T14:00:00Z"
    },
    "error": null,
    "retry_count": 0,
    "created_at": "2025-12-30T14:00:00Z",
    "processed_at": "2025-12-30T14:00:00Z"
  }
}
```

### 4.3 重新发送消息

**接口描述**: 重新发送失败的消息

**请求方式**: `POST /api/v1/bots/{bot_id}/channels/{channel_id}/messages/{message_id}/retry`

**权限要求**: `bot:write`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |
| message_id | String | 消息ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Message retry initiated",
  "data": {
    "message_id": "msg_123456",
    "retry_count": 1,
    "status": "RETRYING",
    "initiated_at": "2025-12-30T16:00:00Z"
  }
}
```

---

## 5. 渠道统计API

### 5.1 获取渠道统计概览

**接口描述**: 获取渠道的统计概览数据

**请求方式**: `GET /api/v1/bots/{bot_id}/channels/{channel_id}/stats`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | String | 否 | 统计周期(today/week/month/all) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "week",
    "message_stats": {
      "total_messages": 1523,
      "incoming_messages": 789,
      "outgoing_messages": 734,
      "avg_response_time_ms": 350,
      "success_rate": 99.2
    },
    "user_stats": {
      "total_users": 235,
      "active_users": 152,
      "new_users": 45
    },
    "error_stats": {
      "total_errors": 12,
      "error_rate": 0.8,
      "common_errors": [
        {
          "error_type": "TIMEOUT",
          "count": 8
        },
        {
          "error_type": "API_ERROR",
          "count": 4
        }
      ]
    },
    "daily_breakdown": [
      {
        "date": "2025-12-24",
        "messages": 215,
        "users": 34
      },
      {
        "date": "2025-12-25",
        "messages": 198,
        "users": 31
      }
    ]
  }
}
```

### 5.2 获取渠道健康检查

**接口描述**: 获取渠道的健康检查结果

**请求方式**: `GET /api/v1/bots/{bot_id}/channels/{channel_id}/health`

**权限要求**: `bot:read`,Bot所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| bot_id | Long | Bot ID |
| channel_id | Long | 渠道配置ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "channel_id": 123,
    "overall_status": "HEALTHY",
    "last_check_at": "2025-12-30T17:00:00Z",
    "checks": [
      {
        "name": "webhook_connectivity",
        "status": "PASS",
        "response_time_ms": 125,
        "last_check_at": "2025-12-30T17:00:00Z"
      },
      {
        "name": "api_access_token",
        "status": "PASS",
        "expires_at": "2025-12-30T18:00:00Z",
        "last_check_at": "2025-12-30T17:00:00Z"
      },
      {
        "name": "message_queue",
        "status": "PASS",
        "queue_size": 23,
        "last_check_at": "2025-12-30T17:00:00Z"
      },
      {
        "name": "error_rate",
        "status": "PASS",
        "error_rate": 0.5,
        "threshold": 5.0,
        "last_check_at": "2025-12-30T17:00:00Z"
      }
    ]
  }
}
```

---

## 6. 渠道配置API

### 6.1 获取渠道配置模板

**接口描述**: 获取指定渠道类型的配置模板

**请求方式**: `GET /api/v1/channels/{channel_type}/config-template`

**权限要求**: 无需认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| channel_type | String | 渠道类型 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "channel_type": "wechat_official_account",
    "channel_name": "微信公众号",
    "config_fields": [
      {
        "name": "app_id",
        "label": "AppID",
        "type": "text",
        "required": true,
        "placeholder": "请输入微信AppID",
        "description": "在微信公众平台-开发-基本配置中获取",
        "validation": {
          "pattern": "^wx[a-zA-Z0-9]{16}$",
          "error_message": "AppID格式不正确"
        }
      },
      {
        "name": "app_secret",
        "label": "AppSecret",
        "type": "password",
        "required": true,
        "description": "在微信公众平台-开发-基本配置中获取"
      },
      {
        "name": "token",
        "label": "Token",
        "type": "text",
        "required": true,
        "description": "自定义Token,用于验证消息",
        "default_value": "cozebot_token",
        "validation": {
          "min_length": 3,
          "max_length": 32
        }
      },
      {
        "name": "encoding_aes_key",
        "label": "消息加密密钥",
        "type": "text",
        "required": false,
        "description": "消息加密方式,可选"
      }
    ],
    "setup_steps": [
      {
        "step": 1,
        "title": "登录微信公众平台",
        "description": "使用公众号管理员账号登录"
      },
      {
        "step": 2,
        "title": "获取AppID和AppSecret",
        "description": "进入开发-基本配置页面"
      },
      {
        "step": 3,
        "title": "配置服务器地址",
        "description": "填写URL和Token"
      }
    ],
    "help_links": [
      {
        "title": "微信公众平台开发文档",
        "url": "https://developers.weixin.qq.com/doc/"
      },
      {
        "title": "Coze多渠道发布指南",
        "url": "https://docs.coze.com/channels/wechat"
      }
    ]
  }
}
```

### 6.2 验证渠道配置

**接口描述**: 验证渠道配置是否正确

**请求方式**: `POST /api/v1/channels/{channel_type}/validate-config`

**权限要求**: 需要认证

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| channel_type | String | 渠道类型 |

**请求体**:

```json
{
  "data": {
    "app_id": "wx1234567890abcdef",
    "app_secret": "xxxxxxxxxxxxxxxxxx",
    "token": "cozebot_token"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Config validation passed",
  "data": {
    "valid": true,
    "warnings": [],
    "test_results": [
      {
        "name": "app_id_format",
        "status": "PASS",
        "message": "AppID格式正确"
      },
      {
        "name": "api_access",
        "status": "PASS",
        "message": "API访问正常,Access Token获取成功",
        "details": {
          "access_token": "access_token_value",
          "expires_in": 7200
        }
      }
    ]
  }
}
```

---

## 附录

### A. 渠道类型枚举

```typescript
enum ChannelType {
  WECHAT_OFFICIAL_ACCOUNT = 'wechat_official_account',  // 微信公众号
  WECHAT_WORK = 'wechat_work',                         // 企业微信
  FEISHU_APP = 'feishu_app',                           // 飞书应用
  DING_TALK = 'ding_talk',                             // 钉钉
  DISCORD_BOT = 'discord_bot',                         // Discord Bot
  SLACK_BOT = 'slack_bot',                             // Slack Bot
  TELEGRAM_BOT = 'telegram_bot',                       // Telegram Bot
}

enum ChannelStatus {
  DRAFT = 'DRAFT',           // 草稿
  PUBLISHED = 'PUBLISHED',   // 已发布
  UNPUBLISHED = 'UNPUBLISHED', // 已取消发布
  ERROR = 'ERROR'            // 错误
}

enum MessageDirection {
  INCOMING = 'INCOMING',     // 接收
  OUTGOING = 'OUTGOING'      // 发送
}

enum MessageType {
  TEXT = 'TEXT',
  IMAGE = 'IMAGE',
  AUDIO = 'AUDIO',
  VIDEO = 'VIDEO',
  FILE = 'FILE',
  CARD = 'CARD',
  EVENT = 'EVENT'
}
```

### B. Go后端实现示例

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// ChannelHandler 渠道处理器
type ChannelHandler struct {
    channelService  ChannelService
    publishService  PublishService
    webhookService  WebhookService
}

// CreateChannelConfig 创建渠道配置
func (h *ChannelHandler) CreateChannelConfig(ctx context.Context, c *app.RequestContext) {
    botID := c.Param("bot_id")
    var req CreateChannelRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    userID := getUserIDFromContext(ctx)

    config, err := h.channelService.CreateConfig(ctx, userID, botID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(201, BotResponse{
        Code:    0,
        Message: "Channel config created",
        Data:    config,
    })
}

// PublishBot 发布Bot到渠道
func (h *ChannelHandler) PublishBot(ctx context.Context, c *app.RequestContext) {
    botID := c.Param("bot_id")
    channelID := c.Param("channel_id")
    var req PublishBotRequest
    c.BindAndValidate(&req)

    userID := getUserIDFromContext(ctx)

    result, err := h.publishService.Publish(ctx, userID, botID, channelID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Bot published to channel",
        Data:    result,
    })
}

// HandleWechatWebhook 处理微信Webhook
func (h *ChannelHandler) HandleWechatWebhook(ctx context.Context, c *app.RequestContext) {
    channelID := c.Param("channel_id")

    // 读取原始请求体
    body := c.GetRawData()

    // 验证签名
    signature := c.GetHeader("X-Wechat-Signature")
    timestamp := c.GetHeader("X-Wechat-Timestamp")
    nonce := c.GetHeader("X-Wechat-Nonce")

    err := h.webhookService.VerifyWechatSignature(ctx, channelID, signature, timestamp, nonce, body)
    if err != nil {
        c.XML(401, BuildWechatErrorResponse("签名验证失败"))
        return
    }

    // 异步处理消息
    err = h.webhookService.HandleWebhookAsync(ctx, "wechat", channelID, body)
    if err != nil {
        c.XML(500, BuildWechatErrorResponse("消息处理失败"))
        return
    }

    // 快速返回成功
    c.XML(200, BuildWechatSuccessResponse())
}
```

### C. 适配器接口定义

```go
// ChannelAdapter 渠道适配器接口
type ChannelAdapter interface {
    // 获取适配器类型
    GetType() ChannelType

    // 验证配置
    ValidateConfig(ctx context.Context, config map[string]interface{}) error

    // 发布Bot
    PublishBot(ctx context.Context, bot *Bot, config map[string]interface{}) (*PublishResult, error)

    // 取消发布Bot
    UnpublishBot(ctx context.Context, channelID string) error

    // 验证Webhook签名
    VerifyWebhook(signature, timestamp, nonce string, body []byte) error

    // 转换接收的消息
    TransformIncomingMessage(rawMessage []byte) (*ChannelMessage, error)

    // 转换发送的消息
    TransformOutgoingMessage(botResponse *BotResponse) (*ChannelOutgoingMessage, error)

    // 发送消息到渠道
    SendMessage(ctx context.Context, message *ChannelOutgoingMessage) error

    // 健康检查
    HealthCheck(ctx context.Context) (*HealthCheckResult, error)
}

// WechatOfficialAccountAdapter 微信公众号适配器
type WechatOfficialAccountAdapter struct {
    httpClient *http.Client
    config     *WechatConfig
}

func (a *WechatOfficialAccountAdapter) GetType() ChannelType {
    return WECHAT_OFFICIAL_ACCOUNT
}

func (a *WechatOfficialAccountAdapter) ValidateConfig(ctx context.Context, config map[string]interface{}) error {
    appID, ok := config["app_id"].(string)
    if !ok || appID == "" {
        return errors.New("app_id is required")
    }

    appSecret, ok := config["app_secret"].(string)
    if !ok || appSecret == "" {
        return errors.New("app_secret is required")
    }

    // 测试API访问
    return a.testAPIAccess(ctx, appID, appSecret)
}
```

### D. 前端TypeScript类型定义

```typescript
// types/channel.ts

interface ChannelConfig {
  channel_id: number;
  bot_id: number;
  channel_type: ChannelType;
  name: string;
  enabled: boolean;
  status: ChannelStatus;
  config: Record<string, any>;
  webhook_url?: string;
  published_at?: string;
  created_at: string;
  updated_at: string;
}

interface PublishResult {
  channel_id: number;
  bot_id: number;
  channel_type: ChannelType;
  status: ChannelStatus;
  webhook_url: string;
  qr_code_url?: string;
  publish_result: {
    server_url_configured: boolean;
    menu_created: boolean;
    access_token?: string;
  };
  next_steps: Array<{
    title: string;
    description: string;
    url?: string;
    config_value?: string;
  }>;
  published_at: string;
}

interface ChannelMessage {
  message_id: string;
  external_message_id: string;
  channel_type: ChannelType;
  direction: 'INCOMING' | 'OUTGOING';
  user: {
    external_user_id: string;
    nickname?: string;
    avatar_url?: string;
  };
  message_type: MessageType;
  content: string;
  processed: boolean;
  bot_response?: {
    response_id: string;
    content: string;
    processed_at: string;
  };
  error?: string;
  created_at: string;
}

// API客户端
class ChannelApiClient {
  async createChannelConfig(
    botId: number,
    data: CreateChannelRequest
  ): Promise<ChannelConfig> {
    const response = await apiClient.post(`/bots/${botId}/channels`, { data });
    return response.data.data;
  }

  async getChannels(botId: number): Promise<ChannelConfig[]> {
    const response = await apiClient.get(`/bots/${botId}/channels`);
    return response.data.data;
  }

  async publishBot(
    botId: number,
    channelId: number,
    data?: PublishBotRequest
  ): Promise<PublishResult> {
    const response = await apiClient.post(
      `/bots/${botId}/channels/${channelId}/publish`,
      { data }
    );
    return response.data.data;
  }

  async getChannelMessages(
    botId: number,
    channelId: number,
    params?: GetMessagesParams
  ): Promise<ChannelMessage[]> {
    const response = await apiClient.get(
      `/bots/${botId}/channels/${channelId}/messages`,
      { params }
    );
    return response.data.data;
  }

  async getChannelStats(
    botId: number,
    channelId: number,
    period?: string
  ): Promise<ChannelStats> {
    const response = await apiClient.get(
      `/bots/${botId}/channels/${channelId}/stats`,
      { params: { period } }
    );
    return response.data.data;
  }
}

export const channelApi = new ChannelApiClient();
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含渠道管理、Bot发布、Webhook接收、消息管理、渠道统计、渠道配置等6个模块的完整API接口定义
