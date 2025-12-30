#!/usr/bin/env python3
"""
组织中心性能测试结果可视化工具

功能:
- 解析K6 JSON测试结果
- 生成性能趋势图
- 对比基线数据
- 输出HTML报告

使用方法:
    # 生成单个测试结果图表
    python plot_results.py results/base_load_test.json

    # 生成所有测试结果图表
    python plot_results.py results/

    # 对比两个测试结果
    python plot_results.py results/before.json results/after.json --compare

@author 研发B (后端工程师)
@version 1.0
@date 2025-01-01
"""

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Dict, List, Tuple

import matplotlib.pyplot as plt
import matplotlib
matplotlib.use('Agg')  # 非交互式后端

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['SimHei', 'Arial Unicode MS']
plt.rcParams['axes.unicode_minus'] = False


class PerformanceAnalyzer:
    """性能数据分析器"""

    def __init__(self, json_file: str):
        self.json_file = json_file
        self.data = self._load_json()

    def _load_json(self) -> Dict:
        """加载JSON结果文件"""
        with open(self.json_file, 'r', encoding='utf-8') as f:
            return json.load(f)

    def get_response_times(self) -> List[float]:
        """获取所有请求的响应时间"""
        return [
            item['data']['values']['http_req_duration']
            for item in self.data
            if 'http_req_duration' in item.get('data', {}).get('values', {})
        ]

    def get_throughput(self) -> float:
        """计算吞吐量 (requests per second)"""
        total_requests = len([item for item in self.data if item.get('type') == 'Point'])

        # 获取测试持续时间
        timestamps = [
            item['data']['time']
            for item in self.data
            if item.get('type') == 'Point'
        ]

        if not timestamps:
            return 0.0

        duration = (max(timestamps) - min(timestamps)) / 1000  # 转换为秒
        return total_requests / duration if duration > 0 else 0.0

    def get_error_rate(self) -> float:
        """计算错误率"""
        total_requests = len([item for item in self.data if item.get('type') == 'Point'])
        failed_requests = len([
            item for item in self.data
            if item.get('type') == 'Point'
            and not item.get('data', {}).get('tags', {}).get('success', True)
        ])

        if total_requests == 0:
            return 0.0

        return (failed_requests / total_requests) * 100

    def get_percentiles(self) -> Dict[str, float]:
        """计算响应时间百分位数"""
        response_times = self.get_response_times()

        if not response_times:
            return {'p50': 0, 'p90': 0, 'p95': 0, 'p99': 0}

        sorted_times = sorted(response_times)
        total = len(sorted_times)

        return {
            'p50': sorted_times[int(total * 0.50)],
            'p90': sorted_times[int(total * 0.90)],
            'p95': sorted_times[int(total * 0.95)],
            'p99': sorted_times[int(total * 0.99)],
        }

    def get_metrics_by_endpoint(self) -> Dict[str, Dict]:
        """按API端点分组统计"""
        endpoints = {}

        for item in self.data:
            if item.get('type') != 'Point':
                continue

            endpoint = item.get('data', {}).get('tags', {}).get('name', 'unknown')

            if endpoint not in endpoints:
                endpoints[endpoint] = {
                    'response_times': [],
                    'errors': 0,
                    'total': 0,
                }

            duration = item.get('data', {}).get('values', {}).get('http_req_duration', 0)
            is_success = item.get('data', {}).get('tags', {}).get('success', True)

            endpoints[endpoint]['response_times'].append(duration)
            endpoints[endpoint]['total'] += 1

            if not is_success:
                endpoints[endpoint]['errors'] += 1

        # 计算统计指标
        result = {}
        for endpoint, data in endpoints.items():
            times = sorted(data['response_times'])
            total = data['total']

            result[endpoint] = {
                'avg': sum(times) / total if total > 0 else 0,
                'p50': times[int(total * 0.50)] if total > 0 else 0,
                'p95': times[int(total * 0.95)] if total > 0 else 0,
                'p99': times[int(total * 0.99)] if total > 0 else 0,
                'error_rate': (data['errors'] / total * 100) if total > 0 else 0,
                'total_requests': total,
            }

        return result


def plot_response_times(analyzer: PerformanceAnalyzer, output_file: str):
    """绘制响应时间分布图"""
    response_times = analyzer.get_response_times()

    plt.figure(figsize=(12, 6))

    # 响应时间直方图
    plt.subplot(1, 2, 1)
    plt.hist(response_times, bins=50, color='steelblue', alpha=0.7, edgecolor='black')
    plt.xlabel('响应时间 (ms)')
    plt.ylabel('频次')
    plt.title('响应时间分布')
    plt.grid(True, alpha=0.3)

    # 响应时间箱线图
    plt.subplot(1, 2, 2)
    plt.boxplot(response_times, vert=True, patch_artist=True)
    plt.ylabel('响应时间 (ms)')
    plt.title('响应时间箱线图')
    plt.grid(True, alpha=0.3)

    plt.tight_layout()
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    plt.close()

    print(f"✅ 响应时间图表已保存: {output_file}")


