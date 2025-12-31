#!/bin/bash
# ================================================================================
# Coze Studio - Secret 生成脚本
# ================================================================================
# 用途：从用户输入自动生成 base64 编码的 Secret YAML
# 使用方法：./generate-secret.sh
# 作者：ZKER DevOps Team
# 创建时间：2025-01-01
# ================================================================================

set -e

SECRET_FILE="secret.yaml"
NAMESPACE="coze-studio"

echo "=================================================="
echo "  Coze Studio - Secret 生成向导"
echo "=================================================="
echo ""
echo "请输入以下敏感信息（输入后将自动 base64 编码）："
echo ""

# 输入函数
input_secret() {
  local prompt=$1
  local default=$2
  local value

  if [ -n "$default" ]; then
    read -p "$prompt [$default]: " value
    value=${value:-$default}
  else
    read -p "$prompt: " value
  fi

  echo "$value" | base64 -w 0
}

# 开始收集信息
echo "========== 数据库配置 =========="
DB_PASSWORD=$(input_secret "数据库密码")
DB_ROOT_PASSWORD=$(input_secret "数据库 Root 密码")

echo ""
echo "========== Redis 配置 =========="
REDIS_PASSWORD=$(input_secret "Redis 密码")

echo ""
echo "========== MinIO 配置 =========="
MINIO_ACCESS_KEY=$(input_secret "MinIO Access Key" "minioadmin")
MINIO_SECRET_KEY=$(input_secret "MinIO Secret Key" "minioadmin")

echo ""
echo "========== AI 模型 API 密钥 =========="
OPENAI_API_KEY=$(input_secret "OpenAI API Key (可选，留空跳过)" "")
AZURE_OPENAI_API_KEY=$(input_secret "Azure OpenAI API Key (可选，留空跳过)" "")

echo ""
echo "========== JWT 配置 =========="
JWT_SIGNING_KEY=$(input_secret "JWT 签名密钥" "$(openssl rand -base64 32)")

echo ""
echo "========== 加密配置 =========="
ENCRYPTION_KEY=$(input_secret "数据加密密钥" "$(openssl rand -base64 32)")

echo ""
echo "========== SMTP 配置 (可选) =========="
SMTP_HOST=$(input_secret "SMTP 主机 (可选，留空跳过)" "")
SMTP_PORT=$(input_secret "SMTP 端口 (可选，留空跳过)" "")
SMTP_USERNAME=$(input_secret "SMTP 用户名 (可选，留空跳过)" "")
SMTP_PASSWORD=$(input_secret "SMTP 密码 (可选，留空跳过)" "")

# 生成 YAML 文件
cat > "$SECRET_FILE" <<EOF
# ================================================================================
# Coze Studio - 生产环境 Secret
# ================================================================================
# 警告：此文件包含敏感信息，请勿提交到 Git 仓库！
# 生成时间：$(date)
# ================================================================================

apiVersion: v1
kind: Secret
metadata:
  name: coze-studio-secret
  namespace: $NAMESPACE
  labels:
    app: coze-studio
    component: secret
type: Opaque
data:
  # 数据库密码
  db-password: "$DB_PASSWORD"
  db-root-password: "$DB_ROOT_PASSWORD"

  # Redis 密码
  redis-password: "$REDIS_PASSWORD"

  # MinIO 凭证
  minio-access-key: "$MINIO_ACCESS_KEY"
  minio-secret-key: "$MINIO_SECRET_KEY"

  # AI 模型 API 密钥
  openai-api-key: "$OPENAI_API_KEY"
  azure-openai-api-key: "$AZURE_OPENAI_API_KEY"

  # JWT 签名密钥
  jwt-signing-key: "$JWT_SIGNING_KEY"

  # 加密密钥
  encryption-key: "$ENCRYPTION_KEY"

  # SMTP 配置
  smtp-host: "$SMTP_HOST"
  smtp-port: "$SMTP_PORT"
  smtp-username: "$SMTP_USERNAME"
  smtp-password: "$SMTP_PASSWORD"
EOF

echo ""
echo "=================================================="
echo "  ✅ Secret 文件已生成: $SECRET_FILE"
echo "=================================================="
echo ""
echo "下一步操作："
echo "  1. 检查生成的文件内容：cat $SECRET_FILE"
echo "  2. 应用到 Kubernetes：kubectl apply -f $SECRET_FILE"
echo "  3. 验证 Secret：kubectl get secret coze-studio-secret -n $NAMESPACE -o yaml"
echo ""
echo "⚠️  重要提示："
echo "  - 请妥善保管此文件，勿上传到 Git 仓库"
echo "  - 建议生产环境使用 External Secrets Operator"
echo "  - 定期轮换密钥"
echo ""
