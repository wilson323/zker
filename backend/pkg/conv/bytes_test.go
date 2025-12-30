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

package conv

import (
	"testing"
)

func TestBytesToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "正常转换",
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			name:     "空字节切片",
			input:    []byte{},
			expected: "",
		},
		{
			name:     "nil字节切片",
			input:    nil,
			expected: "",
		},
		{
			name:     "包含特殊字符",
			input:    []byte("hello\n\t世界"),
			expected: "hello\n\t世界",
		},
		{
			name:     "包含空字符",
			input:    []byte{0, 1, 2, 3},
			expected: string([]byte{0, 1, 2, 3}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToString(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
	}{
		{
			name:     "正常转换",
			input:    "hello",
			expected: []byte("hello"),
		},
		{
			name:     "空字符串",
			input:    "",
			expected: nil,
		},
		{
			name:     "包含特殊字符",
			input:    "hello\n\t世界",
			expected: []byte("hello\n\t世界"),
		},
		{
			name:     "包含Unicode字符",
			input:    "🎉🎊",
			expected: []byte("🎉🎊"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBytes(tt.input)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("StringToBytes() = %v, want nil", result)
				}
			} else {
				if string(result) != string(tt.expected) {
					t.Errorf("StringToBytes() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

func TestBytesToStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]byte
		expected []string
	}{
		{
			name:     "正常转换",
			input:    [][]byte{[]byte("a"), []byte("b"), []byte("c")},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "空切片",
			input:    [][]byte{},
			expected: []string{},
		},
		{
			name:     "nil切片",
			input:    nil,
			expected: []string{},
		},
		{
			name:     "包含空字节切片",
			input:    [][]byte{[]byte("a"), []byte{}, []byte("c")},
			expected: []string{"a", "", "c"},
		},
		{
			name:     "包含nil字节切片",
			input:    [][]byte{[]byte("a"), nil, []byte("c")},
			expected: []string{"a", "", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BytesToStringSlice(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("BytesToStringSlice() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("BytesToStringSlice()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestStringToBytesSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected [][]byte
	}{
		{
			name:     "正常转换",
			input:    []string{"a", "b", "c"},
			expected: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
		},
		{
			name:     "空切片",
			input:    []string{},
			expected: [][]byte{},
		},
		{
			name:     "nil切片",
			input:    nil,
			expected: [][]byte{},
		},
		{
			name:     "包含空字符串",
			input:    []string{"a", "", "c"},
			expected: [][]byte{[]byte("a"), nil, []byte("c")},
		},
		{
			name:     "包含特殊字符",
			input:    []string{"hello\n", "世界", "🎉"},
			expected: [][]byte{[]byte("hello\n"), []byte("世界"), []byte("🎉")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToBytesSlice(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("StringToBytesSlice() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if tt.expected[i] == nil {
					if result[i] != nil {
						t.Errorf("StringToBytesSlice()[%d] = %v, want nil", i, result[i])
					}
				} else {
					if string(result[i]) != string(tt.expected[i]) {
						t.Errorf("StringToBytesSlice()[%d] = %q, want %q", i, result[i], tt.expected[i])
					}
				}
			}
		})
	}
}

// 基准测试
func BenchmarkBytesToString(b *testing.B) {
	data := []byte("hello world, this is a test string for benchmarking")
	for i := 0; i < b.N; i++ {
		_ = BytesToString(data)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	data := "hello world, this is a test string for benchmarking"
	for i := 0; i < b.N; i++ {
		_ = StringToBytes(data)
	}
}

func BenchmarkBytesToStringSlice(b *testing.B) {
	data := [][]byte{
		[]byte("string1"),
		[]byte("string2"),
		[]byte("string3"),
		[]byte("string4"),
		[]byte("string5"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToStringSlice(data)
	}
}

func BenchmarkStringToBytesSlice(b *testing.B) {
	data := []string{"string1", "string2", "string3", "string4", "string5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = StringToBytesSlice(data)
	}
}
