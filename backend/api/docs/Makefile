.PHONY: all build run clean help install-dev

# 变量定义
BINARY_NAME=generate_swagger
BINARY_EXT=$(shell go env GOEXE)
CONFIG_FILE=config.json
OUTPUT_DIR=swagger
GO_FILES=generate_swagger.go

# 默认目标
all: build

## build: 编译生成器
build:
	@echo "🔨 编译 API 文档生成器..."
	go build -o $(BINARY_NAME)$(BINARY_EXT) $(GO_FILES)
	@echo "✅ 编译完成: $(BINARY_NAME)$(BINARY_EXT)"

## run: 运行生成器（使用默认配置）
run:
	@echo "🚀 运行 API 文档生成器..."
	go run $(GO_FILES)

## run-config: 使用配置文件运行
run-config:
	@echo "🚀 使用配置文件运行..."
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "❌ 配置文件不存在: $(CONFIG_FILE)"; \
		exit 1; \
	fi
	go run $(GO_FILES) $(CONFIG_FILE)

## clean: 清理生成的文件
clean:
	@echo "🧹 清理生成的文件..."
	@rm -rf $(OUTPUT_DIR)
	@rm -f $(BINARY_NAME)$(BINARY_EXT)
	@echo "✅ 清理完成"

## view-yaml: 查看生成的 YAML 文件
view-yaml:
	@echo "📄 查看 YAML 文件..."
	@if [ -f $(OUTPUT_DIR)/swagger.yaml ]; then \
		head -100 $(OUTPUT_DIR)/swagger.yaml; \
	else \
		echo "❌ 文件不存在，请先运行 'make run'"; \
	fi

## view-json: 查看生成的 JSON 文件
view-json:
	@echo "📄 查看 JSON 文件..."
	@if [ -f $(OUTPUT_DIR)/swagger.json ]; then \
		head -100 $(OUTPUT_DIR)/swagger.json; \
	else \
		echo "❌ 文件不存在，请先运行 'make run'"; \
	fi

## validate: 验证生成的 OpenAPI 规范
validate:
	@echo "🔍 验证 OpenAPI 规范..."
	@if command -v swagger >/dev/null 2>&1; then \
		swagger validate $(OUTPUT_DIR)/swagger.yaml; \
	else \
		echo "⚠️  swagger 命令未安装，跳过验证"; \
		echo "   安装: go install github.com/go-swagger/go-swagger/cmd/swagger@latest"; \
	fi

## install-swagger: 安装 Swagger 命令行工具
install-swagger:
	@echo "📦 安装 Swagger CLI..."
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	@echo "✅ 安装完成"

## serve: 启动本地 Swagger UI
serve:
	@echo "🌐 启动 Swagger UI..."
	@if command -v docker >/dev/null 2>&1; then \
		docker run -p 8080:8080 \
			-e SWAGGER_JSON=/swagger/swagger.yaml \
			-v "$(PWD)/$(OUTPUT_DIR):/swagger" \
			--rm -it swaggerapi/swagger-ui; \
	else \
		echo "❌ Docker 未安装，请先安装 Docker"; \
	fi

## format: 格式化 Go 代码
format:
	@echo "🎨 格式化代码..."
	gofmt -w $(GO_FILES)
	@echo "✅ 格式化完成"

## lint: 运行代码检查
lint:
	@echo "🔍 代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run $(GO_FILES); \
	else \
		echo "⚠️  golangci-lint 未安装，跳过检查"; \
		echo "   安装: https://golangci-lint.run/usage/install/"; \
	fi

## test: 运行测试
test:
	@echo "🧪 运行测试..."
	go test -v ./...

## install-dev: 安装开发工具
install-dev: install-swagger
	@echo "📦 安装其他开发工具..."
	@echo "   golangci-lint: https://golangci-lint.run/usage/install/"
	@echo "   go install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

## help: 显示帮助信息
help:
	@echo "📖 API 文档生成器 - 使用帮助"
	@echo ""
	@echo "可用命令:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | sort
	@echo ""
	@echo "示例:"
	@echo "  make run           # 使用默认配置生成文档"
	@echo "  make build         # 编译生成器"
	@echo "  make serve         # 启动本地 Swagger UI"
	@echo "  make clean         # 清理生成的文件"