def plot_throughput_over_time(analyzer: PerformanceAnalyzer, output_file: str):
    """绘制吞吐量随时间变化图"""
    # 提取时间序列数据
    timestamps = []
    durations = []

    for item in analyzer.data:
        if item.get('type') != 'Point':
            continue

        timestamp = item.get('data', {}).get('time', 0)
        duration = item.get('data', {}).get('values', {}).get('http_req_duration', 0)

        timestamps.append(timestamp / 1000)  # 转换为秒
        durations.append(duration)

    if not timestamps:
        print("⚠️  没有时间序列数据")
        return

    plt.figure(figsize=(12, 6))

    # 响应时间随时间变化
    plt.plot(timestamps, durations, color='steelblue', alpha=0.6, linewidth=1)
    plt.xlabel('时间 (s)')
    plt.ylabel('响应时间 (ms)')
    plt.title('响应时间随时间变化')
    plt.grid(True, alpha=0.3)

    # 添加移动平均线
    if len(durations) > 10:
        window_size = max(10, len(durations) // 20)
        moving_avg = [
            sum(durations[max(0, i - window_size):i + 1]) /
            min(i + 1, window_size)
            for i in range(len(durations))
        ]
        plt.plot(timestamps, moving_avg, color='red', linewidth=2,
                label=f'移动平均 (窗口={window_size})')
        plt.legend()

    plt.tight_layout()
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    plt.close()

    print(f"✅ 吞吐量图表已保存: {output_file}")


def plot_endpoint_comparison(analyzer: PerformanceAnalyzer, output_file: str):
    """绘制不同API端点性能对比图"""
    endpoints = analyzer.get_metrics_by_endpoint()

    if not endpoints:
        print("⚠️  没有端点数据")
        return

    # 准备数据
    endpoint_names = list(endpoints.keys())
    p95_times = [endpoints[e]['p95'] for e in endpoint_names]
    error_rates = [endpoints[e]['error_rate'] for e in endpoint_names]

    # 创建图表
    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 10))

    # p95响应时间对比
    bars1 = ax1.bar(endpoint_names, p95_times, color='steelblue', alpha=0.7)
    ax1.set_ylabel('p95 响应时间 (ms)')
    ax1.set_title('API端点p95响应时间对比')
    ax1.grid(True, alpha=0.3, axis='y')

    # 添加数值标签
    for bar in bars1:
        height = bar.get_height()
        ax1.text(bar.get_x() + bar.get_width() / 2., height,
                f'{height:.1f}', ha='center', va='bottom', fontsize=8)

    plt.setp(ax1.xaxis.get_majorticklabels(), rotation=45, ha='right')

    # 错误率对比
    bars2 = ax2.bar(endpoint_names, error_rates, color='coral', alpha=0.7)
    ax2.set_ylabel('错误率 (%)')
    ax2.set_title('API端点错误率对比')
    ax2.grid(True, alpha=0.3, axis='y')

    # 添加数值标签
    for bar in bars2:
        height = bar.get_height()
        ax2.text(bar.get_x() + bar.get_width() / 2., height,
                f'{height:.2f}%', ha='center', va='bottom', fontsize=8)

    plt.setp(ax2.xaxis.get_majorticklabels(), rotation=45, ha='right')

    plt.tight_layout()
    plt.savefig(output_file, dpi=300, bbox_inches='tight')
    plt.close()

    print(f"✅ 端点对比图表已保存: {output_file}")


