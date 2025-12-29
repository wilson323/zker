# API接口文档:高级多模态功能模块

**文档编号**: DE-API-2025-DOC-031
**版本**: v1.0.0
**创建日期**: 2025-12-30
**参考设计**: 《31-高级多模态功能_语音与视频.md》

---

## 📋 目录

1. [语音识别API (STT)](#1-语音识别api-stt)
2. [语音合成API (TTS)](#2-语音合成api-tts)
3. [实时语音对话API](#3-实时语音对话api)
4. [视频理解API](#4-视频理解api)
5. [多模态编排API](#5-多模态编排api)
6. [语音记录管理API](#6-语音记录管理api)

---

## 1. 语音识别API (STT)

### 1.1 创建语音识别任务

**接口描述**: 上传音频文件并进行语音识别

**请求方式**: `POST /api/v1/speech/stt`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "audio_url": "https://cdn.coze.com/audio/user_voice_20251230.mp3",
    "language": "zh-CN",
    "model": "whisper-1",
    "options": {
      "enable_punctuation": true,
      "enable_timestamps": false,
      "vad_filter": true,
      "temperature": 0.0
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| audio_url | String | ✅ | 音频文件URL(MP3/WAV/M4A) |
| language | String | ✅ | 语言代码(zh-CN/en-US/ja-JP等) |
| model | String | ✅ | STT模型(whisper-1/whisper-large) |
| options | Object | 否 | 识别选项 |
| options.enable_punctuation | Boolean | 否 | 启用标点符号(默认true) |
| options.enable_timestamps | Boolean | 否 | 启用时间戳(默认false) |
| options.vad_filter | Boolean | 否 | 启用VAD过滤静音(默认true) |
| options.temperature | Float | 否 | 采样温度(0.0-1.0,默认0.0) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Speech recognition completed",
  "data": {
    "record_id": "stt_20251230123456",
    "audio_url": "https://cdn.coze.com/audio/user_voice_20251230.mp3",
    "language": "zh-CN",
    "model": "whisper-1",
    "text": "我想查询订单状态",
    "confidence": 0.98,
    "duration_ms": 3500,
    "segments": [
      {
        "text": "我想查询订单状态",
        "start_time": 0.0,
        "end_time": 3.5,
        "confidence": 0.98
      }
    ],
    "processing_time_ms": 1250,
    "created_at": "2025-12-30T10:00:00Z"
  }
}
```

### 1.2 流式语音识别

**接口描述**: 实时流式语音识别(WebSocket)

**WebSocket URL**: `wss://api.coze.com/api/v1/speech/stt/stream`

**连接参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | String | ✅ | 认证Token |
| language | String | ✅ | 语言代码 |
| model | String | ✅ | STT模型 |

**客户端发送消息**(音频数据):

```json
{
  "type": "audio",
  "data": "base64_encoded_audio_data",
  "sequence": 1,
  "is_final": false
}
```

**服务端推送消息**(识别结果):

```json
{
  "type": "recognition",
  "text": "我想查询",
  "is_final": false,
  "confidence": 0.95,
  "sequence": 1
}
```

**最终识别结果**:

```json
{
  "type": "recognition",
  "text": "我想查询订单状态",
  "is_final": true,
  "confidence": 0.98,
  "sequence": 5,
  "record_id": "stt_stream_20251230123456"
}
```

### 1.3 获取识别记录

**接口描述**: 获取语音识别记录详情

**请求方式**: `GET /api/v1/speech/stt/records/{record_id}`

**权限要求**: `speech:read`,记录所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| record_id | String | 记录ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "record_id": "stt_20251230123456",
    "bot_id": 789,
    "conversation_id": 456,
    "audio_url": "https://cdn.coze.com/audio/user_voice_20251230.mp3",
    "language": "zh-CN",
    "model": "whisper-1",
    "text": "我想查询订单状态",
    "confidence": 0.98,
    "duration_ms": 3500,
    "processing_time_ms": 1250,
    "segments": [
      {
        "text": "我想查询订单状态",
        "start_time": 0.0,
        "end_time": 3.5,
        "confidence": 0.98
      }
    ],
    "created_at": "2025-12-30T10:00:00Z"
  }
}
```

---

## 2. 语音合成API (TTS)

### 2.1 创建语音合成任务

**接口描述**: 将文本转换为语音

**请求方式**: `POST /api/v1/speech/tts`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "text": "您好,我是智能客服助手,请问有什么可以帮您?",
    "voice": "xiaoxiao",
    "language": "zh-CN",
    "model": "azure-neural",
    "options": {
      "speed": 1.0,
      "pitch": 0.0,
      "volume": 1.0,
      "emotion": "cheerful",
      "enable_ssml": false
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| text | String | ✅ | 要合成的文本(最多5000字符) |
| voice | String | ✅ | 音色名称(xiaoxiao/yunyang/jenny等) |
| language | String | ✅ | 语言代码 |
| model | String | ✅ | TTS模型(azure-neural/openai-tts) |
| options | Object | 否 | 合成选项 |
| options.speed | Float | 否 | 语速(0.5-2.0,默认1.0) |
| options.pitch | Float | 否 | 音调(-20到20,默认0) |
| options.volume | Float | 否 | 音量(0.0-2.0,默认1.0) |
| options.emotion | String | 否 | 情感(cheerful/neutral/sad) |
| options.enable_ssml | Boolean | 否 | 是否启用SSML(默认false) |

**响应示例**:

```json
{
  "code": 0,
  "message": "Speech synthesis completed",
  "data": {
    "record_id": "tts_20251230654321",
    "text": "您好,我是智能客服助手,请问有什么可以帮您?",
    "voice": "xiaoxiao",
    "language": "zh-CN",
    "audio_url": "https://cdn.coze.com/audio/tts_20251230654321.mp3",
    "duration_ms": 5200,
    "size_bytes": 83520,
    "format": "mp3",
    "sample_rate": 16000,
    "processing_time_ms": 850,
    "from_cache": false,
    "created_at": "2025-12-30T11:00:00Z"
  }
}
```

### 2.2 使用SSML合成语音

**接口描述**: 使用SSML标记进行高级语音合成

**请求方式**: `POST /api/v1/speech/tts/ssml`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "ssml": "<speak version='1.0' xmlns='http://www.w3.org/2001/10/synthesis' xml:lang='zh-CN'><voice name='xiaoxiao'><prosody rate='1.0' pitch='0%'>您好</prosody><break time='500ms'/><emphasis level='strong'>我是智能客服助手</emphasis><break time='300ms'/>请问有什么可以帮您?</voice></speak>",
    "voice": "xiaoxiao",
    "model": "azure-neural"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "SSML synthesis completed",
  "data": {
    "record_id": "tts_ssml_20251230654322",
    "ssml": "<speak>...</speak>",
    "audio_url": "https://cdn.coze.com/audio/tts_ssml_20251230654322.mp3",
    "duration_ms": 6800,
    "created_at": "2025-12-30T11:05:00Z"
  }
}
```

### 2.3 获取可用音色列表

**接口描述**: 获取所有可用的TTS音色

**请求方式**: `GET /api/v1/speech/tts/voices`

**权限要求**: 无需认证

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| language | String | 否 | 语言筛选 |
| gender | String | 否 | 性别筛选(male/female) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "xiaoxiao",
      "name": "晓晓",
      "display_name": "晓晓(女声)",
      "language": "zh-CN",
      "gender": "female",
      "age": "young",
      "style": ["cheerful", "neutral", "sad"],
      "sample_url": "https://cdn.coze.com/audio/voices/xiaoxiao_sample.mp3",
      "is_neural": true,
      "recommended": true
    },
    {
      "id": "yunyang",
      "name": "云扬",
      "display_name": "云扬(男声)",
      "language": "zh-CN",
      "gender": "male",
      "age": "adult",
      "style": ["neutral", "professional"],
      "sample_url": "https://cdn.coze.com/audio/voices/yunyang_sample.mp3",
      "is_neural": true,
      "recommended": false
    },
    {
      "id": "jenny",
      "name": "Jenny",
      "display_name": "Jenny (English Female)",
      "language": "en-US",
      "gender": "female",
      "age": "young",
      "style": ["cheerful", "neutral"],
      "sample_url": "https://cdn.coze.com/audio/voices/jenny_sample.mp3",
      "is_neural": true,
      "recommended": true
    }
  ]
}
```

---

## 3. 实时语音对话API

### 3.1 建立实时语音对话连接

**接口描述**: 建立WebSocket连接进行实时语音对话

**WebSocket URL**: `wss://api.coze.com/api/v1/speech/conversation`

**连接参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| token | String | ✅ | 认证Token |
| bot_id | Long | ✅ | Bot ID |
| voice | String | ✅ | TTS音色 |
| language | String | ✅ | 语言代码 |

**连接建立后,服务端发送初始化消息**:

```json
{
  "type": "init",
  "conversation_id": "voice_conv_20251230123456",
  "session_id": "sess_xxxxx",
  "server_time": "2025-12-30T12:00:00Z"
}
```

### 3.2 实时对话消息协议

**客户端发送音频数据**:

```json
{
  "type": "audio",
  "data": "base64_encoded_audio_chunk",
  "sequence": 1,
  "sample_rate": 16000,
  "encoding": "pcm"
}
```

**客户端标记语音结束**:

```json
{
  "type": "audio_end",
  "sequence": 25
}
```

**服务端推送中间识别结果**:

```json
{
  "type": "recognition_interim",
  "text": "我想查询订单",
  "is_final": false,
  "confidence": 0.92,
  "sequence": 20
}
```

**服务端推送最终识别结果**:

```json
{
  "type": "recognition_final",
  "text": "我想查询订单状态",
  "is_final": true,
  "confidence": 0.97,
  "sequence": 25
}
```

**服务端推送Bot回复(文本)**:

```json
{
  "type": "bot_response_text",
  "text": "您的订单已发货,预计明天送达",
  "message_id": "msg_123456"
}
```

**服务端推送Bot回复(音频)**:

```json
{
  "type": "bot_response_audio",
  "audio_url": "https://cdn.coze.com/audio/bot_resp_20251230.mp3",
  "text": "您的订单已发货,预计明天送达",
  "duration_ms": 4800,
  "message_id": "msg_123456"
}
```

**服务端推送错误**:

```json
{
  "type": "error",
  "error": {
    "code": "STT_ENGINE_ERROR",
    "message": "语音识别失败,请重试",
    "retryable": true
  }
}
```

### 3.3 关闭对话连接

**客户端主动关闭**:

```json
{
  "type": "close",
  "reason": "user_ended_conversation"
}
```

**服务端关闭连接**:

```json
{
  "type": "close",
  "reason": "server_timeout",
  "duration_seconds": 300
}
```

---

## 4. 视频理解API

### 4.1 创建视频理解任务

**接口描述**: 上传视频并进行分析理解

**请求方式**: `POST /api/v1/video/understand`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "video_url": "https://cdn.coze.com/video/cooking_tutorial.mp4",
    "tasks": [
      "scene_recognition",
      "action_detection",
      "text_extraction",
      "object_detection"
    ],
    "options": {
      "frame_interval": 5,
      "max_frames": 20,
      "enable_audio_analysis": true,
      "language": "zh-CN"
    }
  }
}
```

**请求参数说明**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| video_url | String | ✅ | 视频文件URL |
| tasks | Array | ✅ | 分析任务列表 |
| options | Object | 否 | 分析选项 |
| options.frame_interval | Integer | 否 | 关键帧提取间隔(秒,默认5) |
| options.max_frames | Integer | 否 | 最多提取帧数(默认20) |
| options.enable_audio_analysis | Boolean | 否 | 是否分析音频(默认true) |
| options.language | String | 否 | 文本识别语言 |

**任务类型**:

| 任务 | 说明 |
|------|------|
| scene_recognition | 场景识别(室内/室外/餐厅等) |
| action_detection | 动作检测(烹饪/跑步/演讲等) |
| text_extraction | 文字提取(OCR) |
| object_detection | 物体检测(人物/车辆/物品等) |
| emotion_analysis | 情感分析(人物情绪识别) |

**响应示例**(任务创建):

```json
{
  "code": 0,
  "message": "Video understanding task created",
  "data": {
    "task_id": "video_task_20251230123456",
    "video_url": "https://cdn.coze.com/video/cooking_tutorial.mp4",
    "status": "PROCESSING",
    "estimated_time_seconds": 30,
    "created_at": "2025-12-30T13:00:00Z"
  }
}
```

### 4.2 获取视频理解结果

**接口描述**: 获取视频理解任务的结果

**请求方式**: `GET /api/v1/video/understand/{task_id}`

**权限要求**: `video:read`,任务创建者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| task_id | String | 任务ID |

**响应示例**(处理完成):

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "video_task_20251230123456",
    "video_url": "https://cdn.coze.com/video/cooking_tutorial.mp4",
    "status": "COMPLETED",
    "duration_seconds": 180,
    "processed_frames": 36,
    "results": {
      "scene_recognition": {
        "primary_scene": "kitchen",
        "confidence": 0.95,
        "scenes": [
          {
            "scene": "kitchen",
            "confidence": 0.95,
            "time_range": {
              "start": 0,
              "end": 180
            }
          }
        ]
      },
      "action_detection": {
        "actions": [
          {
            "action": "cooking",
            "confidence": 0.92,
            "time_range": {
              "start": 10,
              "end": 160
            },
            "objects": ["pan", "knife", "ingredients"]
          },
          {
            "action": "mixing",
            "confidence": 0.88,
            "time_range": {
              "start": 90,
              "end": 120
            }
          }
        ]
      },
      "text_extraction": {
        "text_regions": [
          {
            "text": "意大利面制作教程",
            "confidence": 0.98,
            "bbox": [100, 50, 500, 100],
            "timestamp": 0.5
          }
        ]
      },
      "object_detection": {
        "objects": [
          {
            "object": "person",
            "count": 1,
            "confidence": 0.99,
            "appears_in_frames": [0, 5, 10, 15, 20]
          },
          {
            "object": "pan",
            "count": 1,
            "confidence": 0.96,
            "appears_in_frames": [5, 10, 15]
          }
        ]
      },
      "audio_analysis": {
        "speech_detected": true,
        "transcript": "大家好,今天我来教大家制作意大利面",
        "language": "zh-CN",
        "music_detected": false
      }
    },
    "summary": "这是一个烹饪教学视频,展示如何在厨房制作意大利面。视频中包含人物讲解、烹饪和搅拌等动作,时长3分钟。",
    "created_at": "2025-12-30T13:00:00Z",
    "completed_at": "2025-12-30T13:00:28Z"
  }
}
```

