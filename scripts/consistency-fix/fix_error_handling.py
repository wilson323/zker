#!/usr/bin/env python3
"""
ZKER 错误处理统一修复脚本
用途: 批量替换硬编码错误为统一错误码
作者: ZKER开发团队
更新: 2025-01-03
"""

import os
import re
import sys
from pathlib import Path
from typing import List, Tuple

# 颜色定义
class Colors:
    RED = '\033[0;31m'
    GREEN = '\033[0;32m'
    YELLOW = '\033[1;33m'
    BLUE = '\033[0;34m'
    NC = '\033[0m'

# 错误码映射（示例，需要根据实际情况补充）
ERROR_CODE_MAPPING = {
    # 组织管理
    r'organization not found': 'errno.OrgNotFound',
    r'organization already exists': 'errno.ErrOrgAlreadyExists',
    r'invalid organization': 'errno.ErrInvalidOrg',

    # 租户管理
    r'tenant not found': 'errno.ErrTenantNotFound',
    r'tenant already exists': 'errno.ErrTenantAlreadyExists',

    # 权限管理
    r'permission denied': 'errno.ErrPermissionDenied',
    r'role not found': 'errno.ErrRoleNotFound',

    # 工作流
    r'workflow not found': 'errno.ErrWorkflowNotFound',
    r'node not found': 'errno.ErrNodeNotFound',
    r'invalid workflow state': 'errno.ErrInvalidWorkflowState',
    r'node execution failed': 'errno.ErrNodeExecutionFailed',

    # 计费
    r'invoice not found': 'errno.ErrInvoiceNotFound',
    r'invalid amount': 'errno.ErrInvalidAmount',
    r'payment failed': 'errno.ErrPaymentFailed',

    # 通用
    r'database error': 'errno.ErrDatabase',
    r'invalid parameter': 'errno.ErrInvalidParam',
    r'internal server error': 'errno.ErrInternalServer',
}

class ErrorFixer:
    def __init__(self, filepath: str):
        self.filepath = Path(filepath)
        self.content = ""
        self.fixed_content = ""
        self.fixes_count = 0
        self.backup_path = filepath + ".backup"

    def read_file(self) -> bool:
        """读取文件内容"""
        try:
            with open(self.filepath, 'r', encoding='utf-8') as f:
                self.content = f.read()
            self.fixed_content = self.content
            return True
        except Exception as e:
            print(f"{Colors.RED}读取文件失败: {e}{Colors.NC}")
            return False

    def backup_file(self) -> bool:
        """备份原文件"""
        try:
            with open(self.backup_path, 'w', encoding='utf-8') as f:
                f.write(self.content)
            return True
        except Exception as e:
            print(f"{Colors.RED}备份文件失败: {e}{Colors.NC}")
            return False

    def fix_errors_new(self) -> int:
        """修复 errors.New() 调用"""
        pattern = r'errors\.New\("([^"]+)"\)'
        matches = list(re.finditer(pattern, self.fixed_content))

        for match in reversed(matches):  # 反向遍历，避免位置错乱
            error_msg = match.group(1)
            replacement = self.find_error_code(error_msg)

            if replacement:
                old_text = match.group(0)
                new_text = f'errorx.New({replacement}, "{error_msg}")'
                self.fixed_content = self.fixed_content[:match.start()] + new_text + self.fixed_content[match.end():]
                self.fixes_count += 1
                print(f"  {Colors.GREEN}✓{Colors.NC} {old_text} → {new_text}")

        return len(matches)

    def fix_fmt_errorf(self) -> int:
        """修复 fmt.Errorf() 调用"""
        pattern = r'fmt\.Errorf\("([^"]+)",\s*([^)]+)\)'
        matches = list(re.finditer(pattern, self.fixed_content))

        for match in reversed(matches):
            error_msg = match.group(1)
            args = match.group(2)
            replacement = self.find_error_code(error_msg)

            if replacement:
                old_text = match.group(0)
                # 检查是否有格式化参数
                if '%s' in error_msg or '%d' in error_msg or '%v' in error_msg:
                    # 有格式化参数
                    new_text = f'errorx.New({replacement}, "{error_msg}", {args})'
                else:
                    # 无格式化参数
                    new_text = f'errorx.New({replacement}, "{error_msg}")'

                self.fixed_content = self.fixed_content[:match.start()] + new_text + self.fixed_content[match.end():]
                self.fixes_count += 1
                print(f"  {Colors.GREEN}✓{Colors.NC} {old_text} → {new_text}")

        return len(matches)

    def find_error_code(self, error_msg: str) -> str:
        """根据错误信息查找对应的错误码"""
        error_msg_lower = error_msg.lower()

        for pattern, error_code in ERROR_CODE_MAPPING.items():
            if re.search(pattern, error_msg_lower):
                return error_code

        # 未找到匹配的错误码
        return None

    def write_file(self) -> bool:
        """写入修复后的内容"""
        try:
            with open(self.filepath, 'w', encoding='utf-8') as f:
                f.write(self.fixed_content)
            return True
        except Exception as e:
            print(f"{Colors.RED}写入文件失败: {e}{Colors.NC}")
            return False

    def fix(self) -> bool:
        """执行修复流程"""
        print(f"\n{Colors.BLUE}修复文件: {self.filepath}{Colors.NC}")

        if not self.read_file():
            return False

        if not self.backup_file():
            return False

        print(f"{Colors.YELLOW}查找硬编码错误...{Colors.NC}")

        self.fix_errors_new()
        self.fix_fmt_errorf()

        if self.fixes_count > 0:
            print(f"{Colors.GREEN}共修复 {self.fixes_count} 处{Colors.NC}")

            if self.write_file():
                print(f"{Colors.GREEN}✅ 修复完成{Colors.NC}")
                return True
            else:
                return False
        else:
            print(f"{Colors.YELLOW}无需修复{Colors.NC}")
            return True