def generate_html_report(analyzers: List[Tuple[str, PerformanceAnalyzer]], output_file: str):
    """生成HTML性能报告"""
    html_template = """
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>组织中心性能测试报告</title>
    <style>
        body {{
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }}
        .container {{
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }}
        h1 {{
            color: #333;
            border-bottom: 3px solid #4CAF50;
            padding-bottom: 10px;
        }}
        h2 {{
            color: #555;
            margin-top: 30px;
        }}
        table {{
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }}
        th, td {{
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }}
        th {{
            background-color: #4CAF50;
            color: white;
            font-weight: bold;
        }}
        tr:hover {{
            background-color: #f5f5f5;
        }}
        .metric {{
            display: inline-block;
            margin: 10px 20px 10px 0;
            padding: 15px;
            background-color: #f9f9f9;
            border-left: 4px solid #4CAF50;
            border-radius: 4px;
        }}
        .metric-label {{
            font-size: 14px;
            color: #666;
        }}
        .metric-value {{
            font-size: 24px;
            font-weight: bold;
            color: #333;
        }}
        .status-pass {{
            color: #4CAF50;
        }}
        .status-fail {{
            color: #f44336;
        }}
        .chart {{
            margin: 30px 0;
            text-align: center;
        }}
        .chart img {{
            max-width: 100%;
            height: auto;
            border: 1px solid #ddd;
            border-radius: 4px;
        }}
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 组织中心性能测试报告</h1>

        <h2>🎯 总体指标</h2>
        {overall_metrics}

        <h2>📈 API端点性能</h2>
        {endpoint_table}

        <h2>📉 性能图表</h2>
        {charts}
    </div>
</body>
</html>
"""

    overall_metrics_html = ""
    endpoint_table_html = """
        <table>
            <thead>
                <tr>
                    <th>API端点</th>
                    <th>p50 (ms)</th>
                    <th>p95 (ms)</th>
                    <th>p99 (ms)</th>
                    <th>错误率 (%)</th>
                    <th>总请求数</th>
                </tr>
            </thead>
            <tbody>
    """

    charts_html = ""

    for name, analyzer in analyzers:
        percentiles = analyzer.get_percentiles()
        throughput = analyzer.get_throughput()
        error_rate = analyzer.get_error_rate()

        # 总体指标
        overall_metrics_html += f"""
        <h3>{name}</h3>
        <div class="metric">
            <div class="metric-label">吞吐量</div>
            <div class="metric-value">{throughput:.2f} req/s</div>
        </div>
        <div class="metric">
            <div class="metric-label">错误率</div>
            <div class="metric-value {'status-pass' if error_rate < 1 else 'status-fail'}">{error_rate:.2f}%</div>
        </div>
        <div class="metric">
            <div class="metric-label">p95 响应时间</div>
            <div class="metric-value">{percentiles['p95']:.1f} ms</div>
        </div>
        <div class="metric">
            <div class="metric-label">p99 响应时间</div>
            <div class="metric-value">{percentiles['p99']:.1f} ms</div>
        </div>
        """

        # 端点表格
        endpoints = analyzer.get_metrics_by_endpoint()
        for endpoint, metrics in endpoints.items():
            endpoint_table_html += f"""
            <tr>
                <td>{endpoint}</td>
                <td>{metrics['p50']:.1f}</td>
                <td>{metrics['p95']:.1f}</td>
                <td>{metrics['p99']:.1f}</td>
                <td class="{'status-pass' if metrics['error_rate'] < 1 else 'status-fail'}">{metrics['error_rate']:.2f}%</td>
                <td>{metrics['total_requests']}</td>
            </tr>
            """

        # 图表占位符
        chart_name = name.replace(' ', '_').lower()
        charts_html += f"""
        <h3>{name}</h3>
        <div class="chart">
            <img src="{chart_name}_response_times.png" alt="响应时间分布">
            <p>响应时间分布</p>
        </div>
        <div class="chart">
            <img src="{chart_name}_throughput.png" alt="吞吐量趋势">
            <p>吞吐量趋势</p>
        </div>
        <div class="chart">
            <img src="{chart_name}_endpoints.png" alt="端点对比">
            <p>端点性能对比</p>
        </div>
        """

    endpoint_table_html += "</tbody></table>"

    html_content = html_template.format(
        overall_metrics=overall_metrics_html,
        endpoint_table=endpoint_table_html,
        charts=charts_html,
    )

    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(html_content)

    print(f"✅ HTML报告已保存: {output_file}")


def main():
    parser = argparse.ArgumentParser(description='组织中心性能测试结果可视化')
    parser.add_argument('input', help='输入JSON文件或目录')
    parser.add_argument('--output', '-o', default='performance_report',
                       help='输出文件前缀 (默认: performance_report)')
    parser.add_argument('--compare', action='store_true',
                       help='对比多个测试结果')
    parser.add_argument('--html', action='store_true',
                       help='生成HTML报告')

    args = parser.parse_args()

    # 查找所有JSON文件
    input_path = Path(args.input)

    if input_path.is_file():
        json_files = [input_path]
    elif input_path.is_dir():
        json_files = list(input_path.glob('*.json'))
    else:
        print(f"❌ 错误: 输入路径不存在: {args.input}")
        sys.exit(1)

    if not json_files:
        print("⚠️  未找到JSON结果文件")
        sys.exit(0)

    print(f"📊 找到 {len(json_files)} 个结果文件")

    # 分析所有文件
    analyzers = []
    for json_file in json_files:
        print(f"📈 分析文件: {json_file.name}")

        try:
            analyzer = PerformanceAnalyzer(str(json_file))
            analyzers.append((json_file.stem, analyzer))

            # 生成图表
            chart_prefix = f"{args.output}_{json_file.stem}"
            plot_response_times(analyzer, f"{chart_prefix}_response_times.png")
            plot_throughput_over_time(analyzer, f"{chart_prefix}_throughput.png")
            plot_endpoint_comparison(analyzer, f"{chart_prefix}_endpoints.png")

        except Exception as e:
            print(f"❌ 分析文件失败 {json_file.name}: {e}")
            continue

    # 生成HTML报告
    if args.html:
        generate_html_report(analyzers, f"{args.output}.html")
        print(f"✅ HTML报告已生成: {args.output}.html")

    print("✅ 所有图表生成完成!")


if __name__ == '__main__':
    main()