### 4.3 提取视频关键帧

**接口描述**: 提取视频的关键帧图片

**请求方式**: `POST /api/v1/video/extract-frames`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "video_url": "https://cdn.coze.com/video/cooking_tutorial.mp4",
    "interval_seconds": 10,
    "max_frames": 10,
    "output_format": "jpg"
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Frames extracted successfully",
  "data": {
    "task_id": "frame_extract_20251230123456",
    "video_url": "https://cdn.coze.com/video/cooking_tutorial.mp4",
    "frames": [
      {
        "frame_number": 0,
        "timestamp_seconds": 0.0,
        "image_url": "https://cdn.coze.com/video/frames/frame_0.jpg",
        "width": 1920,
        "height": 1080
      },
      {
        "frame_number": 1,
        "timestamp_seconds": 10.0,
        "image_url": "https://cdn.coze.com/video/frames/frame_1.jpg",
        "width": 1920,
        "height": 1080
      }
    ],
    "total_frames": 10,
    "created_at": "2025-12-30T13:05:00Z"
  }
}
```

---

## 5. 多模态编排API

### 5.1 创建多模态编排任务

**接口描述**: 创建包含语音、视频、文本的多模态处理任务

**请求方式**: `POST /api/v1/multimodal/orchestration`

**权限要求**: 需要认证

**请求体**:

```json
{
  "data": {
    "bot_id": 789,
    "name": "多模态客服助手",
    "workflow": {
      "steps": [
        {
          "id": "step_1",
          "type": "stt",
          "name": "语音识别",
          "config": {
            "language": "zh-CN",
            "model": "whisper-1"
          },
          "input_from": "user_audio"
        },
        {
          "id": "step_2",
          "type": "bot_process",
          "name": "Bot处理",
          "config": {
            "bot_id": 789
          },
          "input_from": "step_1"
        },
        {
          "id": "step_3",
          "type": "video_understand",
          "name": "视频理解",
          "config": {
            "tasks": ["scene_recognition", "object_detection"]
          },
          "input_from": "user_video",
          "condition": "user_provided_video"
        },
        {
          "id": "step_4",
          "type": "tts",
          "name": "语音合成",
          "config": {
            "voice": "xiaoxiao",
            "language": "zh-CN"
          },
          "input_from": "step_2"
        }
      ]
    },
    "inputs": {
      "user_audio": "https://cdn.coze.com/audio/user_query.mp3",
      "user_video": null
    }
  }
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "Multimodal orchestration task created",
  "data": {
    "task_id": "mm_task_20251230123456",
    "bot_id": 789,
    "name": "多模态客服助手",
    "status": "PROCESSING",
    "current_step": "step_1",
    "progress": 25,
    "created_at": "2025-12-30T14:00:00Z"
  }
}
```

### 5.2 获取编排任务状态

**接口描述**: 获取多模态编排任务的执行状态

**请求方式**: `GET /api/v1/multimodal/orchestration/{task_id}`

**权限要求**: `multimodal:read`,任务创建者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| task_id | String | 任务ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "task_id": "mm_task_20251230123456",
    "bot_id": 789,
    "name": "多模态客服助手",
    "status": "COMPLETED",
    "steps": [
      {
        "id": "step_1",
        "type": "stt",
        "name": "语音识别",
        "status": "COMPLETED",
        "result": {
          "text": "帮我看看这个视频",
          "confidence": 0.98
        },
        "started_at": "2025-12-30T14:00:00Z",
        "completed_at": "2025-12-30T14:00:01Z"
      },
      {
        "id": "step_2",
        "type": "bot_process",
        "name": "Bot处理",
        "status": "COMPLETED",
        "result": {
          "response": "我来帮您分析这个视频",
          "message_id": "msg_123456"
        },
        "started_at": "2025-12-30T14:00:01Z",
        "completed_at": "2025-12-30T14:00:03Z"
      },
      {
        "id": "step_3",
        "type": "video_understand",
        "name": "视频理解",
        "status": "SKIPPED",
        "skip_reason": "condition_not_met: user_provided_video"
      },
      {
        "id": "step_4",
        "type": "tts",
        "name": "语音合成",
        "status": "COMPLETED",
        "result": {
          "audio_url": "https://cdn.coze.com/audio/tts_20251230.mp3",
          "duration_ms": 4500
        },
        "started_at": "2025-12-30T14:00:03Z",
        "completed_at": "2025-12-30T14:00:04Z"
      }
    ],
    "final_output": {
      "text": "我来帮您分析这个视频",
      "audio_url": "https://cdn.coze.com/audio/tts_20251230.mp3"
    },
    "progress": 100,
    "created_at": "2025-12-30T14:00:00Z",
    "completed_at": "2025-12-30T14:00:04Z"
  }
}
```

