#!/usr/bin/env python3
"""
ZKER API接口实现完整性分析工具
分析33个API文档,检查后端和前端的实现情况
"""

import os
import re
import json
from pathlib import Path
from typing import Dict, List, Tuple, Set
from dataclasses import dataclass, field
from collections import defaultdict

@dataclass
class APIEndpoint:
    """API端点信息"""
    method: str
    path: str
    description: str = ""
    implemented: bool = False
    backend_file: str = ""
    backend_function: str = ""
    frontend_hook: str = ""
    frontend_component: str = ""

@dataclass
class APIModule:
    """API模块信息"""
    name: str
    doc_file: str
    priority: str
    endpoints: List[APIEndpoint] = field(default_factory=list)
    implementation_rate: float = 0.0

class APIAnalyzer:
    """API实现分析器"""

    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.api_docs_dir = self.project_root / "docs/企业级功能完善与统一性设计方案"
        self.backend_handler_dir = self.project_root / "backend/api/handler/coze"
        self.backend_router_dir = self.project_root / "backend/api/router"
        self.frontend_api_dir = self.project_root / "frontend/packages/api-client/src"
        self.frontend_hooks_dir = self.frontend_api_dir / "hooks"

        self.modules: Dict[str, APIModule] = {}
        self.backend_functions: Dict[str, Set[str]] = defaultdict(set)
        self.router_registrations: Dict[str, List[str]] = defaultdict(list)

    def extract_api_from_doc(self, doc_file: Path) -> APIModule:
        """从API文档中提取API端点信息"""
        print(f"正在分析: {doc_file.name}")

        with open(doc_file, 'r', encoding='utf-8') as f:
            content = f.read()

        # 提取模块名称
        module_name_match = re.search(r'# API接口文档[：:]\s*(.+?)\s*\n', content)
        module_name = module_name_match.group(1) if module_name_match else doc_file.stem

        # 提取优先级
        priority_match = re.search(r'\*\*优先级\*\*\s*:\s*(P\d)', content)
        priority = priority_match.group(1) if priority_match else "P3"

        module = APIModule(
            name=module_name,
            doc_file=str(doc_file),
            priority=priority
        )

        # 提取API端点
        # 格式1: **接口地址**: `GET /api/xxx`
        pattern1 = r'\*\*接口地址\*\*\s*:\s*`([A-Z]+)\s+([^`]+)`'
        # 格式2: **请求方式**: `GET /api/xxx`
        pattern2 = r'\*\*请求方式\*\*\s*:\s*`([A-Z]+)\s+([^`]+)`'

        for match in re.finditer(pattern1, content):
            method, path = match.groups()
            module.endpoints.append(APIEndpoint(method=method, path=path.strip()))

        for match in re.finditer(pattern2, content):
            method, path = match.groups()
            # 避免重复
            if not any(ep.method == method and ep.path == path.strip() for ep in module.endpoints):
                module.endpoints.append(APIEndpoint(method=method, path=path.strip()))

        print(f"  发现 {len(module.endpoints)} 个API端点")
        return module

    def scan_backend_handlers(self):
        """扫描后端Handler文件"""
        print("\n扫描后端Handler...")

        if not self.backend_handler_dir.exists():
            print(f"警告: Handler目录不存在: {self.backend_handler_dir}")
            return

        for handler_file in self.backend_handler_dir.glob("*.go"):
            try:
                content = handler_file.read_text(encoding='utf-8', errors='ignore')

                # 提取函数定义
                func_pattern = r'func\s+(\w+)\s*\([^)]*\)\s*'
                for match in re.finditer(func_pattern, content):
                    func_name = match.group(1)
                    if not func_name.startswith('get') or func_name == 'getDigitalEmployeeSvc':
                        self.backend_functions[handler_file.name].add(func_name)

            except Exception as e:
                print(f"  警告: 读取{handler_file}失败: {e}")

    def scan_router_registrations(self):
        """扫描路由注册"""
        print("\n扫描路由注册...")

        if not self.backend_router_dir.exists():
            print(f"警告: Router目录不存在: {self.backend_router_dir}")
            return

        for router_file in self.backend_router_dir.rglob("*.go"):
            try:
                content = router_file.read_text(encoding='utf-8', errors='ignore')

                # 提取路由注册
                # 格式: _route.POST("/path", handler)
                route_pattern = r'_(GET|POST|PUT|DELETE|PATCH)\(["\']([^"\']+)["\']'
                for match in re.finditer(route_pattern, content):
                    method, path = match.groups()
                    self.router_registrations[str(router_file)].append(f"{method} {path}")

            except Exception as e:
                print(f"  警告: 读取{router_file}失败: {e}")

    def check_backend_implementation(self, endpoint: APIEndpoint) -> bool:
        """检查后端实现 - 检查Handler和路由"""
        # 方法1: 检查路由注册
        path_normalized = endpoint.path.replace('{', ':').replace('}', '')

        for file, routes in self.router_registrations.items():
            for route in routes:
                route_method, route_path = route.split(' ', 1)
                if route_method == endpoint.method:
                    # 精确匹配或模式匹配
                    if route_path == path_normalized or self.path_matches(route_path, path_normalized):
                        endpoint.backend_file = file
                        endpoint.implemented = True
                        return True

        # 方法2: 检查Handler函数名
        # 从路径提取关键词,如 /api/bot-store/listings -> bot_store, listings
        path_parts = [p for p in endpoint.path.split('/') if p and not p.startswith('{')]
        if len(path_parts) >= 2:
            # 尝试匹配handler文件名和函数
            for handler_file, funcs in self.backend_functions.items():
                for func in funcs:
                    func_lower = func.lower()
                    # 检查函数名是否包含路径关键词
                    if any(part.replace('-', '_').lower() in func_lower for part in path_parts[-2:]):
                        endpoint.backend_file = handler_file
                        endpoint.backend_function = func
                        endpoint.implemented = True
                        return True

        return False

    @staticmethod
    def path_matches(pattern: str, path: str) -> bool:
        """检查路径是否匹配模式"""
        pattern_parts = pattern.split('/')
        path_parts = path.split('/')

        if len(pattern_parts) != len(path_parts):
            return False

        for p, actual in zip(pattern_parts, path_parts):
            if p.startswith(':') or p.startswith('{'):
                continue
            if p != actual:
                return False

        return True

    def analyze(self):
        """执行完整分析"""
        print("=" * 80)
        print("ZKER API接口实现完整性分析")
        print("=" * 80)

        # 1. 扫描所有API文档
        print("\n[步骤1] 扫描API文档...")
        api_files = sorted(self.api_docs_dir.glob("API接口文档_*.md"))
        print(f"找到 {len(api_files)} 个API文档")

        for api_file in api_files:
            module = self.extract_api_from_doc(api_file)
            if module.endpoints:
                self.modules[module.name] = module

        # 2. 扫描后端代码
        print("\n[步骤2] 扫描后端实现...")
        self.scan_backend_handlers()
        self.scan_router_registrations()

        print(f"  Handler文件: {len(self.backend_functions)}")
        print(f"  路由注册: {sum(len(routes) for routes in self.router_registrations.values())} 条")

        # 3. 检查实现情况
        print("\n[步骤3] 检查实现情况...")
        for module in self.modules.values():
            implemented_count = 0
            for endpoint in module.endpoints:
                if self.check_backend_implementation(endpoint):
                    implemented_count += 1

            module.implementation_rate = implemented_count / len(module.endpoints) if module.endpoints else 0

        # 4. 生成报告
        print("\n[步骤4] 生成分析报告...")
        self.generate_report()

    def generate_report(self):
        """生成分析报告"""
        report = []
        report.append("# API接口实现完整性分析报告\n")
        report.append(f"**生成时间**: {__import__('datetime').datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")
        report.append("---\n")

        # 总体统计
        total_modules = len(self.modules)
        total_endpoints = sum(len(m.endpoints) for m in self.modules.values())
        implemented_endpoints = sum(sum(1 for e in m.endpoints if e.implemented) for m in self.modules.values())
        overall_rate = implemented_endpoints / total_endpoints if total_endpoints > 0 else 0

        report.append("## 总体统计\n")
        report.append(f"- **API模块数**: {total_modules}")
        report.append(f"- **设计API总数**: {total_endpoints}")
        report.append(f"- **已实现API数**: {implemented_endpoints}")
        report.append(f"- **未实现API数**: {total_endpoints - implemented_endpoints}")
        report.append(f"- **整体完成度**: {overall_rate * 100:.1f}%\n")

        # 分模块统计
        report.append("## 分模块统计\n")
        report.append("| 模块 | 优先级 | API数 | 已实现 | 未实现 | 完成度 |\n")
        report.append("|------|--------|-------|--------|--------|--------|\n")

        # 按优先级排序
        sorted_modules = sorted(
            self.modules.values(),
            key=lambda m: (m.priority, m.name)
        )

        for module in sorted_modules:
            implemented = sum(1 for e in module.endpoints if e.implemented)
            not_implemented = len(module.endpoints) - implemented
            rate = module.implementation_rate * 100

            report.append(
                f"| {module.name} | {module.priority} | {len(module.endpoints)} | "
                f"{implemented} | {not_implemented} | {rate:.0f}% |\n"
            )

        # 未实现API详情
        report.append("\n## 未实现API详情\n")

        for module in sorted_modules:
            not_implemented = [e for e in module.endpoints if not e.implemented]
            if not_implemented:
                report.append(f"### {module.name} ({module.priority})\n\n")
                for endpoint in not_implemented:
                    report.append(f"- **{endpoint.method}** `{endpoint.path}`\n")
                report.append("\n")

        # P0/P1优先级未实现API
        report.append("## P0/P1 优先级未实现API\n")

        critical_apis = []
        for module in sorted_modules:
            if module.priority in ['P0', 'P1']:
                for endpoint in module.endpoints:
                    if not endpoint.implemented:
                        critical_apis.append((module, endpoint))

        if critical_apis:
            report.append(f"**共 {len(critical_apis)} 个关键API未实现**\n\n")
            for module, endpoint in critical_apis:
                report.append(f"#### {module.name} ({module.priority})\n")
                report.append(f"- `{endpoint.method} {endpoint.path}`\n\n")
        else:
            report.append("✅ 所有P0/P1 API已实现!\n")

        # 实施建议
        report.append("\n## 实施建议\n")

        # 统计P0/P1未实现
        p0_p1_modules = [m for m in sorted_modules if m.priority in ['P0', 'P1']]
        p0_p1_total = sum(len(m.endpoints) for m in p0_p1_modules)
        p0_p1_implemented = sum(sum(1 for e in m.endpoints if e.implemented) for m in p0_p1_modules)
        p0_p1_rate = p0_p1_implemented / p0_p1_total * 100 if p0_p1_total > 0 else 0

        report.append(f"### 第一阶段: P0/P1核心API ({p0_p1_total}个API, 当前完成度{p0_p1_rate:.0f}%)\n")
        report.append("优先实现以下模块:\n")
        for module in p0_p1_modules:
            not_implemented = [e for e in module.endpoints if not e.implemented]
            if not_implemented:
                report.append(f"- **{module.name}**: {len(not_implemented)}个API未实现\n")

        # 保存报告
        report_file = self.project_root / "API_IMPLEMENTATION_ANALYSIS_REPORT.md"
        report_file.write_text(''.join(report), encoding='utf-8')

        print(f"\n报告已生成: {report_file}")

        # 输出摘要
        print("\n" + "=" * 80)
        print("分析摘要")
        print("=" * 80)
        print(f"总API数: {total_endpoints}")
        print(f"已实现: {implemented_endpoints} ({overall_rate * 100:.1f}%)")
        print(f"未实现: {total_endpoints - implemented_endpoints}")
        print(f"\nP0/P1 API完成度: {p0_p1_rate:.1f}%")

        # Top 5未实现最多的模块
        print("\n未实现API最多的模块:")
        top_unimplemented = sorted(
            self.modules.values(),
            key=lambda m: sum(1 for e in m.endpoints if not e.implemented),
            reverse=True
        )[:5]

        for i, module in enumerate(top_unimplemented, 1):
            not_implemented = sum(1 for e in module.endpoints if not e.implemented)
            print(f"  {i}. {module.name}: {not_implemented}个未实现")

if __name__ == "__main__":
    project_root = Path(__file__).parent.parent
    analyzer = APIAnalyzer(str(project_root))
    analyzer.analyze()
