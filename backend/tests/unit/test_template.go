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

// Package unit 提供企业级单元测试模板和工具
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// =====================================================
// 测试套件模板
// =====================================================

// ServiceTestSuite 服务层测试套件模板
type ServiceTestSuite struct {
	suite.Suite
	Ctx     context.Context
	Cleanup func()
}

// SetupTest 测试前准备
func (s *ServiceTestSuite) SetupTest() {
	s.Ctx = context.Background()
}

// TearDownTest 测试后清理
func (s *ServiceTestSuite) TearDownTest() {
	if s.Cleanup != nil {
		s.Cleanup()
	}
}

// =====================================================
// Mock定义辅助函数
// =====================================================

// MockExpectations 设置Mock期望
type MockExpectations struct {
	Calls []MockCall
}

type MockCall struct {
	Method string
	Args   []interface{}
	Return []interface{}
}

// NewMockExpectations 创建Mock期望配置
func NewMockExpectations() *MockExpectations {
	return &MockExpectations{}
}

// AddCall 添加Mock调用
func (m *MockExpectations) AddCall(method string, args []interface{}, returns []interface{}) {
	m.Calls = append(m.Calls, MockCall{
		Method: method,
		Args:   args,
		Return: returns,
	})
}

// =====================================================
// 测试数据构建器
// =====================================================

// TestDataBuilder 测试数据构建器
type TestDataBuilder struct {
	data map[string]interface{}
}

// NewTestDataBuilder 创建测试数据构建器
func NewTestDataBuilder() *TestDataBuilder {
	return &TestDataBuilder{
		data: make(map[string]interface{}),
	}
}

// Set 设置字段值
func (b *TestDataBuilder) Set(key string, value interface{}) *TestDataBuilder {
	b.data[key] = value
	return b
}

// Build 构建数据
func (b *TestDataBuilder) Build() map[string]interface{} {
	return b.data
}

// =====================================================
// 断言辅助函数
// =====================================================

// AssertNoError 断言无错误，失败时打印详细信息
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	if err != nil {
		t.Helper()
		t.Fatalf("Unexpected error: %v\n%+v", err, msgAndArgs)
	}
}

// AssertErrorType 断言错误类型
func AssertErrorType(t *testing.T, err error, expectedType interface{}) {
	t.Helper()
	assert.NotNil(t, err, "Expected error but got nil")
	assert.IsType(t, expectedType, err)
}

// AssertContains 断言包含
func AssertContains(t *testing.T, container interface{}, contains interface{}) {
	t.Helper()
	assert.Contains(t, container, contains)
}

// AssertEmpty 断言为空
func AssertEmpty(t *testing.T, obj interface{}) {
	t.Helper()
	assert.Empty(t, obj)
}

// AssertNotEmpty 断言非空
func AssertNotEmpty(t *testing.T, obj interface{}) {
	t.Helper()
	assert.NotEmpty(t, obj)
}

// =====================================================
// Mock辅助函数
// =====================================================

// SetupMockReturn 设置Mock返回值
func SetupMockReturn(mockObj interface{}, method string, returns []interface{}) {
	// 使用反射或代码生成设置Mock返回值
	// 这里提供一个示例接口
}

// =====================================================
// 表格驱动测试辅助函数
// =====================================================

// TableTestCase 表格驱动测试用例
type TableTestCase struct {
	Name     string
	Input    interface{}
	Expected interface{}
	Error    error
	Setup    func() interface{}
	Teardown func(interface{})
}

// RunTableTests 运行表格驱动测试
func RunTableTests(t *testing.T, tests []TableTestCase, testFunc func(*testing.T, TableTestCase)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Setup != nil {
				defer tt.Teardown(tt.Setup())
			}
			testFunc(t, tt)
		})
	}
}

// =====================================================
// 上下文辅助函数
// =====================================================

// ContextWithTimeout 创建带超时的上下文
func ContextWithTimeout(t *testing.T) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), testingTimeout)
}

const testingTimeout = 30 // 秒

// =====================================================
// 测试清理辅助函数
// =====================================================

// DeferClean 延迟清理
func DeferClean(t *testing.T, cleanup func()) {
	t.Cleanup(cleanup)
}

// =====================================================
// 示例测试用例
// =====================================================

// ExampleUnitTest 示例单元测试
func ExampleUnitTest() {
	tests := []struct {
		name    string
		input   int
		want    int
		wantErr bool
	}{
		{
			name:    "正常输入",
			input:   10,
			want:    20,
			wantErr: false,
		},
		{
			name:    "错误输入",
			input:   -1,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SomeFunction(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SomeFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SomeFunction() = %v, want %v", got, tt.want)
			}
		})
	}
}

// SomeFunction 示例函数
func SomeFunction(input int) (int, error) {
	return input * 2, nil
}
