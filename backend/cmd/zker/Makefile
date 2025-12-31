# ZKER CLI Makefile
# 用于构建ZKER命令行工具

APP_NAME=zker
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go相关参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

.PHONY: build build-all clean install test help

# 默认目标
all: build

## build: 构建ZKER CLI（当前平台）
build:
	@echo "Building $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME) ./backend/cmd/zker
	@echo "Build complete: bin/$(APP_NAME)"

## build-all: 构建所有平台的二进制文件
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./backend/cmd/zker
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./backend/cmd/zker
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./backend/cmd/zker
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./backend/cmd/zker
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/$(APP_NAME)-windows-amd64.exe ./backend/cmd/zker
	@echo "Build complete for all platforms"

## install: 安装ZKER CLI到$GOPATH/bin
install:
	@echo "Installing $(APP_NAME)..."
	$(GOCMD) install $(LDFLAGS) ./backend/cmd/zker
	@echo "Installation complete"

## clean: 清理构建文件
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	@echo "Clean complete"

## test: 运行测试
test:
	@echo "Running tests..."
	$(GOTEST) -v -cover ./backend/cmd/zker/...

## test-coverage: 运行测试并生成覆盖率报告
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./backend/cmd/zker/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## lint: 代码检查
lint:
	@echo "Running linters..."
	golangci-lint run ./backend/cmd/zker/...

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./backend/cmd/zker/...

## deps: 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## release: 构建发布版本
release: clean build-all
	@echo "Creating release packages..."
	@cd bin && \
	for file in $(APP_NAME)-*; do \
		if [[ $$file == *.exe ]]; then \
			zip $$file.zip $$file; \
		else \
			tar czf $$file.tar.gz $$file; \
		fi \
	done
	@echo "Release packages created in bin/"

## help: 显示帮助信息
help:
	@echo "ZKER CLI - Makefile命令"
	@echo ""
	@echo "使用方法:"
	@echo "  make [target]"
	@echo ""
	@echo "可用命令:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo ""
	@echo "示例:"
	@echo "  make build          # 构建当前平台"
	@echo "  make build-all      # 构建所有平台"
	@echo "  make install        # 安装到系统"
	@echo "  make test           # 运行测试"
	@echo "  make clean          # 清理构建文件"
