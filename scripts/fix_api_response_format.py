#!/usr/bin/env python3
"""
API响应格式统一修复脚本
将所有非标准响应格式替换为企业级统一响应格式

修复规则:
1. c.JSON(consts.StatusOK, resp) -> httputil.BuildSuccessResp(c, resp)
2. invalidParamRequestResponse(c, errMsg) -> httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, errMsg, "参数验证失败", nil)
3. internalServerErrorResponse(ctx, c, err) -> 使用增强错误码

作者: API规范化工具
日期: 2025-01-01
"""

import os
import re
from pathlib import Path

# 需要修复的文件列表
FILES_TO_FIX = [
    "backend/api/handler/coze/conversation_service.go",
    "backend/api/handler/coze/message_service.go",
    "backend/api/handler/coze/knowledge_service.go",
    "backend/api/handler/coze/intelligence_service.go",
    "backend/api/handler/coze/workflow_service.go",
    "backend/api/handler/coze/bot_open_api_service.go",
    "backend/api/handler/coze/config_service.go",
    "backend/api/handler/coze/resource_service.go",
    "backend/api/handler/coze/playground_service.go",
    "backend/api/handler/coze/database_service.go",
]

# 添加httputil导入的模式
IMPORT_PATTERN = r'("github.com/coze-dev/coze-studio/backend/api/model/[^"]+")'
IMPORT_REPLACEMENT = r'\1\n\t"github.com/coze-dev/coze-studio/backend/api/internal/httputil"'

# 替换规则
REPLACEMENTS = [
    # 成功响应
    (r'c\.JSON\(consts\.StatusOK, (\w+)\)', r'httputil.BuildSuccessResp(c, \1)'),
    (r'c\.JSON\(http\.StatusOK, (\w+)\)', r'httputil.BuildSuccessResp(c, \1)'),

    # 参数错误响应
    (
        r'invalidParamRequestResponse\(c, ([^)]+)\)',
        r'httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, \1, "参数验证失败", nil)'
    ),

    # 服务器错误响应 - 使用增强错误码
    (
        r'internalServerErrorResponse\(ctx, c, err\)',
        r'httputil.BuildErrorRespFromEnhanced(c, errno.NewInternalError(ctx, err))'
    ),
]


def fix_file(file_path: str) -> bool:
    """修复单个文件的响应格式"""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()

        original_content = content

        # 检查是否已导入httputil
        if '"github.com/coze-dev/coze-studio/backend/api/internal/httputil"' not in content:
            # 在导入部分添加httputil
            content = re.sub(IMPORT_PATTERN, IMPORT_REPLACEMENT, content, count=1)

        # 应用所有替换规则
        for pattern, replacement in REPLACEMENTS:
            content = re.sub(pattern, replacement, content)

        # 如果有修改,写回文件
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"✅ 已修复: {file_path}")
            return True
        else:
            print(f"⏭️  无需修复: {file_path}")
            return False

    except Exception as e:
        print(f"❌ 修复失败 {file_path}: {e}")
        return False


def main():
    """主函数"""
    print("=" * 80)
    print("API响应格式统一修复脚本")
    print("=" * 80)

    # 统计
    fixed_count = 0
    skipped_count = 0
    failed_count = 0

    # 获取项目根目录
    project_root = Path(__file__).parent.parent
    os.chdir(project_root)

    # 修复所有文件
    for file_path in FILES_TO_FIX:
        if not Path(file_path).exists():
            print(f"⚠️  文件不存在: {file_path}")
            failed_count += 1
            continue

        if fix_file(file_path):
            fixed_count += 1
        else:
            skipped_count += 1

    # 输出统计
    print("\n" + "=" * 80)
    print("修复统计:")
    print(f"  ✅ 已修复: {fixed_count} 个文件")
    print(f"  ⏭️  无需修复: {skipped_count} 个文件")
    print(f"  ❌ 失败: {failed_count} 个文件")
    print("=" * 80)

    if fixed_count > 0:
        print("\n⚠️  请执行以下命令验证修改:")
        print("   cd backend && go build ./...")
        print("   git diff")


if __name__ == "__main__":
    main()
