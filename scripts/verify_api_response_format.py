#!/usr/bin/env python3
"""
API响应格式验证脚本
验证所有Handler文件是否使用统一的响应格式

作者: API规范化工具
日期: 2025-01-01
"""

import os
import re
from pathlib import Path
from collections import defaultdict


def check_file_compliance(file_path: str):
    """检查单个文件的合规性"""
    violations = []
    success_count = 0
    error_count = 0

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            lines = content.split('\n')

        # 检查是否导入了httputil
        has_httputil_import = '"github.com/coze-dev/coze-studio/backend/api/internal/httputil"' in content

        if not has_httputil_import:
            violations.append({
                'type': 'missing_import',
                'severity': 'error',
                'message': '缺少httputil导入'
            })

        # 检查响应格式
        for i, line in enumerate(lines, 1):
            # 检查旧的直接返回格式
            if re.search(r'c\.JSON\(consts\.StatusOK,\s*\w+\)', line):
                violations.append({
                    'type': 'old_success_format',
                    'severity': 'error',
                    'line': i,
                    'code': line.strip(),
                    'message': '使用旧的成功响应格式,应使用httputil.BuildSuccessResp'
                })

            # 检查旧的错误响应格式
            if 'invalidParamRequestResponse(' in line or 'internalServerErrorResponse(' in line:
                violations.append({
                    'type': 'old_error_format',
                    'severity': 'warning',
                    'line': i,
                    'code': line.strip(),
                    'message': '使用旧的错误响应格式,应使用httputil.BuildErrorResp'
                })

            # 统计新格式使用
            if 'httputil.BuildSuccessResp' in line:
                success_count += 1
            if 'httputil.BuildErrorResp' in line:
                error_count += 1

    except Exception as e:
        violations.append({
            'type': 'file_error',
            'severity': 'error',
            'message': f'读取文件失败: {e}'
        })

    return {
        'violations': violations,
        'stats': {
            'success_count': success_count,
            'error_count': error_count,
            'total': success_count + error_count
        }
    }


def main():
    """主函数"""
    print("=" * 80)
    print("API响应格式合规性验证")
    print("=" * 80)

    # 获取项目根目录
    project_root = Path(__file__).parent.parent
    os.chdir(project_root)

    # 检查所有Handler文件
    handler_dir = Path("backend/api/handler/coze")
    handler_files = list(handler_dir.rglob("*.go"))

    # 统计
    total_files = 0
    compliant_files = 0
    total_violations = 0
    total_success_resp = 0
    total_error_resp = 0

    file_results = []

    for file_path in handler_files:
        result = check_file_compliance(str(file_path))
        violations = result['violations']
        stats = result['stats']

        total_files += 1
        total_violations += len(violations)
        total_success_resp += stats['success_count']
        total_error_resp += stats['error_count']

        is_compliant = len([v for v in violations if v['severity'] == 'error']) == 0
        if is_compliant:
            compliant_files += 1

        file_results.append({
            'file': str(file_path),
            'compliant': is_compliant,
            'violations': violations,
            'stats': stats
        })

    # 输出结果
    print("\n📊 总体统计:")
    print(f"  检查文件数: {total_files}")
    print(f"  合规文件数: {compliant_files} ({compliant_files/total_files*100:.1f}%)")
    print(f"  违规总数: {total_violations}")
    print(f"  成功响应: {total_success_resp}")
    print(f"  错误响应: {total_error_resp}")
    print(f"  总响应数: {total_success_resp + total_error_resp}")

    # 输出合规性评分
    compliance_rate = (compliant_files / total_files * 100) if total_files > 0 else 0
    print(f"\n🎯 合规性评分: {compliance_rate:.1f}%")

    if compliance_rate >= 95:
        print("  等级: ⭐⭐⭐⭐⭐ 优秀")
    elif compliance_rate >= 85:
        print("  等级: ⭐⭐⭐⭐ 良好")
    elif compliance_rate >= 70:
        print("  等级: ⭐⭐⭐ 合格")
    else:
        print("  等级: ⭐⭐ 需改进")

    # 输出违规清单
    if total_violations > 0:
        print("\n⚠️  违规清单:")

        # 按类型分组
        by_type = defaultdict(list)
        for result in file_results:
            for v in result['violations']:
                by_type[v['type']].append({
                    'file': result['file'],
                    'violation': v
                })

        for vtype, items in sorted(by_type.items()):
            print(f"\n{vtype} ({len(items)}处):")
            for item in items[:5]:  # 只显示前5个
                v = item['violation']
                print(f"  - {item['file']}")
                if 'line' in v:
                    print(f"    Line {v['line']}: {v['message']}")
                else:
                    print(f"    {v['message']}")
            if len(items) > 5:
                print(f"  ... 还有{len(items)-5}处")

    # 输出文件详情(仅显示不合规的文件)
    non_compliant = [r for r in file_results if not r['compliant']]
    if non_compliant:
        print(f"\n❌ 不合规文件 ({len(non_compliant)}个):")
        for result in non_compliant[:10]:
            print(f"  - {result['file']}")
            print(f"    违规数: {len(result['violations'])}")
        if len(non_compliant) > 10:
            print(f"  ... 还有{len(non_compliant)-10}个文件")

    print("\n" + "=" * 80)

    # 判断是否通过验证
    if total_violations == 0:
        print("✅ 验证通过!所有文件都使用统一的响应格式!")
        return 0
    elif compliance_rate >= 90:
        print("⚠️  基本通过,仍有少数违规需要修复")
        return 0
    else:
        print("❌ 验证失败,需要继续修复响应格式")
        return 1


if __name__ == "__main__":
    exit(main())
