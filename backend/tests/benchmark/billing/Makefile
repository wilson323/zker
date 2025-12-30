# Makefile for Billing Performance Benchmark Tests
# Token计量和预算管理系统 - 性能基准测试
#
# 用法:
#   make bench           # 运行所有基准测试
#   make bench-quick     # 快速测试（跳过压力测试）
#   make profile         # 生成CPU和Memory profile
#   make stress          # 运行压力测试
#   make analyze         # 分析测试结果
#   make clean           # 清理测试结果
#
# @author 研发B
# @date 2025-12-30

# 变量定义
SCRIPT_DIR := $(shell pwd)
PROJECT_ROOT := $(shell cd $(SCRIPT_DIR)/../../../.. && pwd)
BENCHMARK_DIR := $(PROJECT_ROOT)/backend/tests/benchmark/billing
RESULTS_DIR := $(PROJECT_ROOT)/backend/tests/results/benchmark
TIMESTAMP := $(shell date +"%Y%m%d_%H%M%S")

# 颜色定义
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help bench bench-quick profile stress analyze clean test-coverage

# 默认目标
.DEFAULT_GOAL := help

## help: 显示帮助信息
help:
	@echo ""
	@echo "========================================="
	@echo "Token计量和预算管理系统 - 性能测试"
	@echo "========================================="
	@echo ""
	@echo "$(BLUE)可用命令:$(NC)"
	@echo "  $(GREEN)make bench$(NC)           - 运行所有基准测试"
	@echo "  $(GREEN)make bench-quick$(NC)     - 快速测试（跳过压力测试）"
	@echo "  $(GREEN)make profile$(NC)         - 生成CPU和Memory profile"
	@echo "  $(GREEN)make stress$(NC)          - 运行压力测试"
	@echo "  $(GREEN)make analyze$(NC)         - 分析测试结果"
	@echo "  $(GREEN)make clean$(NC)           - 清理测试结果"
	@echo "  $(GREEN)make test-coverage$(NC)   - 运行测试并生成覆盖率报告"
	@echo ""
	@echo "$(BLUE)示例:$(NC)"
	@echo "  make bench"
	@echo "  make analyze RESULTS=benchmark_20251230_120000.txt"
	@echo ""

## bench: 运行所有基准测试
bench:
	@echo "$(BLUE)[INFO]$(NC) 运行所有基准测试..."
	@cd $(BENCHMARK_DIR) && \
	go test -bench=. -benchmem -run=^$ \
		-benchtime=1s \
		-timeout=30m \
		. > $(RESULTS_DIR)/benchmark_$(TIMESTAMP).txt 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) 基准测试完成"
	@echo "结果文件: $(RESULTS_DIR)/benchmark_$(TIMESTAMP).txt"

## bench-quick: 快速测试（跳过压力测试）
bench-quick:
	@echo "$(BLUE)[INFO]$(NC) 运行快速基准测试..."
	@cd $(BENCHMARK_DIR) && \
	go test -bench=. -benchmem -run=^$ \
		-benchtime=500ms \
		-timeout=10m \
		-short \
		. > $(RESULTS_DIR)/benchmark_quick_$(TIMESTAMP).txt 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) 快速测试完成"
	@echo "结果文件: $(RESULTS_DIR)/benchmark_quick_$(TIMESTAMP).txt"

## profile: 生成CPU和Memory profile
profile:
	@echo "$(BLUE)[INFO]$(NC) 生成CPU profile..."
	@cd $(BENCHMARK_DIR) && \
	go test -cpuprofile=$(RESULTS_DIR)/cpu_$(TIMESTAMP).prof \
		-bench=BenchmarkCPUProfile \
		-benchtime=10s \
		-timeout=10m \
		. > /dev/null 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) CPU profile生成完成"
	@echo "Profile文件: $(RESULTS_DIR)/cpu_$(TIMESTAMP).prof"
	@echo ""
	@echo "$(BLUE)[INFO]$(NC) 生成Memory profile..."
	@cd $(BENCHMARK_DIR) && \
	go test -memprofile=$(RESULTS_DIR)/memory_$(TIMESTAMP).prof \
		-bench=BenchmarkMemoryProfile \
		-benchtime=10s \
		-timeout=10m \
		. > /dev/null 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) Memory profile生成完成"
	@echo "Profile文件: $(RESULTS_DIR)/memory_$(TIMESTAMP).prof"
	@echo ""
	@if command -v pprof >/dev/null 2>&1; then \
		echo "$(BLUE)[INFO]$(NC) 分析CPU profile..."; \
		pprof -top -nodecount=10 $(RESULTS_DIR)/cpu_$(TIMESTAMP).prof > $(RESULTS_DIR)/cpu_top10_$(TIMESTAMP).txt; \
		echo "$(GREEN)[SUCCESS]$(NC) CPU profile分析完成"; \
		echo ""; \
		echo "$(BLUE)[INFO]$(NC) 分析Memory profile..."; \
		pprof -top -nodecount=10 $(RESULTS_DIR)/memory_$(TIMESTAMP).prof > $(RESULTS_DIR)/memory_top10_$(TIMESTAMP).txt; \
		echo "$(GREEN)[SUCCESS]$(NC) Memory profile分析完成"; \
	fi