---

## 6. 语音记录管理API

### 6.1 获取语音记录列表

**接口描述**: 获取用户的语音记录(STT/TTS)

**请求方式**: `GET /api/v1/speech/records`

**权限要求**: `speech:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | Integer | 否 | 页码(默认1) |
| page_size | Integer | 否 | 每页数量(默认20) |
| speech_type | String | 否 | 类型筛选(STT/TTS) |
| start_date | String | 否 | 开始日期 |
| end_date | String | 否 | 结束日期 |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "record_id": "stt_20251230123456",
      "speech_type": "STT",
      "bot_id": 789,
      "conversation_id": 456,
      "audio_url": "https://cdn.coze.com/audio/user_voice.mp3",
      "text": "我想查询订单状态",
      "confidence": 0.98,
      "language": "zh-CN",
      "duration_ms": 3500,
      "created_at": "2025-12-30T10:00:00Z"
    },
    {
      "record_id": "tts_20251230654321",
      "speech_type": "TTS",
      "bot_id": 789,
      "conversation_id": 456,
      "text": "您的订单已发货",
      "voice": "xiaoxiao",
      "audio_url": "https://cdn.coze.com/audio/tts_20251230.mp3",
      "duration_ms": 4200,
      "language": "zh-CN",
      "created_at": "2025-12-30T10:00:05Z"
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

### 6.2 删除语音记录

**接口描述**: 删除指定的语音记录

**请求方式**: `DELETE /api/v1/speech/records/{record_id}`

**权限要求**: `speech:write`,记录所有者

**路径参数**:

| 参数 | 类型 | 说明 |
|------|------|------|
| record_id | String | 记录ID |

**响应示例**:

```json
{
  "code": 0,
  "message": "Record deleted",
  "data": {
    "record_id": "stt_20251230123456",
    "deleted_at": "2025-12-30T15:00:00Z"
  }
}
```

### 6.3 获取语音使用统计

**接口描述**: 获取语音功能的使用统计

**请求方式**: `GET /api/v1/speech/stats`

**权限要求**: `speech:read`

**请求参数**:

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| period | String | 否 | 统计周期(today/week/month) |

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "week",
    "stt_stats": {
      "total_records": 152,
      "total_duration_seconds": 532,
      "avg_confidence": 0.96,
      "total_characters": 8950,
      "languages": [
        {
          "language": "zh-CN",
          "count": 120,
          "percentage": 78.9
        },
        {
          "language": "en-US",
          "count": 32,
          "percentage": 21.1
        }
      ]
    },
    "tts_stats": {
      "total_records": 145,
      "total_duration_seconds": 635,
      "total_characters": 12500,
      "top_voices": [
        {
          "voice": "xiaoxiao",
          "count": 85,
          "percentage": 58.6
        },
        {
          "voice": "yunyang",
          "count": 40,
          "percentage": 27.6
        }
      ],
      "cache_hit_rate": 0.35
    },
    "cost_estimates": {
      "stt_cost": 0.53,
      "tts_cost": 0.63,
      "total_cost": 1.16,
      "currency": "USD"
    },
    "daily_breakdown": [
      {
        "date": "2025-12-24",
        "stt_count": 20,
        "tts_count": 18
      }
    ]
  }
}
```

