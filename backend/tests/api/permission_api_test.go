/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/suite"

	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// PermissionAPITestSuite 权限API契约测试套件
type PermissionAPITestSuite struct {
	suite.Suite
	h *server.Hertz
}

// SetupSuite 设置测试套件
func (s *PermissionAPITestSuite) SetupSuite() {
	s.h = server.Default()
}

// TearDownSuite 清理测试套件
func (s *PermissionAPITestSuite) TearDownSuite() {
	// 清理资源
}

// TestPermissionAPITestSuite 运行测试套件
func TestPermissionAPITestSuite(t *testing.T) {
	suite.Run(t, new(PermissionAPITestSuite))
}

// TestListRoles_Success 测试获取角色列表成功
func (s *PermissionAPITestSuite) TestListRoles_Success() {
	tenantID := "test-tenant-001"

	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/roles",
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
		ut.WithQuery("page_size", "20"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证角色列表结构
	data := result["data"].(map[string]interface{})
	_, hasRoles := data["roles"]
	_, hasTotalCount := data["total_count"]

	s.True(hasRoles, "Response should have 'roles' field")
	s.True(hasTotalCount, "Response should have 'total_count' field")
}

// TestCreateRole_Success 测试创建角色成功
func (s *PermissionAPITestSuite) TestCreateRole_Success() {
	tenantID := "test-tenant-002"

	createReq := map[string]interface{}{
		"role_name":  "测试角色",
		"role_code":  "test_role",
		"role_type":  "custom",
		"description": "用于测试的角色",
	}
	reqBody, _ := json.Marshal(createReq)

	w := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 201, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证角色数据结构
	data := result["data"].(map[string]interface{})
	s.assertRoleData(data)
}

// TestCreateRole_Conflict 测试角色编码冲突
func (s *PermissionAPITestSuite) TestCreateRole_Conflict() {
	tenantID := "test-tenant-003"

	createReq := map[string]interface{}{
		"role_name":  "重复角色",
		"role_code":  "duplicate_role",
		"role_type":  "custom",
		"description": "测试重复",
	}
	reqBody, _ := json.Marshal(createReq)

	// 第一次创建
	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)
	assert.DeepEqual(s.T(), 201, w1.Code)

	// 第二次创建（应该冲突）
	w2 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	// 验证冲突响应
	assert.DeepEqual(s.T(), 409, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertErrorResponse(result)
	assert.DeepEqual(s.T(), errno.ErrRoleAlreadyExistsCode, int32(result["code"].(float64)))
}

// TestGetRole_Success 测试获取角色详情成功
func (s *PermissionAPITestSuite) TestGetRole_Success() {
	tenantID := "test-tenant-004"

	// 先创建一个角色
	createReq := map[string]interface{}{
		"role_name":  "详情测试角色",
		"role_code":  "detail_test_role",
		"role_type":  "custom",
		"description": "用于测试详情",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	roleID := createResp["data"].(map[string]interface{})["role_id"].(string)

	// 获取角色详情
	w2 := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/roles/"+roleID,
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	data := result["data"].(map[string]interface{})
	assert.DeepEqual(s.T(), roleID, data["role_id"])
	s.assertRoleData(data)
}

// TestUpdateRole_Success 测试更新角色成功
func (s *PermissionAPITestSuite) TestUpdateRole_Success() {
	tenantID := "test-tenant-005"

	// 先创建角色
	createReq := map[string]interface{}{
		"role_name":  "更新测试角色",
		"role_code":  "update_test_role",
		"role_type":  "custom",
		"description": "用于测试更新",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	roleID := createResp["data"].(map[string]interface{})["role_id"].(string)

	// 更新角色
	updateReq := map[string]interface{}{
		"role_name":  "更新后的角色名",
		"description": "更新后的描述",
	}
	updateBody, _ := json.Marshal(updateReq)

	w2 := ut.PerformRequest(s.h.Engine, "PUT", "/api/v1/roles/"+roleID,
		ut.WithBody(bytes.NewReader(updateBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	assert.DeepEqual(s.T(), 200, w2.Code)

	var result map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &result)

	s.assertSuccessResponse(result)
	data := result["data"].(map[string]interface{})
	assert.DeepEqual(s.T(), "更新后的角色名", data["role_name"])
}

// TestDeleteRole_Success 测试删除角色成功(软删除)
func (s *PermissionAPITestSuite) TestDeleteRole_Success() {
	tenantID := "test-tenant-006"

	// 先创建角色
	createReq := map[string]interface{}{
		"role_name":  "删除测试角色",
		"role_code":  "delete_test_role",
		"role_type":  "custom",
		"description": "用于测试删除",
	}
	reqBody, _ := json.Marshal(createReq)

	w1 := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/roles",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	var createResp map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &createResp)
	roleID := createResp["data"].(map[string]interface{})["role_id"].(string)

	// 删除角色
	w2 := ut.PerformRequest(s.h.Engine, "DELETE", "/api/v1/roles/"+roleID,
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证204响应
	assert.DeepEqual(s.T(), 204, w2.Code)
}

// TestCheckDataPermission_Success 测试数据权限检查成功
func (s *PermissionAPITestSuite) TestCheckDataPermission_Success() {
	tenantID := "test-tenant-007"
	userID := "test-user-001"

	checkReq := map[string]interface{}{
		"resource_type": "bots",
		"resource_id":   "bot-123",
		"permission":    "read",
		"user_id":       userID,
	}
	reqBody, _ := json.Marshal(checkReq)

	w := ut.PerformRequest(s.h.Engine, "POST", "/api/v1/permissions/check-data",
		ut.WithBody(bytes.NewReader(reqBody)),
		ut.WithHeader("Content-Type", "application/json"),
		ut.WithHeader("X-Tenant-ID", tenantID),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证权限检查结果
	data := result["data"].(map[string]interface{})
	_, hasAllowed := data["allowed"]
	_, hasReason := data["reason"]

	s.True(hasAllowed, "Response should have 'allowed' field")
	s.True(hasReason, "Response should have 'reason' field")
}

// TestPermissionMatrix_Success 测试获取权限矩阵成功
func (s *PermissionAPITestSuite) TestPermissionMatrix_Success() {
	tenantID := "test-tenant-008"

	w := ut.PerformRequest(s.h.Engine, "GET", "/api/v1/permissions/matrix",
		ut.WithHeader("X-Tenant-ID", tenantID),
		ut.WithHeader("Authorization", "Bearer test-token"),
	)

	// 验证响应
	assert.DeepEqual(s.T(), 200, w.Code)

	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)

	s.assertSuccessResponse(result)

	// 验证权限矩阵结构
	data := result["data"].(map[string]interface{})
	_, hasMatrix := data["matrix"]
	_, hasRoles := data["roles"]

	s.True(hasMatrix, "Response should have 'matrix' field")
	s.True(hasRoles, "Response should have 'roles' field")
}

// BenchmarkCheckDataPermission 性能基准测试：数据权限检查
func BenchmarkCheckDataPermission(b *testing.B) {
	h := server.Default()
	tenantID := "benchmark-tenant"
	userID := "benchmark-user"

	checkReq := map[string]interface{}{
		"resource_type": "bots",
		"resource_id":   "bot-benchmark",
		"permission":    "read",
		"user_id":       userID,
	}
	reqBody, _ := json.Marshal(checkReq)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := ut.PerformRequest(h.Engine, "POST", "/api/v1/permissions/check-data",
			ut.WithBody(bytes.NewReader(reqBody)),
			ut.WithHeader("Content-Type", "application/json"),
			ut.WithHeader("X-Tenant-ID", tenantID),
		)
		if w.Code != 200 {
			b.Errorf("Expected status 200, got %d", w.Code)
		}
	}
}

// TestPermissionMiddleware_Concurrent 并发测试：验证权限中间件是并发安全的
func TestPermissionMiddleware_Concurrent(t *testing.T) {
	h := server.Default()
	tenantID := "test-tenant-concurrent"
	userID := "test-user-concurrent"

	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			checkReq := map[string]interface{}{
				"resource_type": "bots",
				"resource_id":   "bot-" + string(rune('0'+index)),
				"permission":    "read",
				"user_id":       userID,
			}
			reqBody, _ := json.Marshal(checkReq)

			w := ut.PerformRequest(h.Engine, "POST", "/api/v1/permissions/check-data",
				ut.WithBody(bytes.NewReader(reqBody)),
				ut.WithHeader("Content-Type", "application/json"),
				ut.WithHeader("X-Tenant-ID", tenantID),
			)

			// 验证响应状态码
			assert.DeepEqual(t, w.Code, 200)

			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// assertSuccessResponse 断言成功响应格式
func (s *PermissionAPITestSuite) assertSuccessResponse(result map[string]interface{}) {
	_, hasCode := result["code"]
	_, hasMessage := result["message"]
	_, hasData := result["data"]
	_, hasRequestID := result["request_id"]
	_, hasTimestamp := result["timestamp"]

	s.True(hasCode, "Response should have 'code' field")
	s.True(hasMessage, "Response should have 'message' field")
	s.True(hasData, "Response should have 'data' field")
	s.True(hasRequestID, "Response should have 'request_id' field")
	s.True(hasTimestamp, "Response should have 'timestamp' field")

	s.Equal(int32(0), int32(result["code"].(float64)))
}

// assertErrorResponse 断言错误响应格式
func (s *PermissionAPITestSuite) assertErrorResponse(result map[string]interface{}) {
	_, hasCode := result["code"]
	_, hasMessage := result["message"]
	_, hasRequestID := result["request_id"]
	_, hasTimestamp := result["timestamp"]

	s.True(hasCode, "Error response should have 'code' field")
	s.True(hasMessage, "Error response should have 'message' field")
	s.True(hasRequestID, "Error response should have 'request_id' field")
	s.True(hasTimestamp, "Error response should have 'timestamp' field")

	s.NotEqual(int32(0), int32(result["code"].(float64)))
}

// assertRoleData 断言角色数据结构
func (s *PermissionAPITestSuite) assertRoleData(data map[string]interface{}) {
	requiredFields := []string{
		"role_id",
		"role_name",
		"role_code",
		"role_type",
		"created_at",
		"updated_at",
	}

	for _, field := range requiredFields {
		_, hasField := data[field]
		s.True(hasField, "Role data should have '"+field+"' field")
	}

	// 验证枚举值
	validTypes := []string{"system", "custom"}

	roleType, hasType := data["role_type"].(string)
	if hasType {
		s.Contains(validTypes, roleType, "role_type should be valid")
	}
}
