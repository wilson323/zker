#!/usr/bin/env python3
"""
修复 org/service 包中的 errno 返回错误
将 `return nil, errno.ErrXxx` 替换为 `return nil, errorx.NewByErrorCode(errno.ErrXxx)`
"""

import re
import os

SERVICE_DIR = "backend/domain/org/service"

# 需要替换的模式列表
PATTERNS = [
    (r'\breturn nil, errorx\.New\(errno\.Err([A-Z][a-zA-Z]+)\)', r'return nil, errorx.NewByErrorCode(errno.Err\1)'),
    (r'\breturn nil, errorx\.New\(errno\.Err([A-Z][a-zA-Z]+)\)', r'return nil, errorx.NewByErrorCode(errno.Err\1)'),
    (r'\breturn , errorx\.New\(errno\.Err([A-Z][a-zA-Z]+)\)', r'return , errorx.NewByErrorCode(errno.Err\1)'),
]

def fix_file(filepath):
    """修复单个文件"""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original_content = content

    # 应用所有替换规则
    for pattern, replacement in PATTERNS:
        content = re.sub(pattern, replacement, content)

    # 如果内容有变化，写回文件
    if content != original_content:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed: {filepath}")
        return True
    return False

def main():
    """主函数"""
    fixed_count = 0

    for filename in os.listdir(SERVICE_DIR):
        if filename.endswith('.go'):
            filepath = os.path.join(SERVICE_DIR, filename)
            if fix_file(filepath):
                fixed_count += 1

    print(f"\n总计修复 {fixed_count} 个文件")

if __name__ == '__main__':
    main()