---

## 附录

### A. 错误码参考

| 错误码 | 说明 |
|--------|------|
| 40001 | 请求参数错误 |
| 40002 | 音频文件格式不支持 |
| 40003 | 视频文件格式不支持 |
| 40004 | 文本长度超过限制 |
| 40101 | Token缺失或无效 |
| 40301 | 无权限访问记录 |
| 40401 | 记录不存在 |
| 40901 | 重复的任务请求 |
| 42201 | 音频质量太低 |
| 42202 | 视频处理失败 |
| 42203 | STT识别失败(置信度过低) |
| 50001 | STT引擎错误 |
| 50002 | TTS引擎错误 |
| 50003 | 视频处理引擎错误 |

### B. WebSocket错误码

| 错误码 | 说明 | 可重试 |
|--------|------|--------|
| 4001 | Token无效 | false |
| 4002 | Bot不存在 | false |
| 4003 | 音频格式错误 | false |
| 4004 | 采样率不支持 | false |
| 5001 | STT引擎错误 | true |
| 5002 | TTS引擎错误 | true |
| 5003 | Bot处理错误 | true |
| 5004 | 内部服务错误 | true |

### C. Go后端实现示例

```go
package handler

import (
    "context"
    "github.com/cloudwego/hertz/pkg/app"
)

// SpeechHandler 语音处理器
type SpeechHandler struct {
    speechService  SpeechService
    sttEngine      STTEngine
    ttsEngine      TTSEngine
    videoService   VideoService
}

// CreateSTTTask 创建语音识别任务
func (h *SpeechHandler) CreateSTTTask(ctx context.Context, c *app.RequestContext) {
    var req CreateSTTRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    userID := getUserIDFromContext(ctx)

    result, err := h.speechService.RecognizeSpeech(ctx, userID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Speech recognition completed",
        Data:    result,
    })
}

// CreateTTSTask 创建语音合成任务
func (h *SpeechHandler) CreateTTSTask(ctx context.Context, c *app.RequestContext) {
    var req CreateTTSRequest
    if err := c.BindAndValidate(&req); err != nil {
        c.JSON(400, NewBusinessError(40001, "Invalid request parameters"))
        return
    }

    userID := getUserIDFromContext(ctx)

    // 检查缓存
    cacheKey := fmt.Sprintf("tts:%s:%s", req.Text, req.Voice)
    if cachedAudio := h.cache.Get(ctx, cacheKey); cachedAudio != nil {
        c.JSON(200, BotResponse{
            Code:    0,
            Message: "Speech synthesis completed (cached)",
            Data: map[string]interface{}{
                "audio_url":    cachedAudio.URL,
                "from_cache":   true,
                "record_id":    cachedAudio.RecordID,
            },
        })
        return
    }

    result, err := h.speechService.SynthesizeSpeech(ctx, userID, &req)
    if err != nil {
        handleBusinessError(c, err)
        return
    }

    // 写入缓存
    h.cache.Set(ctx, cacheKey, result, 7*24*time.Hour)

    c.JSON(200, BotResponse{
        Code:    0,
        Message: "Speech synthesis completed",
        Data:    result,
    })
}

// HandleWebSocketConversation 处理实时语音对话WebSocket
func (h *SpeechHandler) HandleWebSocketConversation(ctx context.Context, c *app.RequestContext) {
    // 升级为WebSocket连接
    err := upgrader.Upgrade(c)
    if err != nil {
        c.String(400, "WebSocket upgrade failed")
        return
    }

    conn := getWebSocketConnection(c)
    defer conn.Close()

    userID := getUserIDFromToken(c.Query("token"))
    botID := c.Query("bot_id")
    voice := c.Query("voice")
    language := c.Query("language")

    // 创建会话
    session := h.speechService.CreateConversationSession(ctx, userID, botID, voice, language)

    // 发送初始化消息
    conn.WriteJSON(map[string]interface{}{
        "type": "init",
        "conversation_id": session.ID,
        "session_id": session.SessionID,
        "server_time": time.Now().Format(time.RFC3339),
    })

    // 处理消息循环
    for {
        var msg WebSocketMessage
        err := conn.ReadJSON(&msg)
        if err != nil {
            break
        }

        switch msg.Type {
        case "audio":
            // 处理音频数据
            h.handleAudioData(ctx, session, msg)
        case "audio_end":
            // 处理音频结束
            h.handleAudioEnd(ctx, session, msg)
        case "close":
            // 关闭连接
            break
        }
    }
}
```

