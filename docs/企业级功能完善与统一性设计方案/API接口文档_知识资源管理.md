# API接口文档：知识资源管理模块

**模块名称**: 知识资源管理 (Knowledge Management)
**设计文档**: 09-企业管理中台_知识资源管理.md + 09-知识资源管理_RAG优化补充_完整版.md
**版本**: v2.0.0 (含RAG优化)
**最后更新**: 2025-01-03
**优先级**: P1

---

## 目录

- [1. 模块概述](#1-模块概述)
- [2. 知识库管理API](#2-知识库管理api)
- [3. 文档管理API](#3-文档管理api)
- [4. 文档解析器配置API](#4-文档解析器配置api)
- [5. 分块策略管理API](#5-分块策略管理api)
- [6. 智能分块API](#6-智能分块api)
- [7. 混合检索API](#7-混合检索api)
- [8. RAG知识问答API](#8-rag知识问答api)
- [9. 知识库权限管理API](#9-知识库权限管理api)
- [10. 引用标注API](#10-引用标注api)
- [11. 数据模型](#11-数据模型)
- [12. 错误码定义](#12-错误码定义)

---

## 1. 模块概述

### 1.1 功能说明

知识资源管理是ZKER企业级SaaS平台的核心模块,通过**可扩展文档处理框架 + 智能分块 + 混合检索 + RAG优化**,提供企业级知识库管理能力:

- ✅ **多格式文档支持** - PDF、Word、Excel、PPT、TXT、Markdown、HTML等
- ✅ **智能分块策略** - 固定长度、语义分块、递归分块、混合策略
- ✅ **混合检索** - 向量检索 + 关键词检索(BM25) + 重排序(Reranker)
- ✅ **查询优化** - Query Rewriting + HyDE
- ✅ **RAG问答** - 基于知识库的智能问答,支持引用标注
- ✅ **权限控制** - 基于RBAC的知识库授权

### 1.2 技术架构

**后端技术栈**:
- 框架: CloudWeGo Hertz
- 文档解析: 可扩展解析器框架 (PyPDF2、python-docx、openpyxl等)
- 分块策略: 策略模式 + 数据库配置
- 向量数据库: Milvus 2.5.10
- 关键词检索: Elasticsearch 8.18.0 (BM25)
- 向量化: Embedding Service (text-embedding-ada-002)
- LLM: GPT-4 / Claude / 其他
- 文件存储: MinIO

**前端技术栈**:
- 框架: React 18 + TypeScript
- UI库: Semi Design
- 文件上传: Upload组件
- Markdown渲染: react-markdown

### 1.3 核心特性对比

| 特性 | 原设计 | RAG优化版 | 提升 |
|------|--------|----------|------|
| 分块策略 | 固定长度 | 4种智能策略 | +30%准确率 |
| 检索方式 | 仅向量 | 混合检索+Reranker | +40%召回率 |
| 查询优化 | 无 | Query Rewriting+HyDE | +25%准确率 |
| 引用标注 | 无 | 完整引用系统 | 可追溯性 |

---

## 2. 知识库管理API

### 2.1 创建知识库

**接口地址**: `POST /api/v1/knowledge/bases`

**功能说明**: 创建新知识库

**请求头**:
```http
Content-Type: application/json
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```json
{
  "name": "产品知识库",
  "description": "公司产品相关的文档和手册",
  "category_id": 1,
  "chunk_strategy_id": 1,
  "embedding_model": "text-embedding-ada-002"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "知识库创建成功",
  "data": {
    "id": "kb-20250103-001",
    "tenant_id": "tenant-123",
    "creator_id": 1001,
    "name": "产品知识库",
    "description": "公司产品相关的文档和手册",
    "category_id": 1,
    "chunk_strategy_id": 1,
    "embedding_model": "text-embedding-ada-002",
    "document_count": 0,
    "chunk_count": 0,
    "question_count": 0,
    "hit_rate": 0.0,
    "is_active": true,
    "created_at": "2025-01-03T10:00:00Z"
  }
}
```

### 2.2 查询知识库列表

**接口地址**: `GET /api/v1/knowledge/bases`

**功能说明**: 获取当前租户的知识库列表

**查询参数**:
- page: 页码,默认1
- page_size: 每页数量,默认20
- category_id: 分类ID过滤
- keyword: 搜索关键词

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 15,
    "items": [
      {
        "id": "kb-001",
        "name": "产品知识库",
        "description": "公司产品相关的文档和手册",
        "document_count": 25,
        "chunk_count": 1250,
        "question_count": 156,
        "hit_rate": 0.85,
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 2.3 查看知识库详情

**接口地址**: `GET /api/v1/knowledge/bases/{knowledge_base_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "kb-001",
    "name": "产品知识库",
    "description": "公司产品相关的文档和手册",
    "chunk_strategy": {
      "id": 1,
      "name": "固定长度分块",
      "strategy_type": "fixed_length"
    },
    "statistics": {
      "document_count": 25,
      "chunk_count": 1250,
      "total_size": 52428800,
      "question_count": 156,
      "hit_rate": 0.85
    }
  }
}
```

### 2.4 更新知识库

**接口地址**: `PUT /api/v1/knowledge/bases/{knowledge_base_id}`

**请求参数**:
```json
{
  "name": "产品知识库(更新版)",
  "description": "更新后的描述",
  "chunk_strategy_id": 2
}
```

### 2.5 删除知识库

**接口地址**: `DELETE /api/v1/knowledge/bases/{knowledge_base_id}`

**功能说明**: 软删除知识库及其所有文档

---

## 3. 文档管理API

### 3.1 上传文档

**接口地址**: `POST /api/v1/knowledge/bases/{knowledge_base_id}/documents`

**功能说明**: 上传文档到知识库,自动触发解析、分块、向量化

**请求头**:
```http
Content-Type: multipart/form-data
Authorization: Bearer {jwt_token}
X-Tenant-ID: {tenant_id}
```

**请求参数**:
```
file: <binary> (必需)
chunk_strategy_id: 1 (可选,覆盖知识库默认策略)
```

**响应示例**:
```json
{
  "code": 0,
  "message": "文档上传成功,正在处理中",
  "data": {
    "document_id": "doc-20250103-001",
    "filename": "产品手册.pdf",
    "file_size": 2048576,
    "file_type": "pdf",
    "status": "pending",
    "created_at": "2025-01-03T10:30:00Z"
  }
}
```

**Go后端代码示例**:

```go
// UploadDocument 上传文档
func (s *KnowledgeService) UploadDocument(
    ctx context.Context,
    knowledgeBaseID string,
    fileHeader *multipart.FileHeader,
    chunkStrategyID *int64,
) (*UploadDocumentResponse, error) {
    // 1. 验证知识库权限
    if !s.checkPermission(ctx, knowledgeBaseID, "upload") {
        return nil, errors.New("permission denied")
    }

    // 2. 验证文件类型
    fileType := strings.TrimPrefix(filepath.Ext(fileHeader.Filename), ".")
    parser, err := s.parserRepo.GetByFileType(ctx, fileType)
    if err != nil {
        return nil, fmt.Errorf("unsupported file type: %s", fileType)
    }

    // 3. 验证文件大小
    if fileHeader.Size > parser.MaxFileSize {
        return nil, fmt.Errorf("file too large: max %d bytes", parser.MaxFileSize)
    }

    // 4. 存储文件到 MinIO
    filePath := fmt.Sprintf("knowledge/%s/%s", knowledgeBaseID, generateUniqueFilename(fileHeader.Filename))
    reader, err := fileHeader.Open()
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    if err := s.minioClient.UploadFile(ctx, filePath, reader, fileHeader.Size); err != nil {
        return nil, err
    }

    // 5. 创建文档记录
    document := &model.Document{
        ID:              generateDocumentID(),
        TenantID:        getTenantID(ctx),
        KnowledgeBaseID: knowledgeBaseID,
        UploaderID:      getUserID(ctx),
        Filename:        fileHeader.Filename,
        FileType:        fileType,
        FileSize:        fileHeader.Size,
        FilePath:        filePath,
        Status:          "pending",
    }

    if err := s.documentRepo.Create(ctx, document); err != nil {
        return nil, err
    }

    // 6. 异步处理文档
    go s.processDocumentAsync(context.Background(), document, parser, chunkStrategyID)

    return &UploadDocumentResponse{
        DocumentID: document.ID,
        Filename:   document.Filename,
        FileSize:   document.FileSize,
        Status:     document.Status,
        CreatedAt:  document.CreatedAt,
    }, nil
}
```

**前端代码示例**:

```tsx
import React, { useState } from 'react';
import { Upload, Button, Progress, message } from '@douyinfe/semi-ui';
import { IconUpload } from '@douyinfe/semi-icons';
import { knowledgeAPI } from '@/services/api';

export const KnowledgeBaseUploader: React.FC<{ knowledgeBaseId: string }> = ({ knowledgeBaseId }) => {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);

  const handleUpload = async (fileList: File[]) => {
    setUploading(true);
    setProgress(0);

    try {
      for (let i = 0; i < fileList.length; i++) {
        const file = fileList[i];
        const formData = new FormData();
        formData.append('file', file);

        await knowledgeAPI.uploadDocument(knowledgeBaseId, formData, {
          onUploadProgress: (progressEvent) => {
            const percentCompleted = Math.round(((i * 100 + progressEvent.loaded / progressEvent.total * 100) / fileList.length));
            setProgress(percentCompleted);
          },
        });
      }

      message.success(`成功上传 ${fileList.length} 个文档`);
    } catch (error) {
      message.error('上传失败');
    } finally {
      setUploading(false);
      setProgress(0);
    }
  };

  return (
    <Upload
      draggable
      multiple
      accept=".pdf,.doc,.docx,.xlsx,.ppt,.pptx,.txt,.md"
      showUploadList={false}
      customRequest={({ fileList }) => handleUpload(fileList)}
    >
      <Button icon={<IconUpload />} loading={uploading} theme="solid">
        {uploading ? `上传中... ${progress}%` : '点击或拖拽上传文档'}
      </Button>
      {uploading && <Progress percent={progress} showInfo={false} />}
    </Upload>
  );
};
```

### 3.2 查询文档列表

**接口地址**: `GET /api/v1/knowledge/bases/{knowledge_base_id}/documents`

**查询参数**:
- page: 页码
- page_size: 每页数量
- status: 状态过滤 (pending/parsing/completed/failed)
- file_type: 文件类型过滤

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 25,
    "items": [
      {
        "id": "doc-001",
        "filename": "产品手册.pdf",
        "file_type": "pdf",
        "file_size": 2048576,
        "status": "completed",
        "chunk_count": 125,
        "vector_count": 125,
        "question_count": 10,
        "created_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 3.3 查看文档详情

**接口地址**: `GET /api/v1/knowledge/documents/{document_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "doc-001",
    "filename": "产品手册.pdf",
    "file_size": 2048576,
    "status": "completed",
    "parsed_content": "解析后的文本内容...",
    "chunk_count": 125,
    "chunks": [
      {
        "id": 1,
        "chunk_index": 0,
        "content": "第一章 产品概述...",
        "hit_count": 5
      }
    ]
  }
}
```

### 3.4 重新解析文档

**接口地址**: `POST /api/v1/knowledge/documents/{document_id}/reparse`

**请求参数**:
```json
{
  "chunk_strategy_id": 2
}
```

**功能说明**: 使用新的分块策略重新解析文档

### 3.5 删除文档

**接口地址**: `DELETE /api/v1/knowledge/documents/{document_id}`

---

## 4. 文档解析器配置API

### 4.1 查询解析器列表

**接口地址**: `GET /api/v1/knowledge/parsers`

**查询参数**:
- file_type: 文件类型过滤
- is_active: 是否启用

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "PDFParser",
        "file_types": ["pdf"],
        "parser_type": "library",
        "can_extract_text": true,
        "can_extract_tables": false,
        "can_extract_images": false,
        "max_file_size": 52428800,
        "timeout_seconds": 300
      }
    ]
  }
}
```

### 4.2 查看解析器配置

**接口地址**: `GET /api/v1/knowledge/parsers/{parser_id}`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "PDFParser",
    "config": {
      "library": "PyPDF2",
      "version": "3.0.0",
      "options": {
        "extractText": true,
        "extractImages": false,
        "extractTables": false
      }
    }
  }
}
```

---

## 5. 分块策略管理API

### 5.1 查询分块策略列表

**接口地址**: `GET /api/v1/knowledge/chunk-strategies`

**查询参数**:
- strategy_type: 策略类型 (fixed_length/semantic/section/recursive)
- is_default: 是否默认策略

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "固定长度分块",
        "description": "每块最多1000字符,重叠200字符",
        "strategy_type": "fixed_length",
        "config": {
          "chunkSize": 1000,
          "chunkOverlap": 200,
          "separator": "\n\n"
        },
        "is_default": true
      },
      {
        "id": 2,
        "name": "语义分块",
        "description": "基于语义相似度智能分块",
        "strategy_type": "semantic",
        "config": {
          "embeddingModel": "text-embedding-ada-002",
          "similarityThreshold": 0.7,
          "minChunkSize": 200,
          "maxChunkSize": 2000
        }
      }
    ]
  }
}
```

### 5.2 创建自定义分块策略

**接口地址**: `POST /api/v1/knowledge/chunk-strategies`

**请求参数**:
```json
{
  "name": "自定义语义分块",
  "description": "用于技术文档的语义分块策略",
  "strategy_type": "semantic",
  "config": {
    "embeddingModel": "text-embedding-ada-002",
    "similarityThreshold": 0.6,
    "minChunkSize": 300,
    "maxChunkSize": 1500
  },
  "embedding_model": "text-embedding-ada-002",
  "chunk_overlap": 100
}
```

---

## 6. 智能分块API

### 6.1 执行智能分块

**接口地址**: `POST /api/v1/knowledge/documents/{document_id}/chunk`

**功能说明**: 使用指定的分块策略对文档进行智能分块

**请求参数**:
```json
{
  "chunking_strategy": "semantic",
  "chunk_size": 512,
  "overlap": 50,
  "separator": "\n\n"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "分块完成",
  "data": {
    "document_id": "doc-123",
    "total_chunks": 25,
    "chunks": [
      {
        "id": "chunk-001",
        "chunk_index": 0,
        "chunk_content": "第一章 产品概述...",
        "token_count": 512,
        "embedding_vector": [0.1, 0.2, ...]
      }
    ]
  }
}
```

**Go后端代码示例**:

```go
// SemanticChunking 语义分块
func SemanticChunking(document string, config *ChunkingConfig) []Chunk {
    sentences := splitSentences(document)

    // 使用embedding计算句子相似度
    embeddings := calculateEmbeddings(sentences)

    chunks := make([]Chunk, 0)
    currentChunk := make([]string, 0)
    currentSize := 0

    for i, sentence := range sentences {
        // 检查语义边界
        if len(currentChunk) > 0 && isSemanticBoundary(embeddings[i-1], embeddings[i]) {
            if currentSize >= config.MinChunkSize {
                chunks = append(chunks, Chunk{
                    Content: strings.Join(currentChunk, " "),
                    TokenCount: currentSize,
                    Index: len(chunks),
                })
                currentChunk = make([]string, 0)
                currentSize = 0
            }
        }

        currentChunk = append(currentChunk, sentence)
        currentSize += countTokens(sentence)

        // 达到最大块大小,强制分块
        if currentSize >= config.MaxChunkSize {
            chunks = append(chunks, Chunk{
                Content: strings.Join(currentChunk, " "),
                TokenCount: currentSize,
                Index: len(chunks),
            })
            currentChunk = make([]string, 0)
            currentSize = 0
        }
    }

    return chunks
}

func isSemanticBoundary(emb1, emb2 []float64) bool {
    // 计算余弦相似度
    similarity := cosineSimilarity(emb1, emb2)
    // 如果相似度低,说明是语义边界
    return similarity < 0.5
}
```

---

## 7. 混合检索API

### 7.1 混合检索(向量+关键词+重排序)

**接口地址**: `POST /api/v1/knowledge/search`

**功能说明**: 使用混合检索算法(向量+关键词+Reranker)进行知识检索

**请求参数**:
```json
{
  "query": "怎么重置密码?",
  "vector_weight": 0.7,
  "keyword_weight": 0.3,
  "reranker_enabled": true,
  "top_k": 20,
  "reranker_top_k": 10
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "query": "怎么重置密码?",
    "total_results": 10,
    "results": [
      {
        "chunk_id": "chunk-001",
        "content": "您可以通过以下步骤重置密码...",
        "document_id": "doc-123",
        "document_title": "用户手册",
        "score": 0.95,
        "source": "hybrid",
        "confidence": 0.92,
        "metadata": {
          "page": 15,
          "heading": "账号管理"
        }
      }
    ]
  }
}
```

**Go后端代码示例**:

```go
// HybridSearch 混合检索
func (e *HybridSearchEngine) HybridSearch(
    ctx context.Context,
    query string,
    config *HybridSearchConfig,
) ([]SearchResult, error) {

    // 1. 并行执行向量检索和关键词检索
    var vectorResults, keywordResults []SearchResult
    var wg sync.WaitGroup
    var errV, errK error

    wg.Add(2)
    go func() {
        vectorResults, errV = e.VectorSearch(ctx, query, config)
        wg.Done()
    }()
    go func() {
        keywordResults, errK = e.KeywordSearch(ctx, query, config)
        wg.Done()
    }()
    wg.Wait()

    if errV != nil {
        return nil, errV
    }
    if errK != nil {
        return nil, errK
    }

    // 2. 结果融合(RRF算法)
    fusedResults := e.fuseResults(vectorResults, keywordResults, config)

    // 3. 重排序(如果启用)
    if config.RerankerEnabled {
        fusedResults = e.reranker.Rerank(ctx, query, fusedResults, config.RerankerTopK)
    }

    return fusedResults, nil
}

// fuseResults RRF算法融合结果
func (e *HybridSearchEngine) fuseResults(
    vectorResults []SearchResult,
    keywordResults []SearchResult,
    config *HybridSearchConfig,
) []SearchResult {

    // RRF公式: score = 1/(k+rank_vector) + 1/(k+rank_keyword)
    k := 60.0

    scores := make(map[string]float64)
    chunkMap := make(map[string]*SearchResult)

    // 处理向量检索结果
    for rank, result := range vectorResults {
        rrfScore := 1.0 / (k + float64(rank+1))
        scores[result.ChunkID] += rrfScore * config.VectorWeight
        chunkMap[result.ChunkID] = result
    }

    // 处理关键词检索结果
    for rank, result := range keywordResults {
        rrfScore := 1.0 / (k + float64(rank+1))
        scores[result.ChunkID] += rrfScore * config.KeywordWeight

        if _, exists := chunkMap[result.ChunkID]; !exists {
            chunkMap[result.ChunkID] = result
        }
    }

    // 按分数排序
    sortedResults := make([]*SearchResult, 0, len(chunkMap))
    for chunkID, score := range scores {
        sortedResults = append(sortedResults, &SearchResult{
            ChunkID: chunkID,
            Score:   score,
            Source:  "hybrid",
        })
    }

    sort.Slice(sortedResults, func(i, j int) bool {
        return sortedResults[i].Score > sortedResults[j].Score
    })

    // 转换为[]SearchResult
    results := make([]SearchResult, len(sortedResults))
    for i, r := range sortedResults {
        results[i] = *r
    }

    return results
}
```

---

## 8. RAG知识问答API

### 8.1 知识问答(基础RAG)

**接口地址**: `POST /api/v1/knowledge/bases/{knowledge_base_id}/ask`

**请求参数**:
```json
{
  "question": "如何申请年假?",
  "top_k": 5,
  "include_sources": true
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "answer": "根据公司员工手册,年假申请流程如下:\n\n1. 登录企业OA系统\n2. 进入\"请假申请\"模块\n3. 选择\"年假\"类型\n4. 填写请假时间和事由\n5. 提交审批\n\n年假天数根据工龄计算:\n- 工龄1-10年:每年5天\n- 工龄10-20年:每年10天\n- 工龄20年以上:每年15天",
    "sources": [
      {
        "document_id": "doc-001",
        "document_name": "员工手册.pdf",
        "chunk_id": 123,
        "content": "年假申请流程...",
        "score": 0.92
      }
    ],
    "query_used": "如何申请年假?"
  }
}
```

### 8.2 知识问答(RAG优化版 - 含查询优化)

**接口地址**: `POST /api/v1/knowledge/bases/{knowledge_base_id}/ask/optimized`

**功能说明**: 使用查询优化(Query Rewriting + HyDE) + 混合检索的增强版RAG

**请求参数**:
```json
{
  "question": "怎么休年假?",
  "enable_query_rewrite": true,
  "enable_hyde": true,
  "enable_hybrid_search": true,
  "enable_reranker": true,
  "top_k": 20,
  "reranker_top_k": 10,
  "conversation_history": [
    {
      "role": "user",
      "content": "我想休假"
    },
    {
      "role": "assistant",
      "content": "您想请什么类型的假?"
    }
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "answer": "根据公司员工手册,年假申请流程如下:\n\n1. 登录企业OA系统\n2. 进入\"请假申请\"模块\n3. 选择\"年假\"类型\n4. 填写请假时间和事由\n5. 提交审批\n\n年假天数根据工龄计算:\n- 工龄1-10年:每年5天\n- 工龄10-20年:每年10天\n- 工龄20年以上:每年15天",
    "sources": [
      {
        "document_id": "doc-001",
        "document_name": "员工手册.pdf",
        "chunk_id": 123,
        "content": "年假申请流程...",
        "score": 0.95,
        "confidence": 0.92
      }
    ],
    "optimization_info": {
      "original_query": "怎么休年假?",
      "rewritten_query": "如何申请年假?申请流程是什么?",
      "hypothetical_doc": "年假申请流程:1.登录OA系统 2.进入请假申请 3.选择年假类型...",
      "search_method": "hybrid",
      "reranked": true
    }
  }
}
```

**Go后端代码示例**:

```go
// AskQuestionOptimized RAG优化版问答
func (s *KnowledgeService) AskQuestionOptimized(
    ctx context.Context,
    knowledgeBaseID string,
    req *AskQuestionOptimizedRequest,
) (*AskQuestionResponse, error) {

    var finalQuery string
    var hydeContent string

    // 1. 查询重写(Query Rewriting)
    if req.EnableQueryRewrite && len(req.ConversationHistory) > 0 {
        rewrittenQuery, err := s.queryRewritingService.RewriteQuery(
            ctx,
            req.Question,
            req.ConversationHistory,
        )
        if err == nil {
            finalQuery = rewrittenQuery
        } else {
            finalQuery = req.Question
        }
    } else {
        finalQuery = req.Question
    }

    // 2. HyDE (假设性文档嵌入)
    if req.EnableHyDE {
        hypotheticalDoc, err := s.queryRewritingService.HyDE(ctx, finalQuery)
        if err == nil {
            hydeContent = hypotheticalDoc
        }
    }

    // 3. 混合检索
    var chunks []*Chunk
    if req.EnableHybridSearch {
        // 向量检索 + 关键词检索 + Reranker
        searchConfig := &HybridSearchConfig{
            VectorWeight:    0.7,
            KeywordWeight:   0.3,
            RerankerEnabled: req.EnableReranker,
            TopK:            req.TopK,
            RerankerTopK:    req.RerankerTopK,
        }

        // 如果启用HyDE,将假设性文档与查询合并
        searchQuery := finalQuery
        if hydeContent != "" {
            searchQuery = finalQuery + " " + hydeContent
        }

        results, err := s.hybridSearchEngine.HybridSearch(ctx, searchQuery, searchConfig)
        if err == nil {
            chunkIDs := make([]int64, len(results))
            for i, r := range results {
                chunkIDs[i] = r.ChunkID
            }
            chunks, _ = s.chunkRepo.GetByIDs(ctx, chunkIDs)
        }
    } else {
        // 仅向量检索
        queryVector, _ := s.embeddingService.Embed(ctx, finalQuery)
        searchResults, _ := s.vectorStore.Search(ctx, knowledgeBaseID, queryVector, req.TopK)
        chunks, _ = s.chunkRepo.GetByVectorIDs(ctx, searchResults.VectorIDs())
    }

    if len(chunks) == 0 {
        return &AskQuestionResponse{
            Answer:  "抱歉,我在知识库中没有找到相关信息。",
            Sources: []Source{},
        }, nil
    }

    // 4. 构建提示词
    contextText := s.buildContext(chunks)
    prompt := s.buildRAGPrompt(finalQuery, contextText)

    // 5. 调用 LLM
    llmResp, err := s.llmClient.Chat(ctx, &llm.ChatRequest{
        Model: "gpt-4",
        Messages: []llm.Message{
            {Role: "system", Content: "你是一个专业的知识问答助手,请基于提供的知识库内容回答用户问题。"},
            {Role: "user", Content: prompt},
        },
        Temperature: 0.3,
        MaxTokens:   2000,
    })
    if err != nil {
        return nil, err
    }

    // 6. 构建响应
    sources := make([]Source, 0, len(chunks))
    if req.IncludeSources {
        for i, chunk := range chunks {
            sources = append(sources, Source{
                DocumentID:   chunk.DocumentID,
                DocumentName: chunk.DocumentName,
                ChunkID:      chunk.ID,
                Content:      chunk.Content,
                Score:        chunks[i].Score,
            })
        }
    }

    return &AskQuestionResponse{
        Answer:  llmResp.Content,
        Sources: sources,
    }, nil
}
```

### 8.3 查询问答历史

**接口地址**: `GET /api/v1/knowledge/bases/{knowledge_base_id}/questions`

**查询参数**:
- page: 页码
- page_size: 每页数量

---

## 9. 知识库权限管理API

### 9.1 查看权限列表

**接口地址**: `GET /api/v1/knowledge/bases/{knowledge_base_id}/permissions`

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {
        "id": 1,
        "knowledge_base_id": "kb-001",
        "grantee_type": "role",
        "grantee_name": "全员",
        "permissions": ["read", "ask"],
        "granted_by": 1001,
        "granted_at": "2025-01-03T10:00:00Z"
      }
    ]
  }
}
```

### 9.2 授权访问

**接口地址**: `POST /api/v1/knowledge/bases/{knowledge_base_id}/permissions`

**请求参数**:
```json
{
  "grantee_type": "user",
  "grantee_id": "1001",
  "permissions": ["read", "ask", "manage"]
}
```

### 9.3 撤销权限

**接口地址**: `DELETE /api/v1/knowledge/bases/{knowledge_base_id}/permissions/{permission_id}`

---

## 10. 引用标注API

### 10.1 查询消息引用

**接口地址**: `GET /api/v1/conversations/{conversation_id}/messages/{message_id}/citations`

**功能说明**: 获取AI回复时引用的知识来源

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message_id": "msg-123",
    "citations": [
      {
        "id": "cit-001",
        "chunk_id": "chunk-001",
        "chunk_content": "年假申请流程...",
        "document_id": "doc-123",
        "document_title": "员工手册.pdf",
        "confidence": 0.92,
        "relevance_score": 0.95,
        "snippet": "您可以按以下步骤重置密码..."
      }
    ]
  }
}
```

**Go后端代码示例**:

```go
// TrackCitations 追踪知识引用
func (s *CitationService) TrackCitations(
    ctx context.Context,
    conversationID string,
    messageID string,
    query string,
    retrievedChunks []Chunk,
    response string,
) error {

    // 1. 计算每个chunk的置信度
    for _, chunk := range retrievedChunks {
        // 计算chunk在response中的引用程度
        confidence := s.calculateConfidence(ctx, response, chunk.Content)

        // 计算相关性
        relevance := s.calculateRelevance(ctx, query, chunk.Content)

        // 保存引用记录
        citation := &KnowledgeCitation{
            TenantID:        getTenantID(ctx),
            ConversationID:  conversationID,
            MessageID:       messageID,
            ChunkID:         chunk.ID,
            ChunkContent:    chunk.Content,
            DocumentID:      chunk.DocumentID,
            DocumentTitle:   chunk.DocumentTitle,
            Confidence:      confidence,
            RelevanceScore:  relevance,
        }

        if err := s.db.WithContext(ctx).Create(citation).Error; err != nil {
            return err
        }
    }

    return nil
}
```

---

## 11. 数据模型

### 11.1 KnowledgeBase(知识库)
```typescript
interface KnowledgeBase {
  id: string;
  tenant_id: string;
  creator_id: number;
  name: string;
  description?: string;
  category_id?: number;

  // 配置
  chunk_strategy_id?: number;
  embedding_model: string;

  // 统计
  document_count: number;
  chunk_count: number;
  total_size: number;
  question_count: number;
  hit_rate: number;

  // 状态
  is_active: boolean;

  // 审计
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}
```

### 11.2 Document(文档)
```typescript
interface Document {
  id: string;
  tenant_id: string;
  knowledge_base_id: string;
  uploader_id: number;

  // 文件信息
  filename: string;
  file_type: string;
  file_size: number;
  file_path: string;
  file_hash?: string;

  // 解析结果
  raw_content?: string;
  parsed_content?: string;
  chunk_count: number;
  vector_count: number;

  // 状态
  status: 'pending' | 'parsing' | 'chunking' | 'vectorizing' | 'completed' | 'failed';
  error_message?: string;

  // 统计
  question_count: number;
  hit_count: number;

  // 审计
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}
```

### 11.3 DocumentChunk(文档分块)
```typescript
interface DocumentChunk {
  id: number;
  tenant_id: string;
  document_id: string;
  knowledge_base_id: string;

  // 分块内容
  content: string;
  chunk_index: number;
  start_pos?: number;
  end_pos?: number;

  // 元数据
  metadata?: Record<string, any>;
  vector_id?: string;

  // 统计
  hit_count: number;

  // 索引
  created_at: string;
}
```

### 11.4 ChunkingStrategy(分块策略)
```typescript
type ChunkingStrategy = 'fixed' | 'semantic' | 'recursive' | 'hybrid';

interface ChunkingConfig {
  chunk_size: number;
  overlap: number;
  separator?: string;
  min_chunk_size?: number;
  max_chunk_size?: number;
  similarity_threshold?: number; // for semantic
  embedding_model?: string;
}

interface ChunkStrategy {
  id: number;
  tenant_id?: string;
  name: string;
  description?: string;
  strategy_type: ChunkingStrategy;
  config: ChunkingConfig;
  embedding_model: string;
  chunk_overlap: number;
  is_active: boolean;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}
```

### 11.5 HybridSearchConfig(混合检索配置)
```typescript
interface HybridSearchConfig {
  vector_weight: number;    // 0-1
  keyword_weight: number;   // 0-1
  reranker_enabled: boolean;
  top_k: number;
  reranker_top_k: number;
}

interface SearchResult {
  chunk_id: string | number;
  content: string;
  document_id: string;
  document_title: string;
  score: number;
  source: 'vector' | 'keyword' | 'hybrid';
  confidence?: number;
  metadata?: Record<string, any>;
}
```

### 11.6 Citation(知识引用)
```typescript
interface KnowledgeCitation {
  id: number;
  tenant_id: string;
  conversation_id: string;
  message_id: string;

  // 引用信息
  chunk_id: number;
  chunk_content: string;
  document_id: string;
  document_title: string;

  // 置信度
  confidence: number;         // 0-1
  relevance_score: number;    // 0-1

  // 位置信息
  start_pos?: number;
  end_pos?: number;

  created_at: string;
}
```

---

## 12. 错误码定义

| 错误码 | HTTP状态码 | 说明 |
|--------|-----------|------|
| 40001 | 400 | 知识库不存在 |
| 40002 | 400 | 知识库名称已存在 |
| 40101 | 400 | 文档不存在 |
| 40102 | 400 | 不支持的文件类型 |
| 40103 | 400 | 文件大小超过限制 |
| 40104 | 500 | 文档解析失败 |
| 40105 | 500 | 文档分块失败 |
| 40106 | 500 | 向量化失败 |
| 40201 | 403 | 无知识库访问权限 |
| 40301 | 400 | 分块策略不存在 |
| 40401 | 400 | 检索失败 |
| 40402 | 400 | 未找到相关知识 |

---

**文档结束**
