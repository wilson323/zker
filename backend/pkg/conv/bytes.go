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

// Package conv 提供统一的类型转换工具
// 遵循以下原则：
// - SOLID原则：单一职责，每个函数只做一件事
// - DRY原则：避免重复的类型转换逻辑
// - KISS原则：简单直接的转换，不使用unsafe
// - 类型安全：显式转换，确保类型安全
package conv

// BytesToString 安全的字节到字符串转换
// 空byte slice返回空字符串，避免nil指针问题
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return string(b)
}

// StringToBytes 安全的字符串到字节转换
// 空字符串返回nil byte slice，避免不必要的内存分配
func StringToBytes(s string) []byte {
	if s == "" {
		return nil
	}
	return []byte(s)
}

// BytesToStringSlice 批量转换字节切片到字符串切片
// 空slice返回空结果，不会返回nil
func BytesToStringSlice(bb [][]byte) []string {
	if len(bb) == 0 {
		return []string{}
	}

	result := make([]string, len(bb))
	for i, b := range bb {
		result[i] = BytesToString(b)
	}
	return result
}

// StringToBytesSlice 批量转换字符串切片到字节切片
// 空slice返回空结果，不会返回nil
func StringToBytesSlice(ss []string) [][]byte {
	if len(ss) == 0 {
		return [][]byte{}
	}

	result := make([][]byte, len(ss))
	for i, s := range ss {
		result[i] = StringToBytes(s)
	}
	return result
}
