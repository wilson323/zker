#!/usr/bin/env python3
"""
修复 workflow 包中 errno 类型不匹配的问题

将 vo.WrapError/NewError/WrapWarn(errno.ErrXXX, ...)
改为 vo.WrapError/NewError/WrapWarn(errno.DeprecatedErrXXX, ...)
"""

import os
import re
import sys

# errno 映射表：ErrXXX -> DeprecatedErrXXX
# 注意：只映射那些有对应 Deprecated 版本的错误码
ERRNO_MAPPING = {
    'ErrIDGenError': 'DeprecatedErrIDGenError',
    'ErrDatabaseError': 'DeprecatedErrDatabaseError',
    'ErrRedisError': 'DeprecatedErrRedisError',
    'ErrWorkflowNotFound': 'DeprecatedErrWorkflowNotFound',
    'ErrWorkflowSpecifiedVersionNotFound': 'DeprecatedErrWorkflowSnapshotNotFound',
    'ErrWorkflowSnapshotNotFound': 'DeprecatedErrWorkflowSnapshotNotFound',
    'ErrWorkflowExecuteFail': 'DeprecatedErrWorkflowExecuteFail',
    'ErrSerializationDeserializationFail': 'DeprecatedErrSerializationDeserializationFail',
    'ErrSchemaConversionFail': 'DeprecatedErrSchemaConversionFail',
    'ErrInvalidParameter': 'DeprecatedErrInvalidParameter',
    'ErrInvalidVersionName': 'DeprecatedErrInvalidVersionName',
    'ErrInternalBadRequest': 'DeprecatedErrInternalBadRequest',
    'ErrWorkflowCompileFail': 'DeprecatedErrWorkflowCompileFail',
    'ErrConversationNameIsDuplicated': 'DeprecatedErrConversationNameIsDuplicated',
    'ErrConversationOfAppNotFound': 'DeprecatedErrConversationOfAppNotFound',
    'ErrConversationNodeInvalidOperation': 'DeprecatedErrConversationNodeInvalidOperation',
    'ErrChatFlowRoleOperationFail': 'DeprecatedErrChatFlowRoleOperationFail',
    'ErrMessageNodeOperationFail': 'DeprecatedErrMessageNodeOperationFail',
    'ErrNodeOutputParseFail': 'DeprecatedErrNodeOutputParseFail',
    'ErrInputFieldMissing': 'DeprecatedErrInputFieldMissing',
    'ErrLLMStructuredOutputParseFail': 'DeprecatedErrLLMStructuredOutputParseFail',
    'ErrConversationNodesNotAvailable': 'DeprecatedErrConversationNodesNotAvailable',
    'ErrConversationNodeOperationFail': 'DeprecatedErrConversationNodeOperationFail',
    'ErrAuthorizationRequired': 'DeprecatedErrAuthorizationRequired',
    'ErrPluginIDNotFound': 'DeprecatedErrPluginIDNotFound',
    'ErrPluginAPIErr': 'DeprecatedErrPluginAPIErr',
    'ErrTOSError': 'DeprecatedErrTOSError',
    'ErrMissingRequiredParam': 'DeprecatedErrMissingRequiredParam',
    'ErrVariablesAPIFail': 'DeprecatedErrVariablesAPIFail',
}


def fix_file(file_path):
    """修复单个文件中的 errno 使用"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()

    original_content = content
    changes = []

    # 匹配模式：vo.WrapError(errno.ErrXXX, ...)
    # 替换为：vo.WrapError(errno.DeprecatedErrXXX, ...)
    pattern = r'\bvo\.(WrapError|NewError|WrapWarn)\(errno\.(Err[A-Z][a-zA-Z0-9]*)\b'

    def replace_func(match):
        vo_func = match.group(1)  # WrapError, NewError, or WrapWarn
        err_name = match.group(2)  # ErrXXX

        if err_name in ERRNO_MAPPING:
            deprecated_name = ERRNO_MAPPING[err_name]
            changes.append(f"  {err_name} -> {deprecated_name}")
            return f'vo.{vo_func}(errno.{deprecated_name}'
        return match.group(0)

    content = re.sub(pattern, replace_func, content)

    if content != original_content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
        return changes
    return None


def main():
    workflow_dir = 'backend/domain/workflow'

    if not os.path.exists(workflow_dir):
        print(f"错误：目录 {workflow_dir} 不存在")
        sys.exit(1)

    total_files = 0
    total_changes = 0

    # 遍历所有 Go 文件
    for root, dirs, files in os.walk(workflow_dir):
        for file in files:
            if file.endswith('.go'):
                file_path = os.path.join(root, file)
                changes = fix_file(file_path)
                if changes:
                    total_files += 1
                    total_changes += len(changes)
                    print(f"修复文件: {file_path}")
                    for change in changes:
                        print(change)
                    print()

    print(f"\n修复完成！")
    print(f"总共修复文件数: {total_files}")
    print(f"总共修改数: {total_changes}")


if __name__ == '__main__':
    main()