### D. 前端TypeScript类型定义

```typescript
// types/multimodal.ts

// STT识别结果
interface STTResult {
  record_id: string;
  text: string;
  confidence: number;
  duration_ms: number;
  language: string;
  segments?: Array<{
    text: string;
    start_time: number;
    end_time: number;
    confidence: number;
  }>;
}

// TTS合成结果
interface TTSResult {
  record_id: string;
  text: string;
  audio_url: string;
  duration_ms: number;
  voice: string;
  language: string;
  from_cache: boolean;
}

// 视频理解结果
interface VideoUnderstandingResult {
  task_id: string;
  status: 'PROCESSING' | 'COMPLETED' | 'FAILED';
  results?: {
    scene_recognition?: {
      primary_scene: string;
      confidence: number;
      scenes: Array<{
        scene: string;
        confidence: number;
        time_range: { start: number; end: number };
      }>;
    };
    action_detection?: {
      actions: Array<{
        action: string;
        confidence: number;
        time_range: { start: number; end: number };
        objects: string[];
      }>;
    };
    object_detection?: {
      objects: Array<{
        object: string;
        count: number;
        confidence: number;
        appears_in_frames: number[];
      }>;
    };
  };
}

// WebSocket消息类型
type WebSocketMessage =
  | { type: 'audio'; data: string; sequence: number; sample_rate?: number }
  | { type: 'audio_end'; sequence: number }
  | { type: 'close'; reason?: string }
  | { type: 'init'; conversation_id: string; session_id: string }
  | { type: 'recognition_interim'; text: string; is_final: boolean }
  | { type: 'recognition_final'; text: string; is_final: boolean }
  | { type: 'bot_response_text'; text: string; message_id: string }
  | { type: 'bot_response_audio'; audio_url: string; text: string }
  | { type: 'error'; error: { code: string; message: string } };

// API客户端
class MultimodalApiClient {
  async createSTTTask(audioUrl: string, language: string): Promise<STTResult> {
    const response = await apiClient.post('/speech/stt', {
      data: {
        audio_url: audioUrl,
        language: language,
        model: 'whisper-1',
      },
    });
    return response.data.data;
  }

  async createTTSTask(text: string, voice: string, language: string): Promise<TTSResult> {
    const response = await apiClient.post('/speech/tts', {
      data: {
        text,
        voice,
        language,
        model: 'azure-neural',
      },
    });
    return response.data.data;
  }

  async getVoices(language?: string): Promise<Voice[]> {
    const response = await apiClient.get('/speech/tts/voices', {
      params: { language },
    });
    return response.data.data;
  }

  async understandVideo(
    videoUrl: string,
    tasks: string[]
  ): Promise<{ task_id: string }> {
    const response = await apiClient.post('/video/understand', {
      data: {
        video_url: videoUrl,
        tasks,
      },
    });
    return response.data.data;
  }

  async getVideoResult(taskId: string): Promise<VideoUnderstandingResult> {
    const response = await apiClient.get(`/video/understand/${taskId}`);
    return response.data.data;
  }

  // WebSocket连接
  connectVoiceConversation(
    token: string,
    botId: number,
    voice: string,
    language: string,
    onMessage: (msg: WebSocketMessage) => void
  ): WebSocket {
    const ws = new WebSocket(
      `wss://api.coze.com/api/v1/speech/conversation?token=${token}&bot_id=${botId}&voice=${voice}&language=${language}`
    );

    ws.onmessage = (event) => {
      const msg: WebSocketMessage = JSON.parse(event.data);
      onMessage(msg);
    };

    return ws;
  }
}

export const multimodalApi = new MultimodalApiClient();
```

---

**文档版本历史**:
- v1.0.0 (2025-12-30): 初始版本,包含语音识别(STT)、语音合成(TTS)、实时语音对话、视频理解、多模态编排、语音记录管理等6个模块的完整API接口定义