def scan_directory(directory: str) -> List[Path]:
    """扫描目录，找出需要修复的文件"""
    dir_path = Path(directory)
    go_files = []

    for file in dir_path.rglob("*.go"):
        # 跳过测试文件
        if "test.go" in file.name:
            continue

        # 检查是否包含硬编码错误
        try:
            with open(file, 'r', encoding='utf-8') as f:
                content = f.read()
                if 'errors.New' in content or 'fmt.Errorf' in content:
                    go_files.append(file)
        except Exception:
            continue

    return go_files

def main():
    if len(sys.argv) < 2:
        print(f"{Colors.YELLOW}用法: python3 fix_error_handling.py <file|directory>{Colors.NC}")
        print("")
        print("示例:")
        print("  python3 fix_error_handling.py backend/domain/billing/service/billing_engine.go")
        print("  python3 fix_error_handling.py backend/domain/workflow/")
        print("")
        print("注意: 修复前会自动备份文件 (.backup 后缀)")
        sys.exit(1)

    target = sys.argv[1]
    target_path = Path(target)

    print(f"{Colors.BLUE}========================================{Colors.NC}")
    print(f"{Colors.BLUE}错误处理统一修复工具{Colors.NC}")
    print(f"{Colors.BLUE}========================================{Colors.NC}")
    print("")

    total_files = 0
    fixed_files = 0
    skipped_files = 0

    if target_path.is_file():
        # 单个文件
        total_files = 1
        fixer = ErrorFixer(str(target_path))
        if fixer.fix():
            if fixer.fixes_count > 0:
                fixed_files += 1
            else:
                skipped_files += 1

    elif target_path.is_dir():
        # 目录
        print(f"{Colors.BLUE}扫描目录: {target_path}{Colors.NC}")
        print("")

        go_files = scan_directory(str(target_path))
        total_files = len(go_files)

        print(f"{Colors.BLUE}找到 {total_files} 个需要检查的文件{Colors.NC}")
        print("")

        for file in go_files:
            fixer = ErrorFixer(str(file))
            if fixer.fix():
                if fixer.fixes_count > 0:
                    fixed_files += 1
                else:
                    skipped_files += 1

    else:
        print(f"{Colors.RED}❌ 错误: {target} 不是文件或目录{Colors.NC}")
        sys.exit(1)

    # 汇总
    print("")
    print(f"{Colors.BLUE}========================================{Colors.NC}")
    print(f"{Colors.BLUE}修复完成{Colors.NC}")
    print(f"{Colors.BLUE}========================================{Colors.NC}")
    print("")
    print(f"扫描文件: {total_files}")
    print(f"{Colors.GREEN}已修复:   {fixed_files}{Colors.NC}")
    print(f"{Colors.YELLOW}跳过:     {skipped_files}{Colors.NC}")
    print("")

    if fixed_files > 0:
        print(f"{Colors.YELLOW}⚠️  备份文件已创建 (.backup 后缀){Colors.NC}")
        print(f"{Colors.YELLOW}   请确认修复无误后删除备份文件{Colors.NC}")
        print("")
        print("删除备份命令:")
        print(f"  find \"{target}\" -name '*.backup' -delete")
        print("")

    print(f"{Colors.BLUE}下一步:{Colors.NC}")
    print("  1. 运行测试: cd backend && go test ./...")
    print("  2. 代码审查: git diff")
    print("  3. 提交修复: git commit -am 'fix: 统一错误处理'")
    print("")

if __name__ == "__main__":
    main()