## stress: 运行压力测试
stress:
	@echo "$(YELLOW)[WARNING]$(NC) 压力测试需要较长时间，请耐心等待..."
	@echo "$(BLUE)[INFO]$(NC) 运行内存泄漏检测..."
	@cd $(BENCHMARK_DIR) && \
	go test -v -run=TestMemoryLeak \
		-timeout=10m \
		. > $(RESULTS_DIR)/memory_leak_$(TIMESTAMP).txt 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) 内存泄漏检测完成"
	@echo ""
	@echo "$(BLUE)[INFO]$(NC) 运行持续高负载测试（1小时）..."
	@cd $(BENCHMARK_DIR) && \
	go test -v -run=TestStressTokenRecording \
		-timeout=2h \
		. > $(RESULTS_DIR)/stress_test_$(TIMESTAMP).txt 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) 持续高负载测试完成"
	@echo ""
	@echo "$(BLUE)[INFO]$(NC) 运行峰值负载测试（1分钟）..."
	@cd $(BENCHMARK_DIR) && \
	go test -v -run=TestPeakLoad \
		-timeout=30m \
		. > $(RESULTS_DIR)/peak_load_$(TIMESTAMP).txt 2>&1
	@echo "$(GREEN)[SUCCESS]$(NC) 峰值负载测试完成"

## analyze: 分析测试结果
analyze:
	@if [ -z "$(RESULTS)" ]; then \
		LATEST_RESULT=$$(ls -t $(RESULTS_DIR)/benchmark_*.txt 2>/dev/null | head -1); \
		if [ -z "$$LATEST_RESULT" ]; then \
			echo "$(RED)[ERROR]$(NC) 找不到测试结果文件"; \
			echo "用法: make analyze RESULTS=benchmark_20251230_120000.txt"; \
			exit 1; \
		fi; \
		echo "$(BLUE)[INFO]$(NC) 使用最新测试结果: $$LATEST_RESULT"; \
		RESULT_FILE=$$LATEST_RESULT; \
	else \
		RESULT_FILE=$(RESULTS_DIR)/$(RESULTS); \
	fi
	@echo "$(BLUE)[INFO]$(NC) 分析测试结果: $$RESULT_FILE"
	@$(SCRIPT_DIR)/analyze_results.sh $$RESULT_FILE

## clean: 清理测试结果
clean:
	@echo "$(BLUE)[INFO]$(NC) 清理测试结果..."
	@rm -rf $(RESULTS_DIR)/*.txt
	@rm -rf $(RESULTS_DIR)/*.prof
	@rm -rf $(RESULTS_DIR)/*.md
	@echo "$(GREEN)[SUCCESS]$(NC) 测试结果已清理"

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "$(BLUE)[INFO]$(NC) 运行测试并生成覆盖率报告..."
	@cd $(BENCHMARK_DIR) && \
	go test -coverprofile=$(RESULTS_DIR)/coverage_$(TIMESTAMP).out \
		-covermode=atomic \
		. > $(RESULTS_DIR)/coverage_$(TIMESTAMP).txt 2>&1
	@cd $(BENCHMARK_DIR) && \
	go tool cover -html=$(RESULTS_DIR)/coverage_$(TIMESTAMP).out \
		-o $(RESULTS_DIR)/coverage_$(TIMESTAMP).html
	@echo "$(GREEN)[SUCCESS]$(NC) 测试完成"
	@echo "覆盖率报告: $(RESULTS_DIR)/coverage_$(TIMESTAMP).html"
	@echo "覆盖率数据: $(RESULTS_DIR)/coverage_$(TIMESTAMP).txt"

## install-deps: 安装测试依赖
install-deps:
	@echo "$(BLUE)[INFO]$(NC) 安装测试依赖..."
	@go install github.com/google/pprof@latest
	@echo "$(GREEN)[SUCCESS]$(NC) 依赖安装完成"

## setup: 初始化测试环境
setup:
	@echo "$(BLUE)[INFO]$(NC) 初始化测试环境..."
	@mkdir -p $(RESULTS_DIR)
	@echo "$(GREEN)[SUCCESS]$(NC) 测试环境初始化完成"

## all: 运行完整测试流程
all: setup bench profile stress analyze
	@echo ""
	@echo "========================================="
	@echo "$(GREEN)完整测试流程已完成！$(NC)"
	@echo "========================================="
	@echo ""
	@echo "测试结果目录: $(RESULTS_DIR)"
	@echo ""
	@ls -lh $(RESULTS_DIR)
	@echo ""
