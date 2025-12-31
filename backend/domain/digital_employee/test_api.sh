#!/bin/bash
# 数字员工管理API测试脚本

BASE_URL="http://localhost:8080"
TENANT_ID="test-tenant-$(date +%s)"

echo "======================================"
echo "数字员工管理API测试"
echo "======================================"
echo ""

# 测试1：创建员工画像
echo "📝 测试1：创建员工画像"
EMPLOYEE_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/digital-employees" \
  -H "Content-Type: application/json" \
  -d "{
    \"tenant_id\": \"${TENANT_ID}\",
    \"bot_id\": \"test-bot-001\",
    \"name\": \"技术支持专家\",
    \"role\": \"tech_support\",
    \"skills\": [\"技术支持\", \"故障排查\", \"API开发\"],
    \"status\": \"active\"
  }")
echo "$EMPLOYEE_RESPONSE" | jq '.'
EMPLOYEE_ID=$(echo "$EMPLOYEE_RESPONSE" | jq -r '.data.employee_id // empty')
echo "✅ 员工ID: $EMPLOYEE_ID"
echo ""

# 测试2：列出员工画像
echo "📋 测试2：列出员工画像"
curl -s -X GET "${BASE_URL}/api/v1/digital-employees?tenant_id=${TENANT_ID}&page=1&page_size=10" | jq '.'
echo ""

# 测试3：获取员工详情
if [ -n "$EMPLOYEE_ID" ]; then
  echo "👤 测试3：获取员工详情"
  curl -s -X GET "${BASE_URL}/api/v1/digital-employees/${EMPLOYEE_ID}" | jq '.'
  echo ""
fi

# 测试4：智能任务自动分配
echo "🤖 测试4：智能任务自动分配（技能匹配）"
ASSIGN_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/digital-employees/tasks/auto-assign" \
  -H "Content-Type: application/json" \
  -d "{
    \"task_id\": \"task-$(uuidgen)\",
    \"task_type\": \"tech_support\",
    \"priority\": \"high\",
    \"required_skills\": [\"技术支持\", \"故障排查\"]
  }")
echo "$ASSIGN_RESPONSE" | jq '.'
ASSIGNMENT_ID=$(echo "$ASSIGN_RESPONSE" | jq -r '.data.assignment_id // empty')
echo "✅ 分配ID: $ASSIGNMENT_ID"
echo ""

# 测试5：完成任务
if [ -n "$ASSIGNMENT_ID" ]; then
  echo "✅ 测试5：完成任务"
  curl -s -X PUT "${BASE_URL}/api/v1/digital-employees/tasks/${ASSIGNMENT_ID}/complete" \
    -H "Content-Type: application/json" \
    -d '{
      "result": "任务成功完成，客户问题已解决"
    }' | jq '.'
  echo ""
fi

# 测试6：获取员工绩效
if [ -n "$EMPLOYEE_ID" ]; then
  echo "📊 测试6：获取员工绩效"
  curl -s -X GET "${BASE_URL}/api/v1/digital-employees/${EMPLOYEE_ID}/performance?period=daily" | jq '.'
  echo ""
fi

# 测试7：获取团队绩效
echo "📊 测试7：获取团队绩效"
curl -s -X GET "${BASE_URL}/api/v1/digital-employees/performance/team?tenant_id=${TENANT_ID}&period=daily" | jq '.'
echo ""

# 测试8：更新员工画像
if [ -n "$EMPLOYEE_ID" ]; then
  echo "✏️  测试8：更新员工画像"
  curl -s -X PUT "${BASE_URL}/api/v1/digital-employees/${EMPLOYEE_ID}" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"高级技术支持专家\",
      \"skills\": [\"技术支持\", \"故障排查\", \"API开发\", \"系统优化\"],
      \"specialization\": \"后端系统性能优化\"
    }" | jq '.'
  echo ""
fi

# 测试9：获取员工任务列表
if [ -n "$EMPLOYEE_ID" ]; then
  echo "📋 测试9：获取员工任务列表"
  curl -s -X GET "${BASE_URL}/api/v1/digital-employees/${EMPLOYEE_ID}/tasks?page=1&page_size=10" | jq '.'
  echo ""
fi

# 测试10：删除员工画像（软删除）
if [ -n "$EMPLOYEE_ID" ]; then
  echo "🗑️  测试10：删除员工画像"
  curl -s -X DELETE "${BASE_URL}/api/v1/digital-employees/${EMPLOYEE_ID}" | jq '.'
  echo ""
fi

echo "======================================"
echo "✅ 所有测试完成！"
echo "======================================"
