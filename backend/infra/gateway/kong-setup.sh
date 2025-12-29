#!/bin/bash
# Kong API Gateway初始化配置脚本

KONG_ADMIN_URL="http://localhost:8001"
UPSTREAM_SERVICES=(
    "tenant-service:tenant-service:8080"
    "quota-service:quota-service:8080"
    "subscription-service:subscription-service:8080"
    "auth-service:auth-service:8080"
)

echo "=== 开始配置Kong API Gateway ==="

# 1. 添加上游服务
echo "1. 添加上游服务..."
for service in "${UPSTREAM_SERVICES[@]}"; do
    IFS=':' read -r name host port <<< "$service"

    curl -i -X POST "${KONG_ADMIN_URL}/services" \
      --data name="${name}" \
      --data url="http://${host}:${port}"

    echo "✓ Service ${name} added"
done

# 2. 配置路由规则
echo "2. 配置路由规则..."
# 租户服务路由
curl -i -X POST "${KONG_ADMIN_URL}/services/tenant-service/routes" \
  --data paths[]=/api/v1/tenants \
  --data strip_path=false

# 配额服务路由
curl -i -X POST "${KONG_ADMIN_URL}/services/quota-service/routes" \
  --data paths[]=/api/v1/quota \
  --data strip_path=false

# 订阅服务路由
curl -i -X POST "${KONG_ADMIN_URL}/services/subscription-service/routes" \
  --data paths[]=/api/v1/subscriptions \
  --data strip_path=false

echo "✓ Routes configured"

# 3. 配置限流插件
echo "3. 配置限流插件..."
# 全局限流
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data name=rate-limiting \
  --data config.minute=1000 \
  --data config.hour=10000 \
  --data config.policy=local

# 按租户限流（需要配置ACL插件）
# 在后续步骤中动态添加

echo "✓ Rate limiting configured"

# 4. 配置熔断插件
echo "4. 配置熔断插件..."
for service in "${UPSTREAM_SERVICES[@]}"; do
    IFS=':' read -r name host port <<< "$service"

    curl -i -X POST "${KONG_ADMIN_URL}/services/${name}/plugins" \
      --data name=circuit-breaker \
      --data config.error_threshold=50 \
      --data config.volume_threshold=10 \
      --data config.half_open_timeout=30000

    echo "✓ Circuit breaker for ${name}"
done

# 5. 配置JWT认证
echo "5. 配置JWT认证..."
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data name=jwt \
  --data config.uri_param_names=jwt

echo "✓ JWT authentication configured"

# 6. 配置Prometheus监控
echo "6. 配置Prometheus监控..."
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data name=prometheus \
  --data config.per_consumer=true

echo "✓ Prometheus monitoring configured"

# 7. 配置CORS
echo "7. 配置CORS..."
curl -i -X POST "${KONG_ADMIN_URL}/plugins" \
  --data name=cors \
  --data config.origins=* \
  --data config.methods=GET,POST,PUT,DELETE,OPTIONS \
  --data config.headers=Accept,Accept-Version,Content-Length,Content-MD5,Content-Type,Date,Authorization,Host

echo "✓ CORS configured"

echo "=== Kong API Gateway配置完成 ==="
echo "访问Konga管理UI: http://localhost:1337"
echo "访问Kong Admin API: http://localhost:8001"
