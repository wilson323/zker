#!/usr/bin/env python3
"""
HTTP方法误用分析脚本
识别所有使用POST进行查询操作的API

目标: 将POST查询改为GET,符合RESTful规范

作者: API规范化工具
日期: 2025-01-01
"""

import re
from pathlib import Path
from collections import defaultdict

# POST查询违规模式
VIOLATION_PATTERNS = [
    (r'\.POST\("/([^"]*)"', "POST方法"),
]

# 查询操作关键词
QUERY_KEYWORDS = ['get_', 'list', 'detail', 'search', 'query', 'find']


def analyze_router_file(file_path: str):
    """分析路由文件中的HTTP方法使用"""
    violations = []

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
            lines = content.split('\n')

        for i, line in enumerate(lines, 1):
            # 检测POST方法
            if '.POST(' in line:
                # 提取路径
                match = re.search(r'\.POST\("([^"]+)"', line)
                if match:
                    path = match.group(1)

                    # 检查是否是查询操作
                    is_query = any(kw in path.lower() for kw in QUERY_KEYWORDS)

                    if is_query:
                        # 提取handler函数名
                        handler_match = re.search(r',\s*(\w+)\)', line)
                        handler = handler_match.group(1) if handler_match else "Unknown"

                        violations.append({
                            'line': i,
                            'path': path,
                            'handler': handler,
                            'current': f'POST /{path}',
                            'recommended': f'GET /{path}',
                            'code': line.strip()
                        })

    except Exception as e:
        print(f"❌ 分析失败 {file_path}: {e}")

    return violations


def generate_fix_guide(violations):
    """生成修复指南"""
    if not violations:
        print("✅ 未发现HTTP方法误用问题")
        return

    print("=" * 100)
    print(f"⚠️  发现 {len(violations)} 处HTTP方法误用")
    print("=" * 100)

    # 按路径分组
    by_path = defaultdict(list)
    for v in violations:
        by_path[v['path']].append(v)

    print("\n📋 违规清单:\n")

    for i, v in enumerate(violations, 1):
        print(f"{i}. {v['current']} → {v['handler']}()")
        print(f"   应改为: {v['recommended']}")
        print(f"   位置: Line {v['line']}")
        print(f"   代码: {v['code']}")
        print()

    print("=" * 100)
    print("🔧 修复建议:")
    print("=" * 100)
    print("""
对于查询操作,应该使用GET方法而不是POST:

修复步骤:
1. 修改HTTP方法: .POST() -> .GET()
2. 修改Handler函数: 从c.BindAndValidate(&req)改为从Query参数读取
3. 更新API文档: 更新Swagger注释

示例:

修改前:
    _bot.POST("/get_type_list", GetTypeList)

    func GetTypeList(ctx context.Context, c *app.RequestContext) {
        var req GetTypeListRequest
        c.BindAndValidate(&req)  // 从body读取
        // ...
    }

修改后:
    _bot.GET("/types", GetTypeList)

    func GetTypeList(ctx context.Context, c *app.RequestContext) {
        // 从query参数读取
        page := c.Query("page", "1")
        pageSize := c.Query("page_size", "20")
        // ...
    }
""")


def main():
    """主函数"""
    print("=" * 100)
    print("HTTP方法误用分析工具")
    print("=" * 100)

    # 分析主要路由文件
    router_file = Path("backend/api/router/coze/api.go")

    if not router_file.exists():
        print(f"❌ 路由文件不存在: {router_file}")
        return

    violations = analyze_router_file(router_file)

    # 生成修复指南
    generate_fix_guide(violations)

    # 输出CSV格式供进一步处理
    if violations:
        csv_file = Path("scripts/http_method_violations.csv")
        with open(csv_file, 'w', encoding='utf-8') as f:
            f.write("路径,Handler,当前方法,推荐方法,代码行\n")
            for v in violations:
                f.write(f"{v['path']},{v['handler']},POST,GET,\"{v['code']}\"\n")
        print(f"\n📊 违规清单已导出: {csv_file}")


if __name__ == "__main__":
    main()
