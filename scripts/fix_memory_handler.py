#!/usr/bin/env python3
"""
修复memory_handler.go中的响应格式
"""

import re

file_path = "backend/api/handler/coze/memory/memory_handler.go"

with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

# 替换规则
replacements = [
    # StoreConversationMemory - 参数错误
    (
        r'c\.JSON\(consts\.StatusBadRequest, memory\.StoreConversationMemoryResponse\{\s*Code:.*?\}\s*\)',
        'httputil.BuildErrorResp(c, errno.ErrMemoryInvalidParamCode, err.Error(), "参数验证失败", nil)'
    ),
    # StoreConversationMemory - 服务器错误
    (
        r'c\.JSON\(consts\.StatusInternalServerError, memory\.StoreConversationMemoryResponse\{\s*Code:.*?Message:.*?"Failed to store memory".*?\}\s*\)',
        'httputil.BuildErrorResp(c, errno.ErrMemoryInternalErrorCode, err.Error(), "存储记忆失败", nil)'
    ),
    # StoreConversationMemory - 成功
    (
        r'c\.JSON\(consts\.StatusOK, memory\.StoreConversationMemoryResponse\{\s*Code:.*?Message:.*?"Memory stored successfully".*?Data: &memory\.MemoryData\{.*?\}.*?\}\s*\)',
        'httputil.BuildSuccessResp(c, &memory.MemoryData{\n\t\tMemoryID:  mem.MemoryID,\n\t\tVectorID:  mem.MemoryID,\n\t\tType:      req.MemoryType,\n\t\tCreatedAt: mem.CreatedAt.Format(time.RFC3339),\n\t})'
    ),
]

for pattern, replacement in replacements:
    content = re.sub(pattern, replacement, content, flags=re.DOTALL)

# 手动处理复杂的多行结构
lines = content.split('\n')
new_lines = []
i = 0
while i < len(lines):
    line = lines[i]

    # StoreConversationMemory - 参数错误
    if 'c.JSON(consts.StatusBadRequest, memory.StoreConversationMemoryResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInvalidParamCode, err.Error(), "参数验证失败", nil)')
        # 跳过整个多行响应
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # StoreConversationMemory - 服务器错误
    elif 'c.JSON(consts.StatusInternalServerError, memory.StoreConversationMemoryResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInternalErrorCode, err.Error(), "存储记忆失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # StoreConversationMemory - 成功
    elif 'c.JSON(consts.StatusOK, memory.StoreConversationMemoryResponse{' in line and 'StoreConversationMemoryResponse{' in line:
        new_lines.append('\thttputil.BuildSuccessResp(c, &memory.MemoryData{')
        i += 1
        # 跳过Code, Message, Timestamp行
        while i < len(lines) and ('Data:' not in lines[i]):
            i += 1
        # 保留Data部分
        if i < len(lines) and 'Data:' in lines[i]:
            new_lines.append('\t\t' + lines[i].split('Data: ')[1].rstrip(','))
            i += 1
            # 保留MemoryID, VectorID等字段
            while i < len(lines) and 'Timestamp:' not in lines[i]:
                if lines[i].strip() and not lines[i].strip().startswith('}'):
                    new_lines.append('\t\t' + lines[i].strip().rstrip(','))
                i += 1
            new_lines.append('\t})')
        # 跳过Timestamp和结束括号
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveMemories - 参数错误
    elif 'c.JSON(consts.StatusBadRequest, memory.RetrieveMemoriesResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInvalidParamCode, err.Error(), "参数验证失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveMemories - 服务器错误
    elif 'c.JSON(consts.StatusInternalServerError, memory.RetrieveMemoriesResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInternalErrorCode, err.Error(), "检索记忆失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveMemories - 成功
    elif 'c.JSON(consts.StatusOK, memory.RetrieveMemoriesResponse{' in line and 'RetrieveMemoriesResponse{' in line:
        new_lines.append('\thttputil.BuildSuccessResp(c, &memory.RetrieveMemoriesResponse{')
        i += 1
        # 跳过Code, Message, Timestamp
        while i < len(lines) and ('Memories:' not in lines[i]):
            i += 1
        # 保留数据部分
        while i < len(lines) and 'Timestamp:' not in lines[i]:
            if lines[i].strip():
                new_lines.append('\t\t' + lines[i].strip().rstrip(','))
            i += 1
        new_lines.append('\t})')
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # StoreKnowledge - 参数错误
    elif 'c.JSON(consts.StatusBadRequest, memory.StoreKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInvalidParamCode, err.Error(), "参数验证失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # StoreKnowledge - 服务器错误
    elif 'c.JSON(consts.StatusInternalServerError, memory.StoreKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInternalErrorCode, err.Error(), "存储知识失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # StoreKnowledge - 成功
    elif 'c.JSON(consts.StatusOK, memory.StoreKnowledgeResponse{' in line and 'StoreKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildSuccessResp(c, &memory.StoreKnowledgeResponse{')
        i += 1
        while i < len(lines) and ('KnowledgeID:' not in lines[i]):
            i += 1
        while i < len(lines) and 'Timestamp:' not in lines[i]:
            if lines[i].strip():
                new_lines.append('\t\t' + lines[i].strip().rstrip(','))
            i += 1
        new_lines.append('\t})')
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveKnowledge - 参数错误
    elif 'c.JSON(consts.StatusBadRequest, memory.RetrieveKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInvalidParamCode, err.Error(), "参数验证失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveKnowledge - 服务器错误
    elif 'c.JSON(consts.StatusInternalServerError, memory.RetrieveKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrMemoryInternalErrorCode, err.Error(), "检索知识失败", nil)')
        i += 1
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    # RetrieveKnowledge - 成功
    elif 'c.JSON(consts.StatusOK, memory.RetrieveKnowledgeResponse{' in line and 'RetrieveKnowledgeResponse{' in line:
        new_lines.append('\thttputil.BuildSuccessResp(c, &memory.RetrieveKnowledgeResponse{')
        i += 1
        while i < len(lines) and ('Knowledge:' not in lines[i]):
            i += 1
        while i < len(lines) and 'Timestamp:' not in lines[i]:
            if lines[i].strip():
                new_lines.append('\t\t' + lines[i].strip().rstrip(','))
            i += 1
        new_lines.append('\t})')
        while i < len(lines) and '})' not in lines[i]:
            i += 1
        i += 1
    else:
        new_lines.append(line)
        i += 1

content = '\n'.join(new_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print("✅ 已修复 memory_handler.go")
